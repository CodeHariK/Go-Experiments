# Postgres

```sql
docker compose up -d
docker compose exec godess psql -U postgres

docker compose exec godess sh
psql --help

```

## Psql

```sql
\q : quit
\l : list database
\? : help

CREATE DATABASE test;

-- Connect to database test
1. docker compose exec godess psql -U postgres test
2. psql -h localhost -p 5432 -U postgres test
3. \c test

```
