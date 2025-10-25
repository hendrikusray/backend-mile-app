package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var (
	Cold ColdEnv
	Hot  HotEnv
)

type ColdEnv struct {
	DBMongoURI         string `json:"dBMongoURI" yaml:"dBMongoURI"`
	DBMongoMaxPoolSize uint64 `json:"dBMongoMaxPoolSize" yaml:"dBMongoMaxPoolSize"`
	DBMongoMinPoolSize uint64 `json:"dBMongoMinPoolSize" yaml:"dBMongoMaxPoolSize"`
	DBMongoUsername    string `json:"dBMongoUsername" yaml:"dBMongoUsername"`
	DBMongoPassword    string `json:"dBMongoPassword" yaml:"dBMongoPassword"`
}

type HotEnv struct{}

func StartConfig() error {
	fmt.Println("Config Started ...")
	v := viper.New()
	v.SetConfigType("json")
	v.SetConfigName("config.cold.json")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return err
	}
	if err := v.Unmarshal(&Cold); err != nil {
		return err
	}

	if uri := os.Getenv("MONGO_URI"); uri != "" {
		Cold.DBMongoURI = uri
	}

	fmt.Println("Config generated, COLD value", Cold)
	fmt.Println("Config generated, HOT value", Hot)
	return nil
}
