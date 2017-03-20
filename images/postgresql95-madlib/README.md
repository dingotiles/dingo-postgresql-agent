# Madlib with Dingo PostgreSQL

Thanks to @mwright-pivotal for coding up a madlib version of Dingo PostgreSQL:

To build this image from the root folder:

```
docker build -t dingotiles/dingo-postgresql-madlib images/postgresql95-madlib
```

Then run the following and follow the prompts to run Dingo PostgreSQL/madlib:

```
docker run dingotiles/dingo-postgresql-madlib
```

To access your PostgreSQL, and install the madlib extension:

```
psql $(docker exec -ti dingo-postgresql cat /config/uri) -c "create extension madlib;"
```

To play with your PostgreSQL with madlib:

```
psql $(docker exec -ti dingo-postgresql cat /config/uri)
```
