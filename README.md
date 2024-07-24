## Quickstart (docker compose)

Ensure you have the following requirements

* [Docker Compose](https://docs.docker.com/compose/install/)
* [Golang](https://go.dev/doc/install)
* [Postman](https://www.postman.com/downloads/)

Deploy backend API Platform with following this command:

```
docker-compose up -d
```

Wait for the container to be ready:

```
docker-compose ps
```

Import this file postman collection:

[API Platform - REST APIs](./API Platform - REST APIs.postman_collection.json)

## Monitoring

### Jaeger

Open this link with browser:

```
http://localhost:16686
```

### Prometheus

Open this link with browser:

```
http://localhost:9090
```

### Grafana

Open this link with browser:

```
http://localhost:9000
```

username: `admin`

password: `admin`
