# pgtocsv

`pgtocsv` executes a query on a PostgreSQL database and outputs the results in CSV.

Why not just use psql and `\copy`?

* `pgtocsv` has easier syntax than `\copy`
* `\copy` requires the entire query be given on one line
* `pgtocsv` can read the query from a file

## Installation

The Go tool chain must be installed.

```
$ go get -u github.com/jackc/pgtocsv
```

## Example usage

```
$ pgtocsv -s 'select * from users'
```

```
$ pgtocsv -f query.sql
```
