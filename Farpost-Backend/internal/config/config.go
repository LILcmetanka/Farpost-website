package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env 			string 		`yaml:"env" env-default:"local"`
	StoragePath		string 		`yaml:"storage_path" env-required:"true"`
	HTTPServer					`yaml:"http_server"`
}

type HTTPServer struct {
	Address			string 			`yaml:"address" env-default:"localhost:1234"`
	Timeout			time.Duration	`yaml:"timeout" env-default:"4s"`
	IddleTimeout	time.Duration	`yaml:"iddle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	
	if path == "" {
		panic("config path is empty!")
	} 

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist: " + path)
	}

	var cfg Config

	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		panic("failed to read config" + err.Error())
	} 

	return &cfg
}

func fetchConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "path to config path")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}