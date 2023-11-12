package config

type Config struct {
	Server      Server
	Client      Client
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

type Server struct {
	Host string
	Port int
}

type Client struct {
	Host string
	Port int
}

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
