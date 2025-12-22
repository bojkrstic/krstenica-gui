package requestctx

import "context"

// userContextKey isolates the context value for authenticated users.
type userContextKey struct{}

// User carries authenticated identity data through request handling layers.
type User struct {
	ID       int64
	Username string
	Role     string
	City     string
}

// WithUser attaches the authenticated user to the context.
func WithUser(ctx context.Context, user *User) context.Context {
	if ctx == nil || user == nil {
		return ctx
	}
	return context.WithValue(ctx, userContextKey{}, user)
}

// UserFromContext extracts the authenticated user from context if present.
func UserFromContext(ctx context.Context) (*User, bool) {
	if ctx == nil {
		return nil, false
	}
	if user, ok := ctx.Value(userContextKey{}).(*User); ok && user != nil {
		return user, true
	}
	return nil, false
}

// IsAdmin returns true when the user has admin privileges.
func (u *User) IsAdmin() bool {
	if u == nil {
		return false
	}
	switch u.Role {
	case "admin", "ADMIN", "Admin":
		return true
	}
	return false
}
