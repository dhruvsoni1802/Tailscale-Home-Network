package server

import (
	"context"
	"log"
	"net/http"
)

// unexported type prevents context key collisions with other packages
type contextKey string

const userKey contextKey = "user"
const deviceKey contextKey = "device"

// authMiddleware is a middleware that authenticates the request using the local client
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		who, err := s.localClient.WhoIs(r.Context(), r.RemoteAddr)
		if err != nil || who == nil {
			log.Printf("auth failed for %s: %v", r.RemoteAddr, err)
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		user := who.UserProfile.LoginName
		device := who.Node.Hostinfo.Hostname

		// Attaching identity to request context so handlers can access it
		ctx := context.WithValue(r.Context(), userKey, user)
		ctx = context.WithValue(ctx, deviceKey, device)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helpers handler use to pull identity out of context
func getUser(r *http.Request) string {
	user, _ := r.Context().Value(userKey).(string)
	return user
}

// Helper handler use to pull device out of context
func getDevice(r *http.Request) string {
	device, _ := r.Context().Value(deviceKey).(string)
	return device
}