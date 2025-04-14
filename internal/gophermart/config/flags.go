package config

import (
	"flag"
	"github.com/caarlos0/env/v11"
	"github.com/spf13/cast"
)

var flags Flags

type Flags struct {
	Address       string `env:"RUN_ADDRESS"`
	DatabaseURI   string `env:"DATABASE_URI"`
	AccSystemAddr string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func NewConfig() error {
	flag.StringVar(&flags.Address, "a", "localhost:8080", "address and port to run")
	flag.StringVar(&flags.DatabaseURI, "d", "postgres://user:password@localhost:5432/gophermart?sslmode=disable", "database uri")
	flag.StringVar(&flags.AccSystemAddr, "r", "http://localhost:8090", "path to accrual system")
	flag.Parse()
	if err := env.Parse(&flags); err != nil {
		return err
	}
	return nil
}

func GetAddress() string { return cast.ToString(&flags.Address) }
func GetDatabaseURI() string {
	return cast.ToString(&flags.DatabaseURI)
}
func GetAccSystemAddr() string {
	return cast.ToString(&flags.AccSystemAddr)
}
