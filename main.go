package main

import (
	"assessment/pkg/server"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"

	"flag"
	"fmt"
)

func main() {
	objectSource := flag.String("src", "", "endpoint url for object source, uses mock when not specified")
	useMockDB := flag.Bool("mock_db", false, "whether to use the postgres db or a in-memory mock")
	objectLifespan := flag.Int64("ol", 30, "how long an object will be persisted in seconds, defaults to 30")

	flag.Parse()

	err := godotenv.Load(".env")
	if err != nil {
		panic(fmt.Errorf("error reading env file: %s", err))
	}

	c := server.ServerConfig{
		ObjectSource:   *objectSource,
		UseMockDB:      *useMockDB,
		ObjectLifespan: time.Duration(*objectLifespan) * time.Second,
	}

	server.Serve(":9090", c)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}
