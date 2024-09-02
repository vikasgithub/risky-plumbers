# Risky Plumbers Rest Service

## Setup
This rest services allows the caller to
- Create a new risk
- Get an existing risk information using id
- Get all the risks

### Prerequisites
- Golang version 1.22

### Run the App
- Download/clone the repository
- From the `risky-plumbers` directory, run the below command
```console
    go run cmd/main.go
```

- By default, the application runs on port `8080`. To change the port, create a yaml and provide the below configuration
```yaml
    server:
        port: <port>
```

- To make the custom configuration effective, application must be run using the below command

```console
    go run cmd/main.go --config=<path to the configuration file>
```

### Run the tests

```console
    go test ./...
```

## REST API

The REST API to access the risky plumbers service is described below.

### Get a list of Risks

#### Request

`GET /risks`

    curl -i -H 'Accept: application/json' http://localhost:8080/api/v1/risks

#### Response

    HTTP/1.1 200 OK

    [
        {
            "id": "a3e00a37-f82c-4eef-9f13-2d192cb0bfbe",
            "state": "open",
            "title": "title",
            "description": "desc1"
        },
        {
            "id": "b0d1aa81-e11c-4722-9f71-b89227a540bb",
            "state": "open",
            "title": "title",
            "description": "desc1"
        }
    ]

##### Notes
- Right now pagination is not supported

### Create a new Risk

#### Request

`POST /api/v1/risks`

    curl -XPOST -i -H 'Accept: application/json' -d '{"state":"open", "title":"t", "description": "d"}' http://localhost:8080/api/v1/risks


###### Notes 
- `state`, `title` and `description` in the request body are required
- `state` can only be one of `open`, `closed`, `accepted`, `investigating`
- `title` can have a maximum length of `128` characters
- `description` can have a maximum length of `4096` characters
- Unique Id for the risk object (`id`) is generated and returned as part of the response payload 

#### Response (Risk successfully Created)

    HTTP/1.1 201 Created

    {"id":"7215e2ec-7e2e-46df-92a9-bd34ea958e28","state":"open","title":"t","description":"d"}


#### Response (Invalid Parameters)

    HTTP/1.1 400 Bad Request

    {"status":"Invalid request.","error":"description: the length must be no more than 4096; state: must be a valid value; title: the length must be no more than 128."}


### Get a specific Risk

#### Request

`GET /risks/id`

    curl -i -H 'Accept: application/json' http://localhost:8080/api/v1/risks/1

#### Response (When a Risk is Found)

    HTTP/1.1 200 OK

    {
        "id": "a3e00a37-f82c-4eef-9f13-2d192cb0bfbe",
        "state": "open",
        "title": "title",
        "description": "desc1"
    }

#### Response (When a Risk is Not Found)

    HTTP/1.1 404 OK

    {"status":"Resource not found."}


