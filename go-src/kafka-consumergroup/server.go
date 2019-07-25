package kafkaconsumergroup

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"

	"context"
	"yarencheng/one-tree/go-src/kafka-consumergroup/config"

	"github.com/Shopify/sarama"
)

type Server struct {
	Echo   sarama.ConsumerGroup
	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
	err    error
}

func New() (*Server, error) {

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
		Echo:   consumer,
		ctx:    ctx,
		cancel: cancel,
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

		consumer := Consumer{}

		for {
			if err := s.Echo.Consume(s.ctx, []string{config.Default.Kafka.Topic}, &consumer); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}

			defer func() {
				if err := s.Echo.Close(); err != nil {
					s.err = fmt.Errorf("Failed to shut down consumer cleanly. err: [%v]", err.Error())
					log.Error(s.err)
				}
			}()

			select {
			case <-s.ctx.Done():
				return
			default:
			}
		}
	}()

	return nil
}

type Consumer struct{}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		log.Printf("received message [%s] from topic [%s] partition/offset [%d/%d]",
			string(message.Value),
			message.Topic,
			message.Partition,
			message.Offset,
		)
		session.MarkMessage(message, "")
	}

	return nil
}
