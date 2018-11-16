package main

import (
	"context"
	"time"

	"github.com/Noah-Huppert/kube-git-deploy/api/config"
	"github.com/Noah-Huppert/kube-git-deploy/api/jobs"
	"github.com/Noah-Huppert/kube-git-deploy/api/models"
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
	_, err = etcdKV.Set(ctx, models.KeyDirRepositories, "",
		&etcd.SetOptions{
			Dir:       true,
			PrevExist: etcd.PrevNoExist,
		})
	if err != nil && err.(etcd.Error).Code != etcd.ErrorCodeNodeExist {
		logger.Fatalf("error creating initial empty tracked GitHub "+
			"repositories key: %s", err.Error())
	}

	// Create JobRunner
	jobRunner := jobs.NewJobRunner(ctx, logger.GetChild("job_runner"),
		etcdKV)

	go func() {
		logger.Info("Starting job runner")

		jobRunner.Run()
	}()

	// Run HTTP servers
	serverReturns := make(chan string)

	go func() {
		// Private
		logger.Infof("Starting private HTTP server on :%d",
			cfg.PrivateHTTPPort)

		privServer := server.NewPrivateServer(ctx, logger, cfg, etcdKV)

		err = privServer.Run()
		if err != nil {
			logger.Fatalf("error while running private HTTP server: %s", err.Error())
		}

		serverReturns <- "private"
	}()

	go func() {
		// Public
		logger.Infof("Starting public HTTP server on %s",
			cfg.PublicHTTPHost)

		pubServer := server.NewPublicServer(ctx, logger, cfg, etcdKV,
			jobRunner)

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
