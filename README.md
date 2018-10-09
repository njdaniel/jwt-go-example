# jwt-go-example

Database Postgres:
CREATE TABLE products
(
    id SERIAL,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
    CONSTRAINT products_pkey PRIMARY KEY (id)
)

Docker:
$ env GOOOS=linux GOARCH=amd64 go build --tags netgo --ldfags 'extldflags "-lm lstdc++ -static"'
 