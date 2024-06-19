# How to run migrations?

1. Make sure you are in the root of the repo:

```bash
cd <path_to_the_repo_root>
```

2. Create a migration files:

```bash
migrate create -ext sql -dir migrations -seq <migration_name>
```

3. Update auto-generated files. See `/migrations` dir for examples.

4. Export the `POSTGRESQL_URL` environment variable. For localhost, it is something along these lines.

```bash
export POSTGRESQL_URL='postgres://root:root@localhost:5432/ohno_db?sslmode=disable'
```

5. Apply the migrations.

```bash
migrate -database ${POSTGRESQL_URL} -path migrations up
```

6. Rollback if needed.

```bash
migrate -database ${POSTGRESQL_URL} -path migrations down
```
