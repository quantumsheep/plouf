package plouf_modules

import (
	"fmt"

	"github.com/quantumsheep/plouf"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseConnection struct {
	plouf.Service

	ConfigService *ConfigService

	Database *gorm.DB
}

func (d *DatabaseConnection) Init(self plouf.IInjectable) error {
	var dialector gorm.Dialector

	driver := d.ConfigService.Get("database.driver")
	switch driver {
	case "":
		return fmt.Errorf(`"database.driver" is not configured`)
	case "sqlite":
		path := d.ConfigService.Get("database.path")
		if path == "" {
			return fmt.Errorf(`"database.path" is not configured`)
		}

		dialector = sqlite.Open(path)
	default:
		return fmt.Errorf(`configured "database.driver" ("%s") is not supported`, driver)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}

	d.Database = db

	return nil
}
