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

	// Run HTTP servers
	serverReturns := make(chan string)

	go func() {
		// Private
		logger.Info("Starting private HTTP server")

		privServer := server.NewPrivateServer(ctx, logger, cfg, etcdKV)

		err = privServer.Run()
		if err != nil {
			logger.Fatalf("error while running private HTTP server: %s", err.Error())
		}

		serverReturns <- "private"
	}()

	go func() {
		// Public
		logger.Info("Starting public HTTP server")

		pubServer := server.NewPublicServer(ctx, logger, cfg, etcdKV)

		err = pubServer.Run()
		if err != nil {
			logger.Fatalf("error while running public HTTP server: %s", err.Error())
		}

		serverReturns <- "public"
	}()

	// Wait for both servers to stop
	for i := 0; i < 2; i++ {
		logger.Infof("%s server stopped", <-serverReturns)
	}
}
