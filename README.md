# backupd

backupd is a small service for performing backups. It implements several
datastores that describe how to perform backups for the respective database.

So far we have support for:

* InfluxDB
* MongoDB
* Postgres
* etcd

The procedure is generally the same:

The database needs to be accessible over network and the backup needs to happen
over network. The implementation will start a subprocess, e.g. `pg_dump` or
`influxd backup ...` to create a local backup file (if there are multiple files
they are compressed into an archive) and then upload those files.

## Develop

```
# Build
make

# Lint
make lint

# Tests
make test

# Build the docker image
make docker

# Push the image (if you have access to the repository)
make push
```


## How to restore backups produced by backupd

### InfluxDB

```shell
# Extract tar archive
tar -xvzf influxdb-archive-12345
# This will produce several files named <database>.autogen.nnnnn.nn and meta.nn where
# n is a number between 0 and 9
#
# To restore: Run the following on the database instance
influxd restore -database <database> -datadir /usr/local/var/influxdb/data -metadir /usr/local/var/influxdb/meta .
```

### MongoDB

```shell
# Use the command below and adapt the settings such as host, port and archive
# file to your env
mongorestore -h 192.168.99.100 --port 27017 --archive=mongo-416184046 --gzip
```

### etcd

The etcd backup exporter uses plain JSON files to store the data.

To restore such a backup the easiest way is to use [etcdtool](https://github.com/mickep76/etcdtool)

```shell
etcdtool -p http://192.168.99.100:2379 import -y -r /kubernetes etcd-backup.json
```
