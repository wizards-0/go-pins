package factory

import (
	"fmt"

	"github.com/samber/lo"
)

type Factory interface {
	GetBean(id string) (any, error)
	RegisterFactory(id string, factory func() (any, error))
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

func (f *factory) GetBean(id string) (any, error) {

	if bean, exists := f.beans[id]; exists {
		return bean, nil
	} else {
		_, found := lo.Find(f.requests, func(item string) bool {
			return item == id
		})
		if found {
			cycle := append(f.requests, id)
			return nil, fmt.Errorf("cycle found in dependencies: %v", cycle)
		}
		if factory, exists := f.factories[id]; exists {
			f.requests = append(f.requests, id)
			bean, err := factory()
			if err == nil {
				f.beans[id] = bean
			}
			f.requests = f.requests[:len(f.requests)-1]
			return bean, err
		} else {
			return nil, fmt.Errorf("no factory found for id: %s", id)
		}
	}
}

func (f *factory) RegisterFactory(id string, factory func() (any, error)) {
	f.factories[id] = factory
}
