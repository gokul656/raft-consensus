package common

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

const (
	FollowerHealthCheckDelay = 2 * time.Second
	LeaderHealthCheckDelay   = 2 * time.Second
)

var startupTime = time.Now().UnixMilli()

func ToByte(data interface{}) ([]byte, error) {
	log.SetPrefix("[ERR]")

	marhsalled, err := json.Marshal(data)
	if err != nil {
		return nil, ErrUnableToUnmarshal
	}

	return marhsalled, nil
}

func FromByte[T interface{}](data []byte, result T) (T, error) {
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, ErrUnableToMarshal
	}

	return result, nil
}

func GetLogfileName(dir string) string {
	return fmt.Sprintf("%s/raft-%d.log", dir, startupTime)
}
