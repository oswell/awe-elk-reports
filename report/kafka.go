package report

import (
    "fmt"
    "crypto/tls"
	"io/ioutil"
	"crypto/x509"
    "github.com/Shopify/sarama"
    "github.com/Sirupsen/logrus"
    "github.com/oswell/aws-elk-reports/config"
)

type Kafka struct {
    Config        config.KafkaConfig
    producer      *sarama.SyncProducer
}

func (k *Kafka) getProducer() (*sarama.SyncProducer, error) {
    if k.producer == nil {
        config := sarama.NewConfig()
    	config.Producer.RequiredAcks = sarama.WaitForAll
    	config.Producer.Retry.Max = 10

    	tlsConfig, err := k.createTLSConfiguration() ; if tlsConfig != nil {
    		config.Net.TLS.Config = tlsConfig
    		config.Net.TLS.Enable = true
    	}

    	producer, err := sarama.NewSyncProducer(k.Config.Brokers, config)
    	if err != nil {
            logrus.Info("Error setting up sync producer")
    		return nil, fmt.Errorf("Failed to start Sarama producer: %s", err)
    	}

        k.producer = &producer
    }
	return k.producer, nil
}

func (k *Kafka) Produce(jsonBytes []byte) (error) {

    if k.producer == nil {
        _, err := k.getProducer() ; if err != nil {
            return err
        }
    }

    message := &sarama.ProducerMessage{Topic: k.Config.Topic, Value: sarama.ByteEncoder(jsonBytes)}
    _, _, err := (*k.producer).SendMessage(message)

    if err != nil {
        logrus.Error("FAILED to send message: %s", err)
    }

    return nil
}

// createTLSConfiguration configures TLS support for kafka connections
func (k *Kafka) createTLSConfiguration() (tlsConfig *tls.Config, err error) {

	if k.Config.TLSCertificate != "" && k.Config.TLSKey != "" && k.Config.TLSCACertificate != "" {

    	cert, err := tls.LoadX509KeyPair(k.Config.TLSCertificate, k.Config.TLSKey) ; if err != nil {
            return nil, fmt.Errorf("Error loading key pair, %s", err)
		}

		caCert, err := ioutil.ReadFile(k.Config.TLSCACertificate) ; if err != nil {
			return nil, fmt.Errorf("Error loading CA certificate, %s", err)
		}

		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caCert)

		tlsConfig = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            certPool,
			InsecureSkipVerify: k.Config.TLSVerify,
		}
	}

	return tlsConfig, nil
}
