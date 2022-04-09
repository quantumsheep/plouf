package plouf

import (
	"reflect"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type IModule interface {
	IInjectable

	InitControllersRoutes(self IInjectable, e *echo.Echo) error
}

type Module struct {
	Injectable
}

func (m Module) Init(self IInjectable) error {
	if err := m.Injectable.Init(self); err != nil {
		return err
	}

	if m.ShouldLogInjection(self) {
		logrus.Debugf("Initializing module %s", ReflectTypeName(self))
	}

	return nil
}

func (m Module) ShouldLogInjection(self IInjectable) bool {
	return ReflectValue(self).Type() != reflect.TypeOf(Module{})
}

func (m Module) InitControllersRoutes(self IInjectable, e *echo.Echo) error {
	value := ReflectValue(self)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		if !field.CanInterface() {
			continue
		}

		if _, ok := field.Interface().(IInjectable); ok {
			if controller, ok := field.Interface().(IController); ok {
				logrus.Debugf("Initializing routes for %s", ReflectTypeName(controller))
				controller.InitRoutes(e)
			}

			if module, ok := field.Interface().(IModule); ok {
				if err := module.InitControllersRoutes(module, e); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
