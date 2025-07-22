package factory

import (
	"fmt"

	"github.com/samber/lo"
)

type Factory interface {
	GetService(name string) (any, error)
	RegisterFactory(name string, factory func() (any, error))
}

type factory struct {
	beans     map[string]any
	factories map[string]func() (any, error)
	requests  []string
}

var _factory = &factory{
	beans:     map[string]any{},
	factories: map[string]func() (any, error){},
	requests:  []string{},
}

func GetInstance() Factory {
	return _factory
}

func (f *factory) GetService(name string) (any, error) {

	if bean, exists := f.beans[name]; exists {
		return bean, nil
	} else {
		_, found := lo.Find(f.requests, func(item string) bool {
			return item == name
		})
		if found {
			cycle := append(f.requests, name)
			return nil, fmt.Errorf("cycle found in dependencies: %v", cycle)
		}
		if factory, exists := f.factories[name]; exists {
			f.requests = append(f.requests, name)
			bean, err := factory()
			if err == nil {
				f.beans[name] = bean
			}
			f.requests = f.requests[:len(f.requests)-1]
			return bean, err
		} else {
			return nil, fmt.Errorf("no factory found for id: %s", name)
		}
	}
}

func (f *factory) RegisterFactory(name string, factory func() (any, error)) {
	f.factories[name] = factory
}
