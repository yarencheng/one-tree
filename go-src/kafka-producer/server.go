package kafkaproducer

import (
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"context"
	"yarencheng/one-tree/go-src/kafka-producer/config"

	"github.com/Shopify/sarama"
)

type Server struct {
	Echo    sarama.SyncProducer
	wg      sync.WaitGroup
	stopped chan struct{}
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
		Echo:    producer,
		stopped: make(chan struct{}),
	}

	return s, nil

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
			if err := s.Echo.Close(); err != nil {
				log.Errorf("Failed to shut down data collector cleanly. err: [%v]", err.Error())
			}
		}()

		log.Info("Start loop")

		ticker := time.NewTicker(1 * time.Second)

		for {
			select {
			case <-s.stopped:
				return
			case <-ticker.C:

				message := fmt.Sprintf("It is %s", time.Now())

				partition, offset, err := s.Echo.SendMessage(&sarama.ProducerMessage{
					Topic: config.Default.Kafka.Topic,
					Value: sarama.StringEncoder(message),
				})

				if err != nil {
					log.Warnf("Failed to send message. err: [%s]", err)
				} else {
					log.Debugf("Sent message [%s] to topic [%s] partition/offset [%v / %v]",
						message, config.Default.Kafka.Topic, partition, offset,
					)
				}
			}
		}
	}()

	return nil
}
