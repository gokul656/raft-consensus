package config

import (
	"flag"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var configFile = flag.String("config", "", "path for config.yml")

type Config struct {
	Server ServerConfig `yaml:"server"`
}

type PeerConfig struct {
	FallbackUrl string   `yaml:"fallback_peer_url"`
	Timeout     uint64   `yaml:"health_check_timeout"`
	Hosts       []string `yaml:"hosts"`
}

type ServerConfig struct {
	Peer       PeerConfig `yaml:"peer"`
	InstanceID string     `yaml:"instance_id,omitempty"`
	APIPort    string     `yaml:"api_port,omitempty"`
	RPCPort    string     `yaml:"rpc_port,omitempty"`
	LogDir     string     `yaml:"log_dir,omitempty"`
	Leader     string     `yaml:"leader,omitempty"`
	LeaderID   string     `yaml:"leader_id,omitempty"`
}

func (cfg *Config) LoadConfig() {
	flag.Parse()
	if *configFile == "" {
		log.Fatalln("unble to locate config.yml file")
	}

	file, err := os.ReadFile(*configFile)
	if err != nil {
		log.Fatalln(err.Error())
	}

	yaml.Unmarshal(file, env)
}

var env *Config

func init() {
	env = &Config{}
	env.LoadConfig()

	setupLogDir() // creates a log dir if not exists
}

func GetEnv() ServerConfig {
	return env.Server
}
