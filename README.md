# drones
REST API that allows clients to communicate with drones (i.e. **dispatch controller**).

## Table of Contents

- [API specification](#api_spec)
- [Assigned tasks](#assigned_task)
- [Tech](#tech)
## API specification <a name="api_spec"></a>

The **Drone API server** provides the following API with communicating the **DB**:

| Tag           | Title                              | URL                                      | Query | Method |
| ------------- | ---------------------------------- | ---------------------------------------- | ----- | ---- |
| Auth          | user authentication                | `/api/v1/auth`                           |   -   |`POST`|
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

## Assigned tasks <a name="assigned_task"></a>
- ✅ registering a drone;                   👉🏾 endpoint: `/api/v1/drones  [POST]`
- ✅ loading a drone with medication items; 👉🏾 endpoint: `/api/v1/medicationsitems/:serialNumber [POST]`
- ✅ checking loaded medication items for a given drone; 👉🏾 endpoint: `/api/v1/medicationsitems/:serialNumber [GET]`
- ✅ checking available drones for loading; 👉🏾 endpoint: `drones?state=1 [GET]`
- ✅ check drone battery level for a given drone; 👉🏾 endpoint: `Get a drone by serialNumber [GET]`

## Tech <a name="tech"></a>
* [Iris Web Framework](https://github.com/kataras/iris)
* [Buntdb](https://github.com/tidwall/buntdb)
* [govalidator](https://github.com/asaskevich/govalidator)
* [gocron](https://github.com/go-co-op/gocron)
* [swag](https://github.com/swaggo/swag)