package common

import (
	"errors"
	"log"
)

type CustomError struct {
	Err error
	Msg string
}

var (
	ErrInvalidLeader  = errors.New("leader not found")
	ErrInvalidPeer    = errors.New("invalid peer")
	ErrInvalidRPCPort = errors.New("invalid RPC port")
	ErrInvalidAPIPort = errors.New("invalid API port")
	ErrInvalidLogPath = errors.New("invalid Log path")

	ErrPeerUnavailable = errors.New("unable to establish connection with peer")

	ErrUnableToUnmarshal = errors.New("unable to convert struct to byte[]")
	ErrUnableToMarshal   = errors.New("unable to byte[] struct to struct")
)

func HandlePanic(msg string) {
	if recover := recover(); recover != nil {
		switch err := recover.(type) {
		case CustomError:
			switch err.Err {
			case ErrPeerUnavailable:
				log.Println(msg, "peer unavailable", err.Msg)
			default:
				log.Println(msg, recover)
			}
		default:
			log.Println(err)
		}

	}
}
