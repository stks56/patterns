package options

import (
	"errors"
	"log"
	"net/http"
	"strconv"
)

// spec:
// If port is not set, use default port
// if port is zero, use random port
// If port is negative, print error
// If port is positive, use that port

// procedural pattern
// Level: Poor
func main() {
	port := 8080
	s, err := NewServer("localhost", &port)
	if err != nil {
		log.Println(err)
	}

	s.ListenAndServe()
}

func NewServer(addr string, port *int) (*http.Server, error) {
	if *port < 0 {
		return nil, errors.New("port cannot be negative")
	}
	if port == nil {
		// use default port
	}
	if *port == 0 {
		// use random port
	}

	return &http.Server{
		Addr: addr + ":" + strconv.Itoa(*port),
	}, nil
}

// config struct pattern
// Level: Average
func main() {
	port := 8080
	c := Config{
		Port: &port,
	}
	s, err := NewServer("localhost", &c)
	if err != nil {
		log.Println(err)
	}

	s.ListenAndServe()
}

type Config struct {
	// integer pointer by distinction nil or 0.
	Port *int
}

func NewServer(addr string, cfg *Config) (*http.Server, error) {
	if cfg.Port == nil {
		// use random port
	}
	if *cfg.Port < 0 {
		return nil, errors.New("port cannot be negative")
	}
	if *cfg.Port == 0 {
		// use default port
	}

	return &http.Server{
		Addr: addr + ":" + strconv.Itoa(cfg.Port),
	}, nil
}

// builder pattern
// Level: Good
// cons: Delayed validation, port method can not return error, must assign empty config struct when use default option
func main() {
	builder := ConfigBuilder{}
	builder.Port(8080) // usable method chain
	cfg, err := builder.build()
	if err != nil {
		log.Println(err)
	}

	s, err := NewServer("localhost", cfg)
	if err != nil {
		log.Println(err)
	}

	s.ListenAndServe()
}

type Config struct {
	Post int
}

type ConfigBuilder struct {
	port *int
}

func (b *ConfigBuilder) Port(port int) *ConfigBuilder {
	// Can also write port init logic here.
	b.port = &port
	return b
}

func (b *ConfigBuilder) build() (Config, error) {
	cfg := Config{}
	if b.port == nil {
		// use random port
	}
	if *b.port < 0 {
		return Config{}, errors.New("port cannot be negative")
	}
	if *b.port == 0 {
		// use default port
	}
	cfg.Port = b.port

	return cfg, nil
}

func NewServer(addr string, cfg *Config) (*http.Server, error) {
	return &http.Server{
		Addr: addr + ":" + strconv.Itoa(*cfg.Port),
	}, nil
}

// functional options pattern
// pros: immediate validation eval, lightweight writing, readable, Encapsulation
func main() {
	port := 8080
	// can write default options using like this:
	// s, err := NewServer("localhost")
	s, err := NewServer("localhost", WithPort(port))
	if err != nil {
		log.Println(err)
	}

	s.ListenAndServe()
}

type options struct {
	port *int
}

type Option func(options *options) error

func WithPort(port int) Option {
	return func(options *options) error {
		if port < 0 {
			return errors.New("port cannot be negative")
		}

		options.port = &port
		return nil
	}
}

func NewServer(addr string, opts ...Option) (*http.Server, error) {
	var options options
	for _, opt := range opts {
		err := opt(&options)
		if err != nil {
			return nil, err
		}
	}

	var port int
	if options.port == nil {
		// use default port
	}
	if *options.port == 0 {
		// use random port
	} else {
		port = *options.port
	}

	return &http.Server{Addr: addr + ":" + strconv.Itoa(port)}, nil
}
