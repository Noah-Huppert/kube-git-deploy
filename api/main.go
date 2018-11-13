package main

import (
	"context"
	"time"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"
	"github.com/Noah-Huppert/kube-git-deploy/api/libetcd"
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

	etcdKV := etcd.NewKeysAPI(etcdClient)

	_, err = etcdKV.Set(ctx, "/ping", "pong", nil)
	if err != nil {
		logger.Fatalf("error testing Etcd connection: %s", err.Error())
	}

	// Create initial keys
	_, err = etcdKV.Set(ctx, libetcd.KeyTrackedGitHubRepos, "",
		&etcd.SetOptions{
			Dir:       true,
			PrevExist: etcd.PrevNoExist,
		})
	if err != nil && err.(etcd.Error).Code != etcd.ErrorCodeNodeExist {
		logger.Fatalf("error creating initial empty tracked GitHub "+
			"repositories key: %s", err.Error())
	}

	// Run HTTP server
	logger.Info("Starting HTTP server")

	srv := server.NewServer(ctx, logger, cfg, etcdKV)

	err = srv.Run()
	if err != nil {
		logger.Fatalf("error starting HTTP server: %s", err.Error())
	}
}
