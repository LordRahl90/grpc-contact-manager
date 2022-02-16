package mocks

import (
	"database/sql/driver"
	"strings"
	"time"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyPassword struct{}

func (a AnyPassword) Match(v driver.Value) bool {
	s := v.(string)
	if len(s) < 60 {
		return false
	}
	if !strings.HasPrefix(s, "$") {
		return false
	}
	return true
}
