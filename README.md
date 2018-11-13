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

- Golang API server
- Etcd data backend
- Golang CLI

## Behavior
API server creates GitHub webhooks for repositories. Has no API authentication.  

Receives GitHub webhook requests for every commit.  

Actions on commit depend on contents of `kube-git-deploy.toml` file in
repository. Actions will be executed in order: Docker -> Helm

Golang CLI configures which GitHub repositories webhooks should be created for.

## Configuration File
`kube-git-deploy.toml`.  

Can contain any number of sections.  

Each section has the following fields:

- `helm_chart` (Optional, String)
    - Helm chart to deploy
	- If chart is located in repository value should be in form: `repository@chart`
	- Otherwise should be local path to chart
- `docker` (Optional, Map)
	- `tag` (Required, String)
		- Docker image tag


Go templates can be used in this file. The following variables will
be available:

- `git` (Map)
	- `branch` (String)
		- Commit branch
	- `sha` (String)
		- Commit SHA
