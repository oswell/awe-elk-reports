package config

type KafkaConfig struct {
	Brokers          []string    `json:"brokers",required:"true"`
	Topic            string      `json:"topic",required:"true"`
	TLSKey           string      `json:"tlskey",required:"true"`
	TLSCertificate   string      `json:"tlscertificate",required:"true"`
	TLSCACertificate string      `json:"tlscacertificate",required:"true"`
	TLSVerify        bool        `json:"tlsverify",required:"true"`
}

// Configuration for MultiConfig
type Configuration struct {
	// S3 bucket in which billing reports are stored.
	Bucket string `json:"bucket",required:"true"`

	// Report path value assigned in the billing report configuration
	Prefix string `json:"prefix",required:"false"`

	DBUrl string `json:"dburl",required:"true"`

	Kafka KafkaConfig `json:"kafka,required:"true"`
}
