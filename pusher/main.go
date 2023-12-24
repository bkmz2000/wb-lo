package main

import (
	"fmt"
	"log"
	"time"
	"wb-l0/common"

	"github.com/nats-io/nats.go"

	_ "github.com/lib/pq"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}

	defer nc.Close()

	for {
		object := common.GenerateRandomOrderJSON()
		nc.Publish("insert", []byte(object))
		fmt.Println(object)

		for i := 0; i < 3; i++ {
			time.Sleep(time.Second)
			fmt.Println(i)
		}

		fmt.Println("\n\n\n\n")
	}
}
