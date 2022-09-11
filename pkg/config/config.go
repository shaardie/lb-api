package config

type Config struct {
	DBFilename string
	IP         *string
	Hostname   *string
}

var Cfg Config
