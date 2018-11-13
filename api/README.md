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

# Overview
Kube Git Deploy API.  

# Development
## Configuration
Configuration is passed via the following environment variables:  

- `HTTP_PORT` (Optional, Default `5000`)
- `ETCD_ENDPOINT` (Optional, Default `localhost:2379`)
- `GITHUB_CLIENT_ID`
- `GITHUB_CLIENT_SECRET`

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
## GitHub
### Get Repositories
GET `/api/v0/github/repositories`  

Request: None

Response:

- `repositories` (Array[Map])
	- `tracking` (Boolean)
		- Indicates if server is automatically deploying repository
	- `name` (String)
		- Slug
- `ok` (Boolean)

### Track Repository
POST `/api/v0/github/repositories/:name/track`  

Request: None

Response:

- `ok` (Boolean)

### Untrack Repository
POST `/api/v0/github/repositories/:name/untrack`  

Request: None

Response:

- `ok` (Boolean)

### OAuth Callback
GET `/api/v0/github/oauth_callback?code=:code`  

Request: None

Response: 

- `ok` (Boolean)

### Get GitHub Login URL
GET `/api/v0/github/login_url`  

Request: None

Response:

- `login_url` (String)
	- URL to send user to login
- `ok` (Boolean)

### Webhook
POST `/api/v0/github/repositories/:name/webhook`  

Request: [GitHub Push Event](https://developer.github.com/v3/activity/events/types/#pushevent)

Response: 

- `ok` (Boolean)
