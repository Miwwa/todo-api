# ToDo API Application

## Overview

Made for https://roadmap.sh/projects/todo-list-api

## Installation

Requirements:
- [Go](https://golang.org/dl/) installed
- `migrate` tool installed

`go install -tags 'sqlite3' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`


To get started with this project, clone the repository:
```sh
git clone https://github.com/Miwwa/todo-api
cd todo-api
```

## Configuration

Configuration for the application can be set via environment variables or a configuration file.
Create a `.env` file in the root directory with the following content:

```env
IS_PRODUCTION=true
PORT=3000
SQLITE_DB_PATH=./db.sqlite
JWT_SECRET=mySecret
```

## Usage

Run following command to create a local sqlite database
```sh
migrate -source file://db_migrations -database sqlite3://db.sqlite up
```

To run the application locally, use the following command:

```sh
go run main.go
```

The server will start on `http://localhost:3000` by default.

## License

Distributed under the MIT License. See `LICENSE` for more information.
