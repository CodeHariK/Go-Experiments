# SandSlash

```sql


atlas schema inspect \
    -u $postgres \
    --format '{{ sql . }}' > store/schema/schema.sql

atlas schema apply \
    -u $postgres \
    --to file://store/schema/schema.sql \
    --dev-url "docker+postgres://_/postgres:alpine/dev"

atlas migrate push star \
    --dev-url "docker+postgres://_/postgres:alpine/dev" \
    --dir file://store/schema

atlas migrate hash --dir file://store/schema

atlas schema inspect \
    -u $postgres \
    --web
```
