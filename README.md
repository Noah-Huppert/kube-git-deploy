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

- Etcd data backend
- Golang API server
- Golang CLI

## Behavior
API server creates GitHub webhooks for repositories. Has no API authentication.  

Receives GitHub webhook requests for every commit.  

Actions on commit depend on contents of `kube-git-deploy.toml` file in
repository. Actions will be executed in order: Docker -> Helm

Golang CLI configures which GitHub repositories webhooks should be created for.
