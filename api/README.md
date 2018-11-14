# API
Kube Git Deploy API.

# Table Of Contents
- [Overview](#overview)
- [Development](#development)
	- [Configuration](#configuration)
	- [Dependencies](#dependencies)
	- [Local Etcd](#local-etcd)
	- [GitHub Application](#github-application)
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

Etcd stores data in a directory like structure.  

In the diagram below a key is a directory if it has a **D:** before its key.  
A key is a value file if it has a **V:** before its key.

- **D:** `/github`
	- **D:** `/auth`
		- **V:** `/token`: Jolds a user's GitHub access token
	- **D:** `/repositories/tracked`
		- **D:** `/[USER]/[REPO]`: Each tracked GitHub repository has its own
			directory
			- **V:** `/name`: Name of tracked GitHub repository, in standard 
				`username/repo-name` GitHub format
			- **V:** `/web_hook_id`: ID of GitHub repository web hook
			- **D:** `/jobs`
				- **D:** `/[ID]`: Each job has its own directory
					- **V:** `/status`: Current status of job, values are:
						- `waiting`: Initiated but not started
						- `running`: Running
						- `done`: Successfully completed
						- `error`: Completed but failed
					- **D:** `/output`: Holds output for each module in job
						- **D:** `/[XXX]_[MODULE]`: Each module gets its
							own directory
							- **D:** `/{docker,helm}`
								- **V:** `/status`: Status of build, values
									same as above
								- **V:** `/output`: Build output

## GitHub Access Token
**Key:** `/github/auth/token`  
**Dir:** False

Holds a user's GitHub access token.

## Tracked GitHub Repository
**Key:** `/github/repositories/tracked/[USER]/[REPO]`  
**Dir:** True

Tracked GitHub repositories each have their own directory.  

Data related to this tracked GitHub repository is stored in nodes inside of 
this directory.

### Name
**Key:** `/github/repoositories/tracked/[USER]/[REPO]/name`  
**Dir:** False

Holds the name of the tracked GitHub repository.

### Web Hook ID
**Key:** `/github/repositories/tracked/[USER]/[REPO]/web_hook_id`  
**Dir:** False

Holds the ID of the GitHub repository web hook.

### Jobs
**Key:** `/github/repositories/tracked/[USER]/[REPO]/jobs`  
**Dir:** True

Holds information about build and deploy jobs for a repository.  

Each job will have a unique ID.  

A job is created when a GitHub repository web hook is sent to the server.

#### Job Status
**Key:** `/github/repositories/tracked/[USER]/[REPO]/jobs/[ID]/status`  
**Dir:** False

Holds the status of the job.  

Valid values:

- `triggered`
	- Indicates a job has been triggered but nothing has started
- `running`
	- Indicates a job is currently running
- `done`
	- Indicates a job has completed
- `error`
	- Indicates a job completed but failed

#### Job Output
**Key:** `/github/repositories/tracked/[USER]/[REPO]/jobs/[ID]/output`  
**Dir:** True

Holds job output.  

Inside the output directory can be any number of sub directories. Each sub 
directory represents module.  

##### Job Module Output
**Key:** `/github/repositories/tracked/[USER]/[REPO]/jobs/[ID]/output/xxx_[MODULE]`  
**Dir:** True

Inside each module the following files are present:

- `docker`
- `helm`

They contain all relevant built output for each step.
