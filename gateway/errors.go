package gateway

import "fmt"

var (
	ErrTransportClosed = fmt.Errorf("transport is closed")
)
