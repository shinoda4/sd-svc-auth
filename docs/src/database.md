# Database

By default, when running `docker-composed.yml`, the docker will create a postgresql container with database `sd_auth`.

The username is `sd_auth`, and password is `sd_pass`.

running the following commands to migrate database.

> If had already running composed, recreate database using following commands.
> ```shell
> docker compose down
> docker volume rm deployments_pgdata
> docker compose up -d
> ```

## Migration

```shell
brew install golang-migrate
```

```shell
migrate create -ext sql -dir db/migrations [name]
```

```shell
migrate -source file://db/migrations -database "postgres://sd_auth:sd_pass@localhost:5432/sd_auth?sslmode=disable" up
```
