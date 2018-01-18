package main

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/grid-x/backupd/pkg/backup"
	"github.com/grid-x/backupd/pkg/config"
	"github.com/grid-x/backupd/pkg/datastore"
	"github.com/grid-x/backupd/pkg/storage"
)

func main() {
	var (
		logger = log.New().WithFields(log.Fields{
			"component": "main",
		})

		configFile = flag.String("config-file", "config.yaml", "Config file to use")
	)
	flag.Parse()

	config, err := config.ReadConfigFromFile(*configFile)
	if err != nil {
		logger.Fatalf("Problem reading config file: %s", err.Error())
	}

	if config.Storage.S3 == nil {
		logger.Fatal("Not storage setting given")
	}

	s3 := storage.NewS3(config.Storage.S3.Region, config.Storage.S3.Bucket)

	if len(config.Targets) < 1 {
		logger.Warn("No target provided")
	}

	statusc := make(chan backup.JobStatus)
	var schedules []backup.Schedule

	var ds backup.DataStore
	for _, t := range config.Targets {
		logger.Infof("Processing target %s", t.Name)
		switch t.Type {
		case "etcd":
			endpoint, ok := t.Settings["endpoint"].(string)
			if !ok {
				logger.Errorf("Can't convert etcd endpoint to string for target %s", t.Name)
				continue
			}
			ds = datastore.NewEtcd(endpoint)
		case "influxdb":
			endpoint, ok := t.Settings["endpoint"].(string)
			if !ok {
				logger.Errorf("Can't convert influxdb endpoint to string for target %s", t.Name)
				continue
			}
			database, ok := t.Settings["database"].(string)
			if !ok {
				logger.Errorf("Can't convert influxdb database to string for target %s", t.Name)
				continue
			}

			days, ok := t.Settings["daysToKeep"].(int)
			var last *time.Duration
			if ok {
				last = new(time.Duration)
				*last = 24 * time.Hour * time.Duration(days)
			}
			ds = datastore.NewInflux(endpoint, database, last)
		case "mongodb":
			host, ok := t.Settings["host"].(string)
			if !ok {
				logger.Errorf("Can't convert mongodb host to string for target %s", t.Name)
				continue
			}

			port, ok := t.Settings["port"].(int)
			if !ok {
				logger.Errorf("Can't convert mongodb port to string for target %s", t.Name)
				continue
			}

			user, ok := t.Settings["user"].(string)
			if !ok {
				logger.Errorf("Can't convert mongodb user to string for target %s", t.Name)
				continue
			}
			password, ok := t.Settings["password"].(string)
			if !ok {
				logger.Errorf("Can't convert mongodb password to string for target %s", t.Name)
				continue
			}
			ds = datastore.NewMongoDB(host, port, user, password)
		case "postgres":
			url, ok := t.Settings["url"].(string)
			if !ok {
				logger.Errorf("Can't convert postgres url to string for target %s", t.Name)
				continue
			}
			ds = datastore.NewPostgres(url)
		default:
			logger.Warnf("Target type '%s' currently not supported", t.Type)
			continue
		}
		conf := backup.Config{
			Name:          t.Name,
			TempDir:       config.Settings.TmpDir,
			TempDirPrefix: "",
		}
		schedules = append(schedules, backup.Schedule{
			Spec: t.Schedule,
			Job:  backup.NewJob(ds, s3, conf, statusc),
		})
		logger.Infof("Successfully added job for target: %s", t.Name)
	}

	s, err := backup.NewScheduler(schedules)
	if err != nil {
		logger.Fatal(err)
	}

	go s.Run()
	for s := range statusc {
		if s.Error == nil {
			logger.Infof("Success: Job %s took %s to finish", s.Name, s.Duration)
		} else {
			logger.Errorf("Job %s failed after %s with error %s", s.Name, s.Duration, s.Error.Error())
		}
	}
}
