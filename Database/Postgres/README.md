# Postgres

```sql
docker run --name goddess -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres
docker exec -it goddess psql -U root

docker compose up -d
docker compose exec goddess psql -U postgres

docker compose exec goddess sh
psql --help
```

## Psql

```sql
\q : quit
\l : list database
\? : help
```

### Database
```sql
-- Create
CREATE DATABASE test;

-- Destroy
DROP DATABASE test;
DROP DATABASE IF EXISTS test;

-- Connect to database test
1. docker compose exec goddess psql -U postgres test
2. psql -h localhost -p 5432 -U postgres test
3. \c test
```

### Table

* [Data Types](https://www.postgresql.org/docs/current/datatype.html)

```sql
-- Create
CREATE TABLE <name> (
    <column> <type> <constraints>
    ...
);

CREATE TABLE users (
    id       serial         PRIMARY KEY,
    username varchar (50)   UNIQUE NOT NULL,
    email    varchar (255)  UNIQUE NOT NULL
);

-- List
\d

-- Structure
\d users

-- Destroy
DROP TABLE users;
DROP TABLE IF EXISTS users;

-- Insert
INSERT INTO users 
(username,      email) 
VALUES
('Golang',      'Golang.org'),
('Flutter',     'Flutter.dev');

-- Show
SELECT * FROM users;

-- Update
UPDATE users
SET email = 'Go.dev'
WHERE username = 'Golang';

-- Delete
DELETE FROM users
WHERE username = 'Golang';

-- Query
SELECT DISTINCT * FROM users
WHERE username IN ('Flutter', 'Golang') AND email LIKE '%.dev%' AND (username = 'Flutter' OR email = 'Flutter.dev') AND email LIKE '%F______.dev%' AND email ILIKE 'f%'
ORDER BY id,email desc
OFFSET 0
FETCH FIRST 5 ROW ONLY;

-- Count
SELECT COUNT(*) FROM users;

-- Avg
SELECT AVG(id) FROM users;

-- Join
CREATE TABLE orders (
    order_id serial PRIMARY KEY,
    user_id INT REFERENCES users(id),
    amount DECIMAL NOT NULL
);

INSERT INTO orders (user_id, amount) VALUES
(1, 100.50),
(2, 75.25);

SELECT users.username, orders.amount FROM users
JOIN orders ON users.id = orders.user_id;

-- Alter
ALTER TABLE users
ADD COLUMN complexity DECIMAL(10, 2);

-- Group
SELECT user_id, COUNT(order_id) FROM orders GROUP BY user_id;


-- Operator
= , < , > , >= , <= , !=
LIKE , ILIKE , SIMILAR TO
AND , OR , NOT
+ , - , * , / , % , ^
IS NULL , IS NOT NULL , IN
```
