package godtp

import (
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	waitTime := 100 * time.Millisecond

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

	time.Sleep(waitTime)
	err = server.Stop()
	if err != nil {
		t.Errorf(err.Error())
	}
	if server.serving {
		t.Errorf("Server should not be serving")
	}
}
