# backupd
Microservice for performing backups

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
