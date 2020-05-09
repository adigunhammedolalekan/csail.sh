package http

import "testing"

func TestString(t *testing.T) {
	s := "/app/db-pg-tmp"
	t.Log(s[1:])
	t.Log(s[:1])
}