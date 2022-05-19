package main

import "github.com/kelseyhightower/envconfig"

type config struct {
	RedisURL string `envconfig:"REDIS_URL" default:"localhost:6379"`
}

func newConfig() *config {
	var conf config

	err := envconfig.Process("", &conf)
	if err != nil {
		panic(err)
	}

	return &conf
}
