package config

import (
	"io"
	"log"
	"os"
)

func setupLogDir() {
	dest := GetEnv().LogDir
	_, err := os.Stat(dest)
	if err != nil {
		if err := os.MkdirAll(dest, os.ModePerm); err != nil {
			panic(err)
		}
	}
}

func EnableLogging(filename string) {
	logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	// multi writer for writing logs both to console & file
	mw := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(mw)
}

func DisableLoggin() {
	log.SetOutput(os.Stdout)
}
