package security

import (
	"errors"
	"log"
)

type AuthenticationManager interface {
	Authenticate(user, secret string) (Authentication, error)
}

type AuthenticationHandler func(user, secret string) Authentication

type manager struct {
	handlers []AuthenticationHandler
}

func (m *manager) Authenticate(user, secret string) (Authentication, error) {
	for _, handler := range m.handlers {
		authentication := handler(user, secret)
		if authentication != nil {
			return authentication, nil
		}
	}
	return nil, errors.New("Invalid user name or password.")
}

func NewAuthenticationManager(handlers ...AuthenticationHandler) AuthenticationManager {
	if len(handlers) == 0 {
		e := errors.New("Attemt to create security.AuthenticationManager without AuthenticationHandler!")
		log.Printf("PANIC: %s", e)
		panic(e)
	}
	m := &manager{
		handlers: handlers,
	}
	return m
}
