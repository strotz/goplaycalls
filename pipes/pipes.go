package pipes

import (
	"context"
	"net"
	"os"
	"path/filepath"
)

const unixProtocol = "unix"

// CreateListener creates a new named pipe and attaches listener to it.
func CreateListener(filePath string) (net.Listener, error) {
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	return net.Listen(unixProtocol, filePath)
}

type DialerFunc func(ctx context.Context, network, addr string) (net.Conn, error)

func CreateDialer(filePath string) DialerFunc {
	return func(_ context.Context, _, _ string) (net.Conn, error) {
		return net.Dial(unixProtocol, filePath)
	}
}
