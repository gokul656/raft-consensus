package common

import (
	"errors"
	"log"
)

var (
	INVALID_LEADER   = errors.New("LEADER NOT FOUND")
	INVALID_RPC_PORT = errors.New("INVALID RPC PORT")
	INVALID_API_PORT = errors.New("INVALID API PORT")
	INVALID_LOG_PATH = errors.New("INVALID LOG PATH")
)

func HandlePanic() {
	if recover := recover(); recover != nil {
		log.Fatalln("[ERR]", recover)
	}
}
