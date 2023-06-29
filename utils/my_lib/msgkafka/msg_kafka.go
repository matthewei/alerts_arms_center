package msgkafka

type Configurations struct {
	Enabled             bool   `yaml:"enabled"`
	DebugMode           bool   `yaml:"debug_mode"`
	Version             string `yaml:"version"`
	Brokers             string `yaml:"brokers"`
	Topic               string `yaml:"topic"`
	RequiredAck         int16  `yaml:"required_ack"`
	ReturnSuccesses     bool   `yaml:"return_successes"`
	ReturnErrors        bool   `yaml:"return_errors"`
	RetryMax            int    `yaml:"retry_max"`
	MaxMessageBytes     int    `yaml:"max_message_bytes"`
	AsyncProducersPools int    `yaml:"async_producers_pools"`
}

type Kafka struct {
	configuration Configurations
}

func New(configuration Configurations) Kafka {
	return Kafka{configuration}
}
