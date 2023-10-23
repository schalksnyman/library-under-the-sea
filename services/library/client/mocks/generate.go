package mockery

import (
	"github.com/stretchr/testify/mock"
	library "library-under-the-sea/services/library/domain"
)

//go:generate mockery --dir=../../domain --outpkg=mockery --output=. --case=snake --name=Client --with-expecter

var _ library.Client = (*Client)(nil)

// TSetup will assert the mock expectations once the test completes.
func (_m *Client) TSetup(t mock.TestingT, expectedCalls ...*mock.Call) *Client {
	_m.ExpectedCalls = append(_m.ExpectedCalls, expectedCalls...)
	if t, ok := t.(interface {
		mock.TestingT
		Cleanup(func())
	}); ok {
		t.Cleanup(func() { _m.AssertExpectations(t) })
	}
	return _m
}
