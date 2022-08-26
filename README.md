# 🛰 drones
REST API that allows clients to communicate with drones (i.e. **dispatch controller**).

## Table of Contents

- [API specification](#api_spec)
- [Assigned tasks](#assigned_task)
- [Configuration file](#config_file)
- [Get Started](#get_started)
  * [Deployment ways (2 ways)](#deploy_ways)
    - [Docker way](#docker_way)
    - [Manual way](#manual_way)
- [Tech](#tech)
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

> http://127.0.0.1:7001/swagger/index.html

![swagger ui](/docs/images/swagger-ui.png)

## 📝 Assigned tasks <a name="assigned_task"></a>
|  Done          | Task       | Endpoint                              |
| -------------- | -----------|------------------------- |
| ✅ | registering a drone;                                | 👉🏾 endpoint: `/api/v1/drones  [POST]`
| ✅ | loading a drone with medication items;              | 👉🏾 endpoint: `/api/v1/medicationsitems/:serialNumber [POST]`
| ✅ | checking loaded medication items for a given drone; | 👉🏾 endpoint: `/api/v1/medicationsitems/:serialNumber [GET]`
| ✅ | checking available drones for loading;              | 👉🏾 endpoint: `drones?state=1 [GET]`
| ✅ | check drone battery level for a given drone;        | 👉🏾 endpoint: `Get a drone by serialNumber [GET]`

## 🛠️️ Configuration of conf.yaml <a name="config_file"></a>
👉🏾 ![The config file](/conf/conf.yaml)

|  Param      | Description       | default value   |
| ----------- | -----------|------------------------- |
| ApiDocIP    | IP to expose the api  | 127.0.0.1
| DappPort    | app PORT              | 7001
| StoreDBPath | DB file location      | ./db/data.db
| CronEnabled | active the cron job   | true
| LogDBPath   | DB file event logs    | ./db/event_log.db
| EveryTime   | time interval (in seconds) that the cron task is executed | 300 seconds (every 5 minutes)

By default **StoreDBPath** are using the database file located in the db folder at the root of the project.

The already populated drone DB can be removed if desired. The server exposes the `/api/v1/database/populate` POST endpoint to generate and repopulate the database if necessary.
## ⚡ Get Started <a name="get_started"></a>

Download the drones.restapi project:
```bash
git clone https://github.com/kmilodenisglez/drones.restapi.git
```
Move to the root of the project:
```bash
cd drones.restapi
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

### 🔧 Manual way  <a name="manual_way"></a>

Run:
```bash
go mod download
go mod vendor
go build
```

#### 🌍 Environment variables
The environment variable is exported with the location of the server configuration file.

Run:
```bash
# linux, wsl or darwin
export SERVER_CONFIG=$PWD/conf/conf.yaml
```
but if it is in the windows cmd, then run:
```bash
# windows cmd
set SERVER_CONFIG=%cd%/conf/conf.yaml
```
#### 🏃🏽‍♂️ Start the server
Before it is recommended that you read more about the server configuration file in the section 👉🏾  .

Run the server:
```bash
./drones.restapi
```

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

