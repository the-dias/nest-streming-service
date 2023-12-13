package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"nats-service/app/model"
	"nats-service/app/service"
	"os"
	"strconv"
	"time"
)

const (
	connectWait        = time.Second * 30
	pubAckWait         = time.Second * 30
	interval           = 10
	maxOut             = 5
	maxPubAcksInflight = 25

	host     = "localhost"
	port     = 5438
	user     = "dias"
	password = "dias2502"
	dbname   = "orders"

	NATS_URL  = "nats://0.0.0.0:4222"
	clusterID = "test-cluster"
	clientID  = "publisher"
)

func main() {
	nats := service.New(connectWait, pubAckWait, interval, maxOut, maxPubAcksInflight)

	nc, err := nats.Connect(NATS_URL, clusterID, clientID)

	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	fileName := "model.json"
	jsonFile, err := os.Open(fileName)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully published model.json")

	defer jsonFile.Close()
	byteJson, err := io.ReadAll(jsonFile)

	var order model.Order
	if err := json.Unmarshal(byteJson, &order); err != nil {
		log.Printf("Error decoding JSON: %v", err)

		return
	}

	order.OrderUID += strconv.Itoa(rand.Intn(100))

	if err != nil {
		fmt.Println(err)
	}
	subject := "foo"

	byteObject, err := json.Marshal(order)

	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)

		return
	}

	nats.Publish(nc, subject, byteObject)
}
