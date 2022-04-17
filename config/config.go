package config

import (
	"bufio"
	"cft/log"
	"encoding/json"
	"os"
)

type Config struct {
	CheckpointDir string     `json:"checkpoint_dir"`
	RpcConfig     RpcConfig  `json:"rpc_config"`
	HttpConfig    HttpConfig `json:"http_config"`
}

type RpcConfig struct {
	Port string `json:"port"`
}

type HttpConfig struct {
	Port string `json:"port"`
}

func GetConfig() *Config {
	return cfg
}

var cfg *Config = nil

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err = decoder.Decode(&cfg); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return cfg, nil
}
