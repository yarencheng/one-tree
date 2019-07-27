package main

import (
	"context"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"
	"time"
	kafkaconsumergroup "yarencheng/one-tree/go-src/kafka-consumergroup"
	"yarencheng/one-tree/go-src/kafka-consumergroup/config"
	"yarencheng/one-tree/go-src/pb"

	"flag"
	"strings"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
)

var (
	debug             = flag.Bool("debug", false, "Enable debug mode")
	jsonLog           = flag.Bool("jsonLog", false, "Enable JSON format logger")
	defaultConfigPath = flag.String("defaultConfig", "kafka.config.yaml", "Path of the default configuration file")
	configPath        = flag.String("config", "", "Path of the configuration file")
)

func main() {
	flag.Parse()

	initLog()

	var err error
	if *configPath == "" {
		err = config.Init(*defaultConfigPath)
	} else {
		err = config.Init(*defaultConfigPath, *configPath)
	}
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := kafkaconsumergroup.New(&pb.EchoEvent{})
	if err != nil {
		log.Fatal(err)
	}

	err = consumer.Start()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for message := range consumer.In() {
			log.Infof("received message [%#v]]", message)
		}
	}()

	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	wait := make(chan int, 1)

	go func() {
		if err := consumer.Shutdown(ctx); err != nil {
			log.Fatal("Server Shutdown:", err)
		}
		close(wait)
	}()

	select {
	case <-ctx.Done():
	case <-wait:
	}

	e := ctx.Err()
	if e != nil {
		log.Warnf("Server failed. err=[%v]", e)
	} else {
		log.Println("Server exited")
	}
}

func initLog() {
	// log.SetFormatter(&log.JSONFormatter{})

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetReportCaller(true)

	if *jsonLog {
		log.SetFormatter(&log.JSONFormatter{
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				s := strings.Split(f.Function, ".")
				funcname := s[len(s)-1]
				_, filename := path.Split(f.File)
				return funcname, filename
			},
		})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				s := strings.Split(f.Function, ".")
				funcname := s[len(s)-1]
				_, filename := path.Split(f.File)
				return funcname, filename
			},
		})
	}

	sarama.Logger = log.StandardLogger()
}
