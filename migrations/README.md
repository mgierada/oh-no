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

# How to manually change/remove/list triggers and trigger functions

To list all trigger function, use this query:

```sql
SELECT n.nspname as schema,
       p.proname as function_name
FROM pg_proc p
JOIN pg_namespace n ON p.pronamespace = n.oid
WHERE p.prorettype = 'pg_catalog.trigger'::pg_catalog.regtype;
```

To list all triggers, use this query:

```sql
SELECT event_object_schema as table_schema,
       event_object_table as table_name,
       trigger_name,
       action_timing as trigger_time,
       event_manipulation as event,
       action_statement as definition
FROM information_schema.triggers
ORDER BY table_schema, table_name, trigger_name;
```

To drop trigger (if any), use this mutation:

```sql
DROP TRIGGER IF EXISTS your_trigger_name ON your_table_name;
```

To drop trigger function, use this mutation:

```sql
DROP FUNCTION IF EXISTS your_schema_name.your_trigger_function_name();
```
