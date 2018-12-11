Project Status: Development Stopped. Found project with similar goals and more support: [Draft](https://draft.sh)

# Kube Git Deploy
Automatically deploy git repositories to Kubernetes.

# Table Of Contents
- [Overview](#overview)
- [Design](#design)

# Overview
Automatically builds Docker images and deploys Kubernetes resources on Git 
pushes.

# Design
## Components

- [Etcd data backend](./api/README.md#data)
- [Golang API server](./api/README.md#endpoints)
- [Dashboard website](#./frontend/README.md)
- Golang CLI

## Behavior
API server creates GitHub web hooks for repositories.  

Receives GitHub webhook requests for every commit.  

Actions on commit depend on contents of `kube-git-deploy.toml` file in
repository. 

Golang CLI configures which GitHub repositories webhooks should be created for.
