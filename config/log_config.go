package config

import (
	"os"
)

func setupLogDir() {
	dest := GetEnv().TmpDir
	_, err := os.Stat(dest)
	if err != nil {
		if err := os.MkdirAll(dest, os.ModePerm); err != nil {
			panic(err)
		}
	}
}
