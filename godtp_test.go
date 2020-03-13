package godtp

import "testing"

func TestServer(t *testing.T) {
	server := NewServer(nil, nil, nil, false, false)
	if server.serving {
		t.Errorf("Server should not be serving")
	}
	err := server.StartDefault()
	if err != nil {
		t.Errorf(err.Error())
	}
	if !server.serving {
		t.Errorf("Server should be serving")
	}
}
