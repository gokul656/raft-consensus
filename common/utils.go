package common

import (
	"encoding/json"
	"log"
)

func ToByte(data interface{}) ([]byte, error) {
	log.SetPrefix("[ERR]")

	marhsalled, err := json.Marshal(data)
	if err != nil {
		return nil, UnableToUnmarshal
	}

	return marhsalled, nil
}

func FromByte[T interface{}](data []byte, result T) (T, error) {
	err := json.Unmarshal(data, &result)
	if err != nil {
		return result, UnableToMarshal
	}

	return result, nil
}
