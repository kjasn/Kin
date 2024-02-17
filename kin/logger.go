package kin

import (
	"log"
	"time"
)

// serve as debug... (maybe
func Logger() HandlerFunc {
	return func(ctx *Context) {
		// Start timer
		t := time.Now()
		// Process request
		ctx.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", ctx.StatusCode, ctx.Req.RequestURI, time.Since(t))
	}
}