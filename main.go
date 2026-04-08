package main

import (
	"fmt"
	"log"

	"github.com/jasonsoprovich/rss-aggregator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Println("Config before:", cfg)

	err = cfg.SetUser("jason")
	if err != nil {
		log.Fatalf("error setting user: %v", err)
	}

	cfg, err = config.Read()
	if err != nil {
		log.Fatalf("error reading config again: %v", err)
	}
	fmt.Println("Config after:", cfg)
}
