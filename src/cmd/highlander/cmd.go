package main

import (
	"flag"
	"time"

	"github.com/segmentio/ksuid"
)

type Cmd struct {
	bind           string
	remote         string
	healthTick     time.Duration
	healthExpiry   time.Duration
	instanceID     string
	clusterConnect string
	clusterBind    string
}

func parseCmd() Cmd {
	var cmd Cmd
	// L7 router settings
	flag.DurationVar(&cmd.healthTick, "t", 1*time.Second, "health check every n second")
	flag.DurationVar(&cmd.healthExpiry, "e", 5*time.Second, "disable stream after n seconds of inactivity")

	// service
	flag.StringVar(&cmd.bind, "l", "0.0.0.0:9091", "listen on ip:port")
	flag.StringVar(&cmd.remote, "r", "http://localhost:9090", "reverse proxy addr")
	flag.StringVar(&cmd.instanceID, "i", "", "InstanceID is a uuid, used among other things, for clustering")

	// clustering
	flag.StringVar(&cmd.clusterConnect, "c", "", "comma separated list of ip:port combinations that can be use for joining the cluster. Not joining if empty string")
	flag.StringVar(&cmd.clusterBind, "b", "0.0.0.0:6746", "0.0.0.0:6746 formatted cluster bind")

	flag.Parse()
	if cmd.instanceID == "" {
		cmd.instanceID = ksuid.New().String()
	}
	return cmd
}
