package config

type Config struct {
	LogLevel    string
	ListenPort  int
	PrivKey     string
	ExpTokenSec int64
	BookClient  BookConfig
	Consumer    Consumer
}

type RoutingKey string

const (
	RoutingTest RoutingKey = "test_routing"
)

type BookConfig struct {
	Host string
}

type Consumer struct {
	Address     string
	Exchange    string
	QueueName   string
	RoutingKeys []RoutingKey
	Concurent   int
}
