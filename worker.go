package plouf

import (
	"net/http"
	"reflect"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
)

type Worker struct {
	e *echo.Echo

	injectables map[reflect.Type]*injectableStatus
	mainModule  IModule
}

func NewWorker(mainModule IModule) (*Worker, error) {
	e := echo.New()
	e.HideBanner = true
	e.Validator = NewDTOValidator()

	e.Use(middleware.Recover())

	worker := &Worker{
		e:           e,
		injectables: make(map[reflect.Type]*injectableStatus),
		mainModule:  mainModule,
	}

	if err := worker.Inject(mainModule); err != nil {
		return nil, err
	}

	if err := mainModule.Init(mainModule); err != nil {
		return nil, err
	}

	if err := mainModule.InitControllersRoutes(mainModule, e); err != nil {
		return nil, err
	}

	return worker, nil
}

func (w *Worker) Start(address string) error {
	return w.e.Start(address)
}

type injectableStatus struct {
	Injectable  IInjectable
	Initialized bool
}

func (w *Worker) markInjectableAsInitialized(field reflect.Value) (wasInitialized bool) {
	if field.Kind() == reflect.Ptr {
		if status, ok := w.injectables[field.Type().Elem()]; ok {
			if status.Initialized {
				return true
			}

			status.Initialized = true
		}
	}

	return false
}

func (w *Worker) Inject(injectable IInjectable) error {
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
			if value, ok := w.injectables[typ.Elem()]; ok {
				field.Set(reflect.ValueOf(value))
				continue
			}

			value = reflect.New(typ.Elem())
		} else {
			value = reflect.New(typ)
		}

		// Save new value to cache and set field value
		if typ.Kind() == reflect.Ptr {
			w.injectables[typ.Elem()] = &injectableStatus{
				Injectable:  value.Interface().(IInjectable),
				Initialized: false,
			}
			field.Set(value)
		} else {
			field.Set(value.Elem())
		}

		// Inject new value dependencies
		child := value.Interface().(IInjectable)
		if err := w.Inject(child); err != nil {
			return err
		}
	}

	return nil
}

type DTOValidator struct {
	validator *validator.Validate
}

func (cv *DTOValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func NewDTOValidator() *DTOValidator {
	return &DTOValidator{validator: validator.New()}
}
