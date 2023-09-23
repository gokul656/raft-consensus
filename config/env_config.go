package config

import (
	"flag"

	"github.com/gokul656/raft-consensus/common"
)

var rpcPortFlag = flag.String("grpc-port", "51505", "provide port for gRPC connection")
var apiPortFlag = flag.String("api-port", "3000", "provide port for API connection")
var tmpDirFlag = flag.String("raft-dir", "./tmp/raft-logs", "provide directory to be used for storing logs")

type EnvConfig struct {
	APIPort string
	RPCPort string
	TmpDir  string
}

func (cfg *EnvConfig) LoadConfig() {
	flag.Parse()

	cfg.RPCPort = *rpcPortFlag
	cfg.APIPort = *apiPortFlag
	cfg.TmpDir = *tmpDirFlag
}

var env *EnvConfig

func init() {
	defer common.HandlePanic()

	env = &EnvConfig{}
	env.LoadConfig()

	setupLogDir() // creates a log dir if not exists
}

func GetEnv() EnvConfig {
	return *env
}
