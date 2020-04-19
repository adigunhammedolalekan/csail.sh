package fn

import "testing"

func TestFn(t *testing.T) {
	value := "v123"
	v := value[1:]
	if v != "123" { t.Fatal() }
}