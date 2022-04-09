package plouf

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

type IInjectable interface {
	Init(self IInjectable) error

	ShouldLogInjection(self IInjectable) bool
}

type Injectable struct{}

func (i *Injectable) Init(self IInjectable) error {
	value := reflectValue(self)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		if !field.CanInterface() {
			continue
		}

		if injectable, ok := field.Interface().(IInjectable); ok {
			if wasInitialized := markInjectableAsInitialized(field); wasInitialized {
				continue
			}

			if err := injectable.Init(injectable); err != nil {
				return err
			}
		}
	}

	return nil
}

func (i *Injectable) ShouldLogInjection(self IInjectable) bool {
	return false
}

type injectableStatus struct {
	Injectable  IInjectable
	Initialized bool
}

var injectables = make(map[reflect.Type]*injectableStatus)

func markInjectableAsInitialized(field reflect.Value) (wasInitialized bool) {
	if field.Kind() == reflect.Ptr {
		if status, ok := injectables[field.Type().Elem()]; ok {
			if status.Initialized {
				return true
			}

			status.Initialized = true
		}
	}

	return false
}

func Inject(injectable IInjectable) error {
	if injectable.ShouldLogInjection(injectable) {
		logrus.Debugf("Injecting %s", reflect.TypeOf(injectable).Elem().Name())
	}

	value := reflect.ValueOf(injectable).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		if !field.CanInterface() {
			continue
		}

		if _, ok := field.Interface().(IInjectable); !ok {
			continue
		}

		typ := field.Type()

		var value reflect.Value

		// Use cache if possible, else assign new value
		if typ.Kind() == reflect.Ptr {
			if value, ok := injectables[typ.Elem()]; ok {
				field.Set(reflect.ValueOf(value))
				continue
			}

			value = reflect.New(typ.Elem())
		} else {
			value = reflect.New(typ)
		}

		// Save new value to cache and set field value
		if typ.Kind() == reflect.Ptr {
			injectables[typ.Elem()] = &injectableStatus{
				Injectable:  value.Interface().(IInjectable),
				Initialized: false,
			}
			field.Set(value)
		} else {
			field.Set(value.Elem())
		}

		// Inject new value dependencies
		child := value.Interface().(IInjectable)
		Inject(child)
	}

	return nil
}
