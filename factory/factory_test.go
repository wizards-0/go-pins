package factory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var f Factory = New()

func TestGetService(t *testing.T) {
	assert := assert.New(t)

	f.RegisterFactory("service1", NewService1)
	f.RegisterFactory("service2", NewService2)
	f.RegisterFactory("service3", NewService3)
	f.RegisterFactory("service4", NewService4)
	f.RegisterFactory("service5", NewService5)

	s1, err := f.GetBean("service1")
	_, castable := s1.(Service1)
	assert.True(castable)
	assert.Nil(err)

	s2, err := f.GetBean("service2")
	_, castable = s2.(Service2)
	assert.True(castable)
	assert.Nil(err)

	_, err = f.GetBean("service3")
	assert.ErrorContains(err, "cycle")

	_, err = f.GetBean("service4")
	assert.ErrorContains(err, "cycle")

	_, err = f.GetBean("service5")
	assert.ErrorContains(err, "cycle")

	_, err = f.GetBean("service6")
	assert.ErrorContains(err, "no bean / factory found")

	_, err = GetBean("service6")
	assert.ErrorContains(err, "no bean / factory found")

	RegisterBean("service6", []string{})
	_, err = GetBean("service6")
	assert.Nil(err)

	f.RegisterBean("service6", []string{})
	_, err = f.GetBean("service6")
	assert.Nil(err)

	RegisterFactory("service7", NewService1)
	_, err = GetBean("service7")
	assert.Nil(err)
}

type Service1 interface {
	do1()
}

type service1 struct {
}

func (s *service1) do1() {
	//NOOP
}

func NewService1() (any, error) {
	return &service1{}, nil
}

type Service2 interface {
	do2()
}

type service2 struct {
	s1 Service1
}

func (s *service2) do2() {
	//NOOP
}

func NewService2() (any, error) {

	s1, _ := f.GetBean("service1")
	return &service2{
		s1: s1.(Service1),
	}, nil
}

func NewService3() (any, error) {

	_, err := f.GetBean("service5")
	return nil, err
}

func NewService4() (any, error) {

	_, err := f.GetBean("service3")
	return nil, err
}

func NewService5() (any, error) {

	_, err := f.GetBean("service4")
	return nil, err
}
