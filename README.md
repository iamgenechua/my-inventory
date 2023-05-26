# My Inventory

This project is a simple inventory management application built with Go, using Gorm for database ORM, REST API for communication, and includes testing.

## Prerequisites

Before running the project, make sure you have the following dependencies installed:

- Go: https://golang.org/doc/install
- Docker: https://docs.docker.com/get-docker/

## Getting Started

Follow these steps to run the project:

1. Start a PostgreSQL container:

```
docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres
```

2. Build go project
```
go build
```

3. Run the executable
```
./my-inventory
```

The application should now be running and accessible at http://localhost:10000/projects.
