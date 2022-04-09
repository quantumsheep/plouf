package plouf_modules

import (
	"reflect"

	"github.com/quantumsheep/plouf"
	"github.com/sirupsen/logrus"
)

type Repository[T interface{}] struct {
	plouf.Injectable

	Connection *DatabaseConnection
}

func (r *Repository[T]) ShouldLogInjection(self plouf.IInjectable) bool {
	return plouf.ReflectValue(self).Type() != reflect.TypeOf(plouf.Service{})
}

func (r *Repository[T]) Init(self plouf.IInjectable) error {
	if err := r.Injectable.Init(self); err != nil {
		return err
	}

	if r.ShouldLogInjection(self) {
		logrus.Debugf("Migrating entity %s", reflect.TypeOf(new(T)).Elem().Name())
		return r.Connection.Database.AutoMigrate(new(T))
	}

	return nil
}

func (r *Repository[T]) Find(conds ...interface{}) ([]*T, error) {
	var values []*T
	tx := r.Connection.Database.Find(&values, conds...)
	return values, tx.Error
}

func (r *Repository[T]) FindOne(conds ...interface{}) (*T, error) {
	var value *T
	tx := r.Connection.Database.First(&value, conds...)
	return value, tx.Error
}

func (r *Repository[T]) Create(value *T) error {
	tx := r.Connection.Database.Create(value)
	return tx.Error
}

func (r *Repository[T]) Save(value *T) error {
	tx := r.Connection.Database.Save(value)
	return tx.Error
}

func (r *Repository[T]) Delete(conds ...interface{}) error {
	tx := r.Connection.Database.Delete(new(T), conds...)
	return tx.Error
}
