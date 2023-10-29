package mocks

import (
	"github.com/stretchr/testify/mock"
	"net/http"
)

type MockHttpClient struct {
	mock.Mock
}

func (m *MockHttpClient) Do(request *http.Request) (*http.Response, error) {
	res := m.Called(request)
	return res.Get(0).(*http.Response), res.Error(1)
}
