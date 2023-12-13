package main

import (
	"encoding/json"
	"fmt"
	"log"
	"nats-service/app/cache"
	"nats-service/app/config"
	"nats-service/app/database"
	"nats-service/app/handler"

	"nats-service/app/model"
	"nats-service/app/service"
	cache_utils "nats-service/app/utils/cache"
	db_utils "nats-service/app/utils/sql"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
)

var (
	cacheMutex  sync.RWMutex
	cacheHolder *cache.Cache
)

const (
	clientID   = "subscriber"
	staticFile = "./static/index.html"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conf := config.New()
	// Cache
	cacheHolder = cache.New()

	// Load data to cache
	cache_utils.GetDataToCache(conf, cacheHolder)

	httpServer := handler.New(cacheHolder, &cacheMutex, staticFile)
	httpServer.StartHttpServer(
		conf.HttpConfig.Port,
		conf.HttpConfig.PatternServer,
		conf.HttpConfig.PatternStatic,
		conf.HttpConfig.StaticDir,
	)

	// Database
	database := database.New(
		conf.DatabaseConfig.User,
		conf.DatabaseConfig.Password,
		conf.DatabaseConfig.DatabaseName,
		conf.DatabaseConfig.Host,
		conf.DatabaseConfig.Port,
	)
	// db, err := database.Open(user, password, dbname, host, port)
	db, err := database.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Nats-Server connect
	nats := service.New(
		conf.NatsConfig.ConnectWait,
		conf.NatsConfig.PubAckWait,
		conf.NatsConfig.Interval,
		conf.NatsConfig.MaxOut,
		conf.NatsConfig.MaxPubAcksInflight,
	)
	nc, err := nats.Connect(conf.NatsConfig.Url, conf.NatsConfig.ClusterID, clientID)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Nats-Server Subscribe
	subject := "foo"
	sub, err := nats.Subscribe(nc, subject, func(msg *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(msg.Data))

		var order model.Order

		validate := validator.New()

		err := json.Unmarshal(msg.Data, &order)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)

			return
		}

		// Валидируем структуру

		if err = validate.Struct(order); err != nil {
			fmt.Println("Validation error:", err)
			return
		}

		cacheHolder.Set(cacheHolder.Len()+1, order)

		db_utils.CreateTableIfNotExist(db)
		db_utils.InsertToTable(db, msg.Data)

	})

	if err != nil {
		log.Println("Can't subscribe to nats: ")

	}

	log.Println("Nats-service is running. Press Ctrl+C to stop.")

	// Waiting stop signal
	<-stopChan

	sub.Unsubscribe()
	log.Println("Nats-service stopped.")
}
