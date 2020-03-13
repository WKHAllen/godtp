package godtp

import "testing"

func TestServer(t *testing.T) {
	server := NewServer(nil, nil, nil, false, false)
	if server.serving {
		t.Errorf("Server should not be serving")
	}
}
