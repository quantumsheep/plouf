package plouf

type IInjectable interface {
	Init(self IInjectable) error

	ShouldLogInjection(self IInjectable) bool
	IsInitialized() bool
}

type Injectable struct {
	initialized bool
}

func (i *Injectable) IsInitialized() bool {
	return i.initialized
}

func (i *Injectable) Init(self IInjectable) error {
	if i.IsInitialized() {
		return nil
	}

	i.initialized = true

	value := ReflectValue(self)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)

		if !field.CanInterface() {
			continue
		}

		if injectable, ok := field.Interface().(IInjectable); ok {
			if injectable.IsInitialized() {
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
