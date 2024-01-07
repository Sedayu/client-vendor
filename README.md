# Client Vendor API

## How To Run

### Requirements

To be able to run the database, you need to install these tools:

1. Golang
2. Golang Migrate CLI


### Run The Application

Copy .env file

```
cp env.sample .env
```
Create database in your postgresql, with dbeaver same as `POSTGRES_DB` env in your .env. After that run the migration

```
migrate -source file://migrations -database postgres://username:password@localhost:5432/client_vendor_local?sslmode=disable up
```

Run the applications

```
go run main.go
```

## Folder Structure

```
.
├── entity
├── handler
├── migrations
├── repository
└── service
```

### Entity

Business entity should be written here.

### Handler

HTTP handler should be written here.

### Migrations

All migrations file should be written here.

### Repository

All file related to get data from database should be written here.

### Service

All business logic.