package kafkaproducer

import (
	"fmt"
	"sync"

	"github.com/gogo/protobuf/proto"
	log "github.com/sirupsen/logrus"

	"context"
	"yarencheng/one-tree/go-src/kafka-producer/config"

	"github.com/Shopify/sarama"
)

type Server struct {
	producer sarama.SyncProducer
	wg       sync.WaitGroup
	stopped  chan struct{}
	err      error
	out      chan proto.Message
}

func New() (*Server, error) {

	sconfig := sarama.NewConfig()
	sconfig.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	sconfig.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	sconfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(config.Default.Kafka.Brokers, sconfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to start Sarama producer. err: [%s]", err)
	}

	s := &Server{
		producer: producer,
		stopped:  make(chan struct{}),
		out:      make(chan proto.Message, 10),
	}

	return s, nil

}

func (s *Server) Out() chan<- proto.Message {

	return s.out

}

func (s *Server) Shutdown(ctx context.Context) error {

	log.Infof("Server exiting")

	close(s.stopped)

	done := make(chan struct{})

	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}

	e := ctx.Err()
	if e != nil {
		log.Warnf("Close server failed. err=[%s]", e)
		return fmt.Errorf("Close server failed. err=[%s]", e)
	} else if s.err != nil {
		log.Warnf("Close server failed. err=[%s]", s.err)
		return s.err
	} else {
		log.Infof("Server exited")
	}

	return nil
}

func (s *Server) Start() error {

	s.wg.Add(1)

	go func() {

		defer s.wg.Done()

		defer func() {
			if err := s.producer.Close(); err != nil {
				s.err = fmt.Errorf("Failed to shut down producer cleanly. err: [%v]", err.Error())
				log.Error(s.err)
			}
		}()

		defer close(s.out)

		log.Info("Start loop")

		for {
			select {
			case <-s.stopped:
				return
			case message := <-s.out:

				payloadBytes, err := proto.Marshal(message)
				if err != nil {
					log.Warnf("Skip message due to [%s] ", err)
					continue
				}

				partition, offset, err := s.producer.SendMessage(&sarama.ProducerMessage{
					Topic: config.Default.Kafka.Topic,
					Value: sarama.StringEncoder(payloadBytes),
				})

				if err != nil {
					log.Warnf("Failed to send message. err: [%s]", err)
				} else {
					log.Debugf("Sent message [%#v] to topic [%s] partition/offset [%v / %v]",
						message, config.Default.Kafka.Topic, partition, offset,
					)
				}
			}
		}
	}()

	return nil
}
