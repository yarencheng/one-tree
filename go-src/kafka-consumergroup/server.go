package kafkaconsumergroup

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/gogo/protobuf/proto"
	log "github.com/sirupsen/logrus"

	"context"
	"yarencheng/one-tree/go-src/kafka-consumergroup/config"

	"github.com/Shopify/sarama"
)

type Server struct {
	consumer sarama.ConsumerGroup
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	err      error
	in       chan proto.Message
	tvpe     reflect.Type
}

func (s *Server) In() <-chan proto.Message {

	return s.in

}

func New(payload proto.Message) (*Server, error) {

	version, err := sarama.ParseKafkaVersion(config.Default.Kafka.Version)
	if err != nil {
		e := fmt.Errorf("Error parsing Kafka version: %v", err)
		return nil, e
	}

	/**
	 * Construct a new Sarama configuration.
	 * The Kafka cluster version has to be defined before the consumer/producer is initialized.
	 */
	sconfig := sarama.NewConfig()
	sconfig.Version = version

	consumer, err := sarama.NewConsumerGroup(config.Default.Kafka.Brokers, config.Default.Kafka.Group, sconfig)
	if err != nil {
		return nil, fmt.Errorf("Failed to create Sarama consumer. err: [%s]", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	s := &Server{
		consumer: consumer,
		ctx:      ctx,
		cancel:   cancel,
		in:       make(chan proto.Message, 10),
		tvpe:     reflect.TypeOf(payload).Elem(),
	}

	return s, nil

}

func (s *Server) Shutdown(ctx context.Context) error {

	log.Infof("Server exiting")

	s.cancel()

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

		for {
			if err := s.consumer.Consume(s.ctx, []string{config.Default.Kafka.Topic}, s); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}

			defer func() {
				if err := s.consumer.Close(); err != nil {
					s.err = fmt.Errorf("Failed to shut down consumer cleanly. err: [%v]", err.Error())
					log.Error(s.err)
				}
			}()

			defer close(s.in)

			select {
			case <-s.ctx.Done():
				return
			default:
			}
		}
	}()

	return nil
}

func (s *Server) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (s *Server) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (s *Server) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {

		pb := reflect.New(s.tvpe).Interface().(proto.Message)

		err := proto.Unmarshal(message.Value, pb)
		if err != nil {
			log.Warnf("Skip message due to [%s] ", err)
			continue
		}

		s.in <- pb

		log.Debugf("received message [%#v] from topic [%s] partition/offset [%d/%d]",
			pb,
			message.Topic,
			message.Partition,
			message.Offset,
		)
		session.MarkMessage(message, "")
	}

	return nil
}
