# 🛰 drones
REST API that allows clients to communicate with drones (i.e. **dispatch controller**).

> **NOTE**: Drones app has been tested on **Ubuntu 18.04** and on **Windows 10 with WSL** and Golang 1.16 was used.

## Table of Contents

- [API specification](#api_spec)
- [Assigned tasks](#assigned_task)
- [Configuration file](#config_file)
- [Get Started](#get_started)
  * [Deployment ways (2 ways)](#deploy_ways)
    - [Docker way](#docker_way)
    - [Manual way](#manual_way)
- [Tech and packages](#tech)
- [Architecture](#arch)
## ⚙️API specification <a name="api_spec"></a>

The **Drone API server** provides the following API with communicating the **DB**:

| Tag           | Title                              | URL                                      | Query | Method |
| ------------- | ---------------------------------- | ---------------------------------------- | ----- | ---- |
| Auth          | user authentication (Using JWT)    | `/api/v1/auth`                           |   -   |`POST`|
| Auth          | user logout                        | `/api/v1/auth/logout`                    |   -   |`GET` |
| Auth          | get user authenticated             | `/api/v1/auth/user`                      |   -   |`GET` |
| Database      | Populate DB with fake data         | `/api/v1/database/populate`              |   -   |`POST`|
| Drones        | Get all drones or filters for State| `/api/v1/drones`                         |?state=|`GET` |
| Drones        | Registers or update a drone        | `/api/v1/drones`                         |   -   |`POST`|
| Drones        | Get a drone by serialNumber        | `/api/v1/drones/:serialNumber`           |   -   |`GET` |
| Logs          | Get event logs                     | `/api/v1/logs`                           |   -   |`GET` |
| Medications   | Get medications                    | `/api/v1/medications`                    |   -   |`GET` |
| Medications   | Checking loaded items for a drone  | `/api/v1/medications/items/:serialNumber`|   -   |`GET` |
| Medications   | Load a drone with medication items | `/api/v1/medications/items/:serialNumber`|   -   |`POST`|

To see the API specifications in more detail, run the app and visit the swagger docs:

> http://localhost:7001/swagger/index.html

![swagger ui](/docs/images/swagger-ui.png)

## 📝 Assigned tasks <a name="assigned_task"></a>
|  Done          | Task       | Endpoint                              |
| -------------- | -----------|------------------------- |
| ✅ | registering a drone;                                | 👉🏾 endpoint: `/api/v1/drones  [POST]`
| ✅ | loading a drone with medication items;              | 👉🏾 endpoint: `/api/v1/medicationsitems/:serialNumber [POST]`
| ✅ | checking loaded medication items for a given drone; | 👉🏾 endpoint: `/api/v1/medicationsitems/:serialNumber [GET]`
| ✅ | checking available drones for loading;              | 👉🏾 endpoint: `/api/v1/drones?state=1 [GET]`
| ✅ | check drone battery level for a given drone;        | 👉🏾 endpoint: `/api/v1/drones/:serialNumber [GET], Get a drone by serialNumber`

> The endpoints `/api/v1/drones  [POST]` and `/api/v1/medicationsitems/:serialNumber [POST]` can also be used to update.

| Done | Functional and Non-functional requirements |
| -------------- | -----------|
| ✅ | periodic task to check drones battery levels and create event log;
| ✅ | prevent the drone from being loaded with more weight that it can carry;
| ✅ | prevent the drone from being in LOADING state if the battery level is **below 25%**;
| ✅ | Your project must be buildable and runnable;
| ✅ | Your project must have a README file with build/run/test instructions (use DB file);
| ✅ | Required data must be preloaded in the database.
| ✅ | a bit of unit and end-to-end testing
| ✅ | show us how you work through your commit history.



## 🛠️️ Configuration file (conf.yaml) <a name="config_file"></a>
👉🏾 ![The config file](/conf/conf.yaml)

|  Param      | Description       | default value   |
| ----------- | -----------|------------------------- |
| APIDocIP    | IP to expose the api (unused)  | 127.0.0.1
| DappPort    | app PORT              | 7001
| StoreDBPath | DB file location      | ./db/data.db
| CronEnabled | active the cron job   | true
| LogDBPath   | DB file event logs    | ./db/event_log.db
| EveryTime   | time interval (in seconds) that the cron task is executed | 300 seconds (every 5 minutes)

By default, **StoreDBPath** generates the database file in the /db folder at the root of the project.

The server exposes the `/api/v1/database/populate` POST endpoint to generate and repopulate the database whenever necessary.
## ⚡ Get Started <a name="get_started"></a>

Download the drones.restapi project and move to root of project:
```bash
git clone https://github.com/kmilodenisglez/drones.restapi.git && cd drones.restapi 
```

### 🚀 Deployment ways (2 ways)  <a name="deploy_ways"></a>
You can start the server in 2 ways, the first is using **docker** and **docker-compose** and the second is **manually**
#### 📦 Docker way <a name="docker_way"></a>
You will need docker and docker-compose in your system.

To builds Docker image from  Dockerfile, run:
```bash
docker build --no-cache --force-rm --tag drones_restapi .
```
Use docker-compose to start the container:
```bash
docker-compose up
```

#### 🔧 Manual way  <a name="manual_way"></a>

Run:
```bash
go mod download
go mod vendor
go build
```

#### 🌍 Environment variables
The environment variable is exported with the location of the server configuration file.

If you have 🐧Linux or 🍎Dash, run:
```bash
export SERVER_CONFIG=$PWD/conf/conf.yaml
```
but if it is in the windows cmd, then run:
```bash
set SERVER_CONFIG=%cd%/conf/conf.yaml
```
#### 🏃🏽‍♂️ Start the server
Before it is recommended that you read more about the server configuration file in the section 👉🏾  .

Run the server:
```bash
./drones.restapi
```

and visit the swagger docs:

> http://localhost:7001/swagger/index.html

The first endpoint to execute must be /api/v1/database/populate [POST], to populate the database. That endpoint does not need authentication.

![swagger ui](/docs/images/populate_endpoint.png)

You can then authenticate and test the remaining endpoints.

### 🧪 Unit or End-To-End Testing
Run:
```bash
go test -v
```

## 🔨 Tech and packages <a name="tech"></a>
* [Iris Web Framework](https://github.com/kataras/iris)
* [Buntdb](https://github.com/tidwall/buntdb)
* [govalidator](https://github.com/asaskevich/govalidator)
* [gocron](https://github.com/go-co-op/gocron)
* [swag](https://github.com/swaggo/swag)
* [Docker](https://docs.docker.com)
* [docker-compose](https://docs.docker.com/compose/)

## 📐 Architecture <a name="arch"></a>
This project has 3 layer :

- Controller Layer (Presentation)
- Service Layer (Business)
- Repository Layer (Persistence)


Tag | Path | Layer |
--- | ---- | ----- |
Auth     | api/endpoints/end_auth.go | Controller | 
Drones   | api/endpoints/end_drones.go |  Controller |
EventLog | api/endpoints/end_eventlog.go |  Controller |
 |  |  |
Auth     | service/auth/svc_authentication.go | Service | 
Drones   | service/svc_drones.go |  Service |
EventLog | service/cron/svc_eventlog.go |  Service |
 |  |  |
Auth     | repo/db/repo_drones.go | Repository | 
Drones   | repo/db/repo_drones.go |  Repository |
EventLog | repo/db/repo_eventlog.go |  Repository |