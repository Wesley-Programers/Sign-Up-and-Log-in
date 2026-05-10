package mocks

import "github.com/stretchr/testify/mock"

type SecurityMock struct {
	mock.Mock
}

func (m *SecurityMock) GenerateToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *SecurityMock) TokenHash(token string) string {
	args := m.Called(token)
	return args.String(0)
}