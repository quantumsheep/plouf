package plouf

import (
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type IController interface {
	IInjectable

	InitRoutes(e *echo.Echo)
}

type Controller struct {
	Injectable
}

func (c Controller) Init(self IInjectable) error {
	if err := c.Injectable.Init(self); err != nil {
		return err
	}

	if c.ShouldLogInjection(self) {
		logrus.Debugf("Initializing controller %s", ReflectTypeName(self))
	}

	return nil
}

func (c Controller) ShouldLogInjection(self IInjectable) bool {
	return ReflectValue(self).Type() != reflect.TypeOf(Controller{})
}
