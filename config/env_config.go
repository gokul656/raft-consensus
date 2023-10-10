package config

import (
	"flag"
)

var instaceIDFlag = flag.String("name", "peer-1", "provide peer name")
var leaderIDFlag = flag.String("leader-name", "peer-0", "provide leader name")
var rpcPortFlag = flag.String("grpc-port", "51500", "provide port for gRPC connection")
var apiPortFlag = flag.String("api-port", "3000", "provide port for API connection")
var knownLeader = flag.String("leader", "", "provide leader address")
var logDirFlag = flag.String("raft-dir", "./tmp/raft-logs", "provide directory to be used for storing logs")

type EnvConfig struct {
	InstanceID string
	APIPort    string
	RPCPort    string
	LogDir     string
	Leader     string
	LeaderID   string
}

func (cfg *EnvConfig) LoadConfig() {
	flag.Parse()

	cfg.InstanceID = *instaceIDFlag
	cfg.RPCPort = *rpcPortFlag
	cfg.APIPort = *apiPortFlag
	cfg.LogDir = *logDirFlag
	cfg.Leader = *knownLeader
	cfg.LeaderID = *leaderIDFlag
}

var env *EnvConfig

func init() {
	env = &EnvConfig{}
	env.LoadConfig()

	setupLogDir() // creates a log dir if not exists
}

func GetEnv() EnvConfig {
	return *env
}
