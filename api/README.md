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

A file contains modules. Modules are individual items which can be built
and deployed.

Each module contains steps. Steps define the actions a job performs. 

There are two types of steps: `docker` and `helm`.  

If a module defines a Docker and Helm step, the Docker step will be
executed first.

### Step Definitions
**Docker:**
- `directory` (String): Directory Dockerfile is located in
- `tag` (String): Docker image tag 

**Helm:**  
- `chart` (String): Helm chart to deploy
	- If chart is located in repository value should be in form: `repository@chart`
	- Otherwise should be local path to chart

### Templating
Go templating can be used inside the file. The following data can be access:

 - `git` (Map)
	- `branch` (String)
		- Commit branch
	- `sha` (String)
		- Commit SHA


### Syntax
Modules are TOML sections. Steps are module sub-sections. Step parameters are 
key value pairs.

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

- `repositories` (Array[String])
	- Repository names
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

# Data
Data is stored in Etcd.  

Some keys hold regular string values. While other keys hold serialized
JSON models.

## Key Structure
Etcd stores data in a tree like a file system.

- `/github/auth` (Directory)
	- `/token` (String): Holds a user's GitHub access token
	- `/repositories/tracked/[USER]/[REPO]` (Directory)
		- `/information` ([Repository Model](#repository-model))
		- `/jobs/[ID]` ([Job Model](#job-model))

## Repository Model
Tracked GitHub repository information.  

Fields:

- `owner` (String): GitHub repository owner
- `name` (String): GitHub repository name
- `web_hook_id` (Integer): ID of GitHub web socket

## Job Model
Repository build and deploy job.

See [Repository Configuration File](#repository-configuration-file) for more 
details on the structure of ajob.

Fields:

- `id` (Integer): ID of job
- `modules` (Array[Job Module]): Modules in repository
	- `configuration` (Object): Raw module configuration
		- `docker` (Object): Docker configuration with keys from
			[Step Definitions](#step-definitions)
		- `helm` (Object): Helm configuration with keys from 
			[Step Definitions](#step-definitions)
	- `state` (Object): Steps in module
		- `{docker,helm}`
			- `status` (String): Status of step, allowed values:
				- `waiting`: Initiated but not started
				- `running`: Running
				- `success`: Successfully completed
				- `error`: Completed but failed
			- `output` (String): Raw output of step
- `metadata` (Object): Information about event which triggered job
	- `owner` (String): GitHub repository owner
	- `name` (String): GitHub repository name
	- `branch` (String): Name of branch of commit which triggered job
	- `commit_sha` (String): Git commit sha which triggered job
