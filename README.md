Check out the other branches for a MongoDB template.

# GO + Docker REST API Template

A template of a REST API written in Go with a development environment run using Docker.

## Requirements:
- [Docker](https://www.docker.com/products/docker-desktop/)
    - Docker Compose
- [Go](https://go.dev/)
- [Make](https://www.gnu.org/software/make/)

## Local Dev Environment
### Running the server
```
make start-server
```
- This will use Docker Compose to spool up all the required Docker containers
- The Go server will be run using [air](https://github.com/cosmtrek/air) which automatically restarts the server when code changes are made
- Server environment variables are located in the `docker-compose.yaml` file.

## Database Migrations
Database migrations are managed by [Golang-Migrate](https://github.com/golang-migrate/migrate) using docker containers. You may need to configure the Make scripts depending on your Docker network and Postgres configurations defined in the `docker-compose.yaml` file.

### Creating Migrations
```
make mm migration_name=
```
- Migrations will be created in `migrations/sql`
### Apply all migrations
```
make m
```
### Revert a single migration
```
make migrate-down
```
