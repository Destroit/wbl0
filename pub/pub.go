package main

import (
	stan "github.com/nats-io/stan.go"
	"log"
	"os"
)

func main() {
	sc, err := stan.Connect("test-cluster", "test-server")
	if err != nil {
		log.Fatal(err)
	}

	files := []string{"orders.json", "orders2.json", "not.json", "invalid.json"}

	for _, s := range files {
		jsonFile, err := os.ReadFile(s)
		if err != nil {
			log.Fatal(err)
		}
		sc.Publish("orders", jsonFile)
	}
}
