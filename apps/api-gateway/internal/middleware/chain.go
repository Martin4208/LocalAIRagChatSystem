// ミドルウェア結合
package middleware

import "net/http"

// Chain creates a middleware chain by wrapping handlers in reverse order
// The first middleware in the slice will be the outermost wrapper
func Chain(middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		// Apply middlewares in reverse order
		// This ensures the first middleware wraps all others
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i](final)
		}
		return final
	}
}

// Example usage:
// stack := middleware.Chain(
//     middleware.Recovery,  // Outermost - catches panics from everything
//     middleware.Logger,    // Logs all requests
//     middleware.CORS,      // Handles CORS headers
// )
// router.Use(stack)
