# API
Kube Git Deploy API.

# Table Of Contents
- [Overview](#overview)
- [Development](#development)
	- [Configuration](#configuration)
	- [Dependencies](#dependencies)
	- [Local Etcd](#local-etcd)
	- [GitHub Application](#github-application)
- [User Manual](#user-manual)
	- [Repository Configuration File](#repository-configuration-file)
- [Endpoints](#endpoints)
- [Data](#data)

# Overview
Kube Git Deploy API.  

# Development
## Configuration
Configuration is passed via the following environment variables:  

- `PRIVATE_HTTP_PORT` (Optional, Default `5000`)
	- Port private API will respond to requests on
- `PUBLIC_HTTP_PORT` (Optional, Default `5001`)
	- Port public API will respond to requests on
- `PUBLIC_HTTP_HOST`
	- Full URI that public API can be reached at
	- Should include a schema, host and port (if needed)
	- Ex: `http://kube-git-deploy.example.com:5001`
- `PUBLIC_HTTP_SSL_ENABLED` (Optional, Default: `false`)
	- Indicates if the public API can be reached using SSL
- `ETCD_ENDPOINT` (Optional, Default `localhost:2379`)
	- URI of Etcd server
	- Should include host and port
- `GITHUB_CLIENT_ID`
	- GitHub API application client ID
- `GITHUB_CLIENT_SECRET`
	- GitHub API application client secret

## Dependencies
[Dep](https://github.com/golang/dep) is used to manage dependencies.

Install dependencies with:

```
dep ensure
```

## Local Etcd
Start a local Etcd server by running:

```
make etcd
```

## GitHub Application
Create a GitHub application with an authorization callback URL of: 

```
http://localhost:5000/api/v0/github/oauth_callback
```

# User Guide
## Repository Configuration File
### Structure Overview
Anytime code is pushed to a repository Kube Git Deploy will run a job.  

The behavior of jobs is defined by a file in the repository root named 
`kube-git-deploy.toml`.

A file contains units. Units are individual items which can be built
and deployed.

Each unit contains actions. Actions define what a job does for that unit. 

There are two types of actions: `docker` and `helm`.  

If a unit defines a Docker and Helm action, the Docker action will be
executed first.

### Action Definitions
[Docker](https://godoc.org/github.com/Noah-Huppert/kube-git-deploy/api/models#DockerActionConfig)  
[Helm](https://godoc.org/github.com/Noah-Huppert/kube-git-deploy/api/models#HelmActionConfig)

### Templating
Go templating can be used inside the file.  

[This data is available in Go templates](https://godoc.org/github.com/Noah-Huppert/kube-git-deploy/api/models#JobTarget).

### Syntax
Units are TOML sections. Actions are unit sub-sections. Action parameters are 
key value pairs.

### Example
#### Basic Example
Example file:

```toml
[api]
[[docker]]
directory = "./api"
tag = "noahhuppert/example-api:{{ .git.sha }}"

[[helm]]
chart = "./api/deploy"

[ui]
[[docker]]
directory = "./ui"
tag = "noahhuppert/example-ui:{{ .git.sha }}"

[[helm]]
chart = "./ui/deploy"
```

This will create 2 units named `api` and `ui`.  

The `api` unit will have 2 actions. The Docker action will build a Docker an
image in the `./api` directory and tag it with
`noahhuppert/example-api:{{ .git.sha }}`.

Notice how the Docker tag uses Go templating to get the Git commit's sha.  

The Helm action will deploy a Helm chart in `./api/deploy`.  

The `ui` unit is similar.

#### Template Example
Example file:

```toml
{{ if .git.branch == "master" }}
[api]
[[docker]]
directory = "."
tag = "noahhuppert/example-api:{{ .git.sha }}

[[helm]]
chart = "./deploy"
{{ end }}
```

This will define an `api` unit, which will only be deployed on the
master branch.

# Endpoints
The server provides a public and private API.  

The public API is accessible from the internet.  
The private API is only accessibly internally in Kubernetes.

## Get Tracked Repositories
GET `/api/v0/github/repositories/tracked`  

**API:** Private

**Actions:**

- Return a list of tracked GitHub repositories

**Request:** None

**Response:**

- `repositories` (Array[Repository])
	- Repository objects
- `ok` (Boolean)

## Track Repository
POST `/api/v0/github/repositories/:user/:repo`  

**API:** Private

**Actions:**

- Use the GitHub API to create web hook in repository
- Save repository as tracked in Etcd

**Request:**

- `:user` (String)
	- Repository GitHub user
- `:repo` (String)
	- Repository name

**Response:**

- `ok` (Boolean)

## Untrack Repository
DELETE `/api/v0/github/repositories/:user/:repo`  

**API:** Private

**Actions:**

- Use GitHub API to delete web hook in repository
- Delete repository in Etcd

**Request:**

- `:user` (String)
	- Repository GitHub user
- `:repo` (String)
	- Repository name

**Response:**

- `ok` (Boolean)

## OAuth Callback
GET `/api/v0/github/oauth_callback?code=:code`  

**API:** Private

**Actions:**

- Exchanges a temporary GitHub authentication code for a longer lived access 
	token
- Saves longer lived GitHub access token in etcd

**Request:** None

**Response:**

- `ok` (Boolean)

## Get GitHub Login URL
GET `/api/v0/github/login_url`  

**API:** Private

**Actions:**

- Returns the URL a user should visit to login with GitHub

**Request:** None

**Response:**

- `login_url` (String)
	- URL to send user to login
- `ok` (Boolean)

## Webhook
POST `/api/v0/github/repositories/:user/:repo/web_hook`  

**API:** Public

**Actions:**

- Triggers a build and deploy of the repository

**Request:**

- [GitHub Push Event](https://developer.github.com/v3/activity/events/types/#pushevent)
- `:user` (String)
	- Repository GitHub user
- `:repo` (String)
	- Repository name

**Response:**

- `ok` (Boolean)

## Health Check
GET `/healthz`

**API:** Public, Private

**Actions:** None

**Request:** None

**Response:**

- `ok` (Boolean)
- `server` (String)
	- Name of server, either `public` or `private`

# Data
Data is stored in Etcd.  

Some keys hold regular string values. While other keys hold serialized
JSON models.

Etcd stores data in a tree like a file system.

- `/github/auth` (Directory)
	- `/token` (String): Holds a user's GitHub access token
	- `/repositories/tracked/[USER]/[REPO]` (Directory)
		- `/information` ([Repository Model](https://godoc.org/github.com/Noah-Huppert/kube-git-deploy/api/models#Repository))
		- `/jobs/[ID]` ([Job Model](https://godoc.org/github.com/Noah-Huppert/kube-git-deploy/api/models#Job))
