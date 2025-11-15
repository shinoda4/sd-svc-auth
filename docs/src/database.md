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

```shell
psql -U sd_auth -h 127.0.0.1 -f sql/init.sql
```