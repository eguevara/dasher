package middleware

import (
	"context"
	"net/http"

	"github.com/eguevara/dasher/common"
	"github.com/pborman/uuid"
)

// AppHandler decorates the next func with context information.
func AppHandler(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

	// Lets use the existing context part of http.Request
	ctx := req.Context()

	// Lets add header information to context.
	authorization := req.Header.Get("Authorization")
	if authorization != "" {
		ctx = context.WithValue(ctx, common.ContextKeyRequestAuthorization, authorization)
	}

	requestID := req.Header.Get("X-Request-Id")
	if requestID == "" {
		requestID = uuid.New()
	}
	ctx = context.WithValue(ctx, common.ContextKeyRequestXRequestID, requestID)

	// Create a new context to call with next.
	req = req.WithContext(ctx)

	// Adding x-request-id to response header.
	res.Header().Set("X-Request-Id", requestID)

	// Call the next http handler with new req including context.
	next(res, req)
}
