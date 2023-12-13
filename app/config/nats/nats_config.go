package nats

import "time"

type NatsConfig struct {
	ConnectWait        time.Duration
	PubAckWait         time.Duration
	Interval           int
	MaxOut             int
	MaxPubAcksInflight int
	Url                string
	ClusterID          string
}
