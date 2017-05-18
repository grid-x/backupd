package main

import (
	"flag"

	log "github.com/Sirupsen/logrus"

	"github.com/grid-x/backupd/pkg/backup"
	"github.com/grid-x/backupd/pkg/config"
	"github.com/grid-x/backupd/pkg/datastore"
	"github.com/grid-x/backupd/pkg/storage"
)

func main() {
	var (
		logger = log.New()

		configFile = flag.String("config-file", "config.yaml", "Config file to use")
	)

	config, err := config.ReadConfigFromFile(*configFile)
	if err != nil {
		logger.Fatalf("Problem reading config file: %s", err.Error())
	}

	if config.Storage.S3 == nil {
		log.Fatal("Not storage setting given")
	}

	s3 := storage.NewS3(config.Storage.S3.Region, config.Storage.S3.Bucket)

	if len(config.Targets) < 1 {
		log.Warn("No target provided")
	}

	statusc := make(chan backup.BackupJobStatus)
	var schedules []backup.Schedule

	var ds backup.DataStore
	for _, t := range config.Targets {
		switch t.Type {
		case "etcd":
			endpoint, ok := t.Settings["endpoint"].(string)
			if !ok {
				log.Errorf("Can't convert etcd endpoint to string for target %s", t.Name)
				continue
			}
			ds = datastore.NewEtcd(endpoint)
		case "influxdb":
			endpoint, ok := t.Settings["endpoint"].(string)
			if !ok {
				log.Errorf("Can't convert influxdb endpoint to string for target %s", t.Name)
				continue
			}
			database, ok := t.Settings["database"].(string)
			if !ok {
				log.Errorf("Can't convert influxdb database to string for target %s", t.Name)
				continue
			}
			ds = datastore.NewInflux(endpoint, database)
		case "mongodb":
			host, ok := t.Settings["host"].(string)
			if !ok {
				log.Errorf("Can't convert mongodb host to string for target %s", t.Name)
				continue
			}

			port, ok := t.Settings["port"].(int)
			if !ok {
				log.Errorf("Can't convert mongodb port to string for target %s", t.Name)
				continue
			}

			useri, exists := t.Settings["user"]
			user := ""
			if exists {
				user, ok = useri.(string)
				if !ok {
					log.Errorf("Can't convert mongodb user to string for target %s", t.Name)
					continue
				}
			}

			passwordi, exists := t.Settings["password"]
			password := ""
			if exists {
				password, ok = passwordi.(string)
				if !ok {
					log.Errorf("Can't convert mongodb password to string for target %s", t.Name)
					continue
				}
			}
			ds = datastore.NewMongoDB(host, port, user, password)
		default:
			log.Warnf("Storage type '%s' currently not supported", t.Type)
			continue
		}
		schedules = append(schedules, backup.Schedule{
			Spec:      t.Schedule,
			BackupJob: backup.NewBackupJob(ds, s3, t.Name, statusc),
		})
		log.Infof("Successfully added job for target: %s", t.Name)
	}

	s, err := backup.NewScheduler(schedules)
	if err != nil {
		log.Fatal(err)
	}

	go s.Run()
	for s := range statusc {
		if s.Error == nil {
			log.Infof("Success: Job %s took %s to finish", s.Name, s.Duration)
		} else {
			log.Errorf("Job %s failed after %s with error %s", s.Name, s.Duration, s.Error.Error())
		}
	}
}
