package plouf

import (
	"github.com/labstack/echo/v4"
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
