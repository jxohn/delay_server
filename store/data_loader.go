package store

import (
	"os"
	"time"
)

type DataLoader struct {
	logDir os.File
}

func (d *DataLoader) LoadNextHour(time time.Time) error {
	d.logDir.Readdir()
}