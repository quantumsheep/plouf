package plouf

import (
	"reflect"

	"github.com/sirupsen/logrus"
)

type IService interface {
	IInjectable
}

type Service struct {
	Injectable
}

func (s Service) Init(self IInjectable) error {
	if err := s.Injectable.Init(self); err != nil {
		return err
	}

	if s.ShouldLogInjection(self) {
		logrus.Debugf("Initializing service %s", ReflectTypeName(self))
	}

	return nil
}

func (s Service) ShouldLogInjection(self IInjectable) bool {
	return ReflectValue(self).Type() != reflect.TypeOf(Service{})
}
