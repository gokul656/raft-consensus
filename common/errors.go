package common

import (
	"errors"
	"log"
)

var (
	InvalidLeader  = errors.New("leader not found")
	InvalidRPCPort = errors.New("invalid RPC port")
	InvalidAPIPort = errors.New("invalid API port")
	InvalidLogPath = errors.New("invalid Log path")

	UnableToUnmarshal = errors.New("unable to convert struct to byte[]")
	UnableToMarshal   = errors.New("unable to byte[] struct to struct")
)

func HandlePanic() {
	if recover := recover(); recover != nil {
		log.Fatalln("[ERR]", recover)
	}
}
