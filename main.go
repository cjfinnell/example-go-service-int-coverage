package main

import "log"

func main() {
	conf := newConfig()

	srv := newServer(conf)

	if err := srv.run(); err != nil {
		log.Fatalf("server exited with error: %s", err)
	}
}
