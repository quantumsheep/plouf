package plouf_modules

import (
	"fmt"

	"github.com/quantumsheep/plouf"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DatabaseDriverInitializer func(driver string, config *ConfigService) (gorm.Dialector, error)

var DatabaseDrivers = map[string]DatabaseDriverInitializer{
	"sqlite": DatabaseDriverSQLite,
}

func DatabaseDriverSQLite(driver string, config *ConfigService) (gorm.Dialector, error) {
	path := config.Get("database.path")
	if path == "" {
		return nil, fmt.Errorf(`"database.path" is not configured`)
	}

	return sqlite.Open(path), nil
}

type DatabaseConnection struct {
	plouf.Service

	ConfigService *ConfigService

	Database *gorm.DB
}

func (d *DatabaseConnection) Init(self plouf.IInjectable) error {
	var dialector gorm.Dialector

	driver := d.ConfigService.Get("database.driver")
	if driver == "" {
		return fmt.Errorf(`"database.driver" is not configured`)
	}

	initializer, ok := DatabaseDrivers[driver]
	if !ok {
		return fmt.Errorf(`configured "database.driver" ("%s") is not supported`, driver)
	}

	dialector, err := initializer(driver, d.ConfigService)
	if err != nil {
		return err
	}

	d.Database, err = gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}
