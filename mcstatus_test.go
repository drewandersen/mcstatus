package mcstatus

import "testing"

func TestNewNoServer(t *testing.T) {
	config := Config{
		Addr: "",
	}

	_, err := New(config)
	if err == nil {
		t.Errorf("expected no server address to fail")
	}

}
