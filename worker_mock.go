package plouf

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type WorkerMock struct {
	*Worker
}

func NewWorkerMock(mainModule IModule) (*WorkerMock, error) {
	e := echo.New()
	e.Validator = NewDTOValidator()

	worker := &WorkerMock{
		Worker: &Worker{
			e:           e,
			injectables: make(map[reflect.Type]*injectableStatus),
			mainModule:  mainModule,
		},
	}

	if err := worker.Inject(mainModule); err != nil {
		return nil, err
	}

	if err := mainModule.Init(mainModule); err != nil {
		return nil, err
	}

	return worker, nil
}

func (w *WorkerMock) NewContext(req *http.Request, res http.ResponseWriter) echo.Context {
	return w.e.NewContext(req, res)
}
