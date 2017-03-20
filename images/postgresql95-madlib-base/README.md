# Madlib with Dingo PostgreSQL

Thanks to @mwright-pivotal for coding up a madlib version of Dingo PostgreSQL:

To build this image from the root folder:

```
docker build -t dingotiles/dingo-postgresql-madlib-base images/postgresql95-madlib-base
```

Then run the following and follow the prompts to run Dingo PostgreSQL/madlib:

```
docker run dingotiles/dingo-postgresql-madlib
```

To access your PostgreSQL, and install the madlib extension:

```
uri=$(docker exec -ti dingo-postgresql cat /config/uri)
psql $uri -c "create extension plpythonu;"
psql $uri -c "create extension madlib;"
```

To play with your PostgreSQL with madlib:

```
psql $uri
```
