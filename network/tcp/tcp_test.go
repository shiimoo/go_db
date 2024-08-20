package tcp

import (
	"context"
	"testing"
)

func TestTcpListen(t *testing.T) {

	server, err := NewServer(context.Background(), "test", "0.0.0.0:8080")
	if err != nil {
		t.Error(err)
	}
	server.Start()
}
