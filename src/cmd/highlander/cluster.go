package main

import (
	"strconv"
	"strings"

	"github.com/hashicorp/memberlist"
	"github.com/pkg/errors"
)

func clusterJoin() (*memberlist.Memberlist, error) {
	var err error

	cfg := memberlist.DefaultLocalConfig()

	cfg.Name = cmd.instanceID

	if len(cmd.clusterBind) > 0 {
		b := strings.Split(cmd.clusterBind, ":")
		cfg.BindAddr = b[0]
		cfg.BindPort, err = strconv.Atoi(b[1])
	}

	clusterNodes, err := memberlist.Create(cfg)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create empty cluster memberlist: ")
	}

	// Join an existing cluster by specifying at least one known member.
	if len(cmd.clusterConnect) > 0 {
		_, err = clusterNodes.Join(strings.Split(cmd.clusterConnect, ","))
		if err != nil {
			return nil, errors.Wrapf(err, "Failed to join cluster: ")
		}
	}
	return clusterNodes, nil
}
