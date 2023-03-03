# pgtocsv

`pgtocsv` executes a query on a PostgreSQL database and outputs the results in CSV.

## Project Retired

As of PostgreSQL 12, `psql` can do everything `pgtocsv` can. For example:

```
psql --quiet --no-psqlrc --csv  -c 'select * from users' > users.csv
```

Or for TSV output to a file instead of stdout:

```
psql --quiet --no-psqlrc --csv -c "\pset csv_fieldsep '\t'" -c 'select * from users' -o users.tsv
```

## Installation

The Go tool chain must be installed.

```
$ go get -u github.com/jackc/pgtocsv
```

## Configuring Database Connection

`pgtocsv` supports the standard `PG*` environment variables. In addition, the `-d` flag can be used to specify a database URL.

## Example usage

```
$ pgtocsv -s 'select * from users'
id,name
1,jack
```

```
$ pgtocsv -f query.sql
id,first_name,last_name,sex,birth_date,weight,height,update_time
1,Hunter,Halvorson,male,2004-06-13,66,61,2006-05-24 11:32:33-05
2,Sigrid,Kub,male,2002-06-09,10,70,2010-11-18 08:50:07-06
3,Alta,Luettgen,male,2006-07-31,336,73,1986-07-16 12:10:13-05
4,Nestor,Schulist,female,1999-03-12,171,24,2004-04-08 21:09:24-05
5,Carolyn,Yundt,female,2003-06-18,275,72,2013-05-25 08:31:50-05
```

## Related

See also the sibling project [csvtopg)](https://github.com/jackc/csvtopg) which simplifies importing CSV data.
