package data

import "context"

type Controller interface {
	WriteEventMessage(ctx context.Context, msg []byte) error
}
