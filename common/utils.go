package common

import (
	"context"
	"log"
)

type contextKey int

// Context keys used to add to context.
const (
	ContextKeyRequestXRequestID contextKey = iota
	ContextKeyRequestAuthorization
)

// PrintContext is a helper function to test Context.Value().
func PrintContext(ctx context.Context) {
	if err := ctx.Value(ContextKeyRequestXRequestID); err != nil {
		log.Println("X-Request-ID", ctx.Value(ContextKeyRequestXRequestID).(string))
	}
	if err := ctx.Value(ContextKeyRequestAuthorization); err != nil {
		log.Println("Authorization", ctx.Value(ContextKeyRequestAuthorization).(string))
	}
}
