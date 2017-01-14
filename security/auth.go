package security

import "strings"

type Authentication interface {
	User() string

	HasGroup(group string) bool
}

type authentication struct {
	user   string
	groups map[string]interface{}
}

func (a *authentication) User() string {
	return a.user
}

func (a *authentication) HasGroup(group string) bool {
	group = strings.TrimSpace(group)
	_, ok := a.groups[group]
	return ok
}

func NewAuthentication(user string, groups ...string) Authentication {
	a := &authentication{user: user, groups: make(map[string]interface{}, 0)}
	for _, group := range groups {
		group = strings.TrimSpace(group)
		a.groups[group] = nil
	}
	return a
}
