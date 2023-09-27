package config

import (
	"os"
)

var (
	MongoDB_URI string
)

func GetEnvDefault(key string, defVal string) string {

	val, ex := os.LookupEnv(key)
	if !ex {
		val = defVal
	}
	return val

}

func SetEnvionment() {

	MongoDB_URI = GetEnvDefault("MONGODB_URI", "mongodb://localhost:27017")

}
