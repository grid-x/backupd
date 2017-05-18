package main

import (
	"fmt"

	"github.com/grid-x/backupd/pkg/backup"
	"github.com/grid-x/backupd/pkg/datastore"
	"github.com/grid-x/backupd/pkg/storage"
)

func main() {

	s3 := storage.NewS3("eu-central-1", "gridx-de-staging-k8s-backups")

	etcd := datastore.NewEtcd("http://localhost:4001")
	influx := datastore.NewInflux("localhost:8088", "test")
	mongo := datastore.NewMongoDB("localhost", 27017, "", "")

	statusc := make(chan backup.BackupJobStatus)
	schedules := []backup.Schedule{
		backup.Schedule{
			Spec:      "@every 20s",
			BackupJob: backup.NewBackupJob(etcd, s3, "etcd-events-backup", statusc),
		},
		backup.Schedule{
			Spec:      "@every 30s",
			BackupJob: backup.NewBackupJob(influx, s3, "influx-backup", statusc),
		},
		backup.Schedule{
			Spec:      "@every 34s",
			BackupJob: backup.NewBackupJob(mongo, s3, "mongo-backup", statusc),
		},
	}

	s, err := backup.NewScheduler(schedules)
	if err != nil {
		fmt.Println(err)
		return
	}
	go s.Run()

	for s := range statusc {
		fmt.Println(s)
	}
}
