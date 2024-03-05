package scope

type Scope int8

const (
	ScopeAdmin Scope = iota
	ScopeUser
)

// String returns the string representation of the Scope value.
func (s Scope) String() string {
	return [...]string{"admin", "user"}[s]
}

// IsValid checks if the scope is valid.
// It returns true if the scope is either ScopeAdmin or ScopeUser, otherwise it returns false.
func (s Scope) IsValid() bool {
	switch s {
	case ScopeAdmin, ScopeUser:
		return true
	}
	return false
}

// Has checks if the current scope has any of the specified scopes.
// It returns true if the current scope matches any of the specified scopes, and false otherwise.
func (s Scope) Has(scopes ...Scope) bool {
	for _, scope := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}

// Parse parses the given string and returns the corresponding Scope value.
// If the string does not match any known scope, it returns -1.
func Parse(s string) Scope {
	switch s {
	case "admin":
		return ScopeAdmin
	case "user":
		return ScopeUser
	}
	return -1
}
