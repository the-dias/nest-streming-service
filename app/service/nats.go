package service

import (
	"database/sql"
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

type Nats struct {
	connectWait        time.Duration
	pubAckWait         time.Duration
	interval           int
	maxOut             int
	maxPubAcksInflight int
	// NATS_URL           string
	// clusterID          string
	// clientID           string
	conn stan.Conn
	db   *sql.DB
}

func New(connectWait, pubAckWait time.Duration, interval, maxOut, maxPubAcksInflight int) *Nats {
	nats := Nats{
		connectWait:        connectWait,
		pubAckWait:         pubAckWait,
		interval:           interval,
		maxOut:             maxOut,
		maxPubAcksInflight: maxPubAcksInflight,
	}

	return &nats
}

func (n *Nats) Connect(NATS_URL, clusterID, clientID string) (stan.Conn, error) {
	nc, err := stan.Connect(
		clusterID,
		clientID,
		stan.NatsURL(NATS_URL),
		stan.ConnectWait(n.connectWait),
		stan.PubAckWait(n.pubAckWait*10),
		stan.Pings(n.interval, n.maxOut),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("Connection lost, reason: %v", reason)
		}),
		stan.MaxPubAcksInflight(n.maxPubAcksInflight),
	)

	n.conn = nc

	if err != nil {
		return nil, err
	}

	return nc, nil
}

func (n *Nats) Subscribe(nc stan.Conn, subject string, cb stan.MsgHandler) (stan.Subscription, error) {
	sub, err := nc.Subscribe(subject, cb, stan.DurableName(subject))

	if err != nil {
		return nil, err
	}
	n.conn = nc
	return sub, nil
}

func (n *Nats) Publish(nc stan.Conn, subject string, data []byte) error {
	err := nc.Publish(subject, data)

	if err != nil {
		return err
	}
	n.conn = nc
	return nil
}

func (n *Nats) Close() {
	if n.conn != nil {
		if err := n.conn.Close(); err != nil {
			log.Println("Error closing database connection:", err)
		}
	}
}
