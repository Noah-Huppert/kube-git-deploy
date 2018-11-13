package main

import (
	"context"
	"time"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"
	"github.com/Noah-Huppert/kube-git-deploy/api/server"

	"github.com/Noah-Huppert/golog"
	etcd "go.etcd.io/etcd/client"
)

func main() {
	// Get context
	ctx := context.Background()

	// Setup logger
	logger := golog.NewStdLogger("api")

	// Load configuration
	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatalf("error loading configuration: %s", err.Error())
	}

	// Connect to Ectd
	etcdClient, err := etcd.New(etcd.Config{
		Endpoints:               []string{cfg.EtcdEndpoint},
		Transport:               etcd.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	})
	if err != nil {
		logger.Fatalf("error connecting to Etcd: %s", err.Error())
	}

	//etcdKV := etcdClient.NewKeysAPI(etcdClient)
	_ = etcd.NewKeysAPI(etcdClient)

	// Run HTTP server
	logger.Info("Starting HTTP server")

	err = server.RunServer(ctx, cfg)
	if err != nil {
		logger.Fatalf("error starting HTTP server: %s", err.Error())
	}
}
