#

- Chi, Sqlc, JWT auth

# Prepare Postgresql container

- `Postgresql` container will use `port 5432` in the host machine
- `pgadmin` will be accessed from `localhost:80`
  - username: `admin@example.com`
  - password: `password`
- When register pg server in pgAdmin, the `Host name/address` field is the `Postgresql's container name: postgres_container`
- refs
  - [docker-compose file](https://hevodata.com/learn/pgadmin-docker/#:~:text=pgAdmin%20is%20an%20excellent%20tool,our%20environment%20up%20within%20minutes.)
  - pgadmin
    - [deployment](https://www.pgadmin.org/docs/pgadmin4/latest/container_deployment.html)
    - [usage](https://www.youtube.com/watch?v=WFT5MaZN6g4)

```bash
cd ./
docker compose up -d # the Postgresql service
```

# TODO

- TBD
