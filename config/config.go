package config

import (
	"fmt"
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

type HotEnv struct {
}

func StartConfig() error {
	fmt.Println("Config Started ...")
	v := viper.New()
	v.SetConfigType("json")
	v.SetConfigName("config.cold.json")
	v.AddConfigPath(".")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return err
	}
	err = v.Unmarshal(&Cold)
	if err != nil {
		return err
	}

	fmt.Println("Config generated, COLD value", Cold)
	fmt.Println("Config generated, HOT value", Hot)
	return nil
}
