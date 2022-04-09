package plouf

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Worker struct {
	e *echo.Echo

	mainModule IModule
}

func NewWorker(mainModule IModule) (*Worker, error) {
	if err := Inject(mainModule); err != nil {
		return nil, err
	}

	if err := mainModule.Init(mainModule); err != nil {
		return nil, err
	}

	e := echo.New()
	e.HideBanner = true
	e.Validator = NewDTOValidator()

	e.Use(middleware.Recover())

	if err := mainModule.InitControllersRoutes(mainModule, e); err != nil {
		return nil, err
	}

	return &Worker{
		e:          e,
		mainModule: mainModule,
	}, nil
}

func (w *Worker) Start(address string) error {
	return w.e.Start(address)
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
