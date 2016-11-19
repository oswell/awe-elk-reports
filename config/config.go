package config

type KafkaConfig struct {
	Brokers          []string    `required:"true"`
	Topic            string      `default:"aws-reports"`
	TLSKey           string
	TLSCertificate   string
	TLSCACertificate string      
	TLSVerify        bool        `default:false`
}

// Configuration for MultiConfig
type Cfg struct {
	// S3 bucket in which billing reports are stored.
	Bucket   string   `required:"true"`

	// Report path value assigned in the billing report configuration
	Prefix   string   `default:""`

	DBUrl    string   `required:"true"`

	LogLevel string   `default:"info"`

	Kafka KafkaConfig `required:"true"`
}
