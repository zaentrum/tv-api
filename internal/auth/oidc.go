package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/coreos/go-oidc/v3/oidc"
)

type ctxKey int

const (
	subjectKey ctxKey = iota
)

type Verifier struct {
	enabled  bool
	audience string
	once     sync.Once
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	initErr  error
	issuer   string
}

func NewVerifier(issuer, audience string, enabled bool) *Verifier {
	return &Verifier{enabled: enabled, audience: audience, issuer: issuer}
}

func (v *Verifier) init(ctx context.Context) error {
	v.once.Do(func() {
		p, err := oidc.NewProvider(ctx, v.issuer)
		if err != nil {
			v.initErr = err
			return
		}
		v.provider = p
		v.verifier = p.Verifier(&oidc.Config{ClientID: v.audience, SkipClientIDCheck: v.audience == ""})
	})
	return v.initErr
}

func (v *Verifier) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !v.enabled {
			next.ServeHTTP(w, r)
			return
		}
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "missing bearer token", http.StatusUnauthorized)
			return
		}
		raw := strings.TrimPrefix(auth, "Bearer ")
		if err := v.init(r.Context()); err != nil {
			http.Error(w, "oidc provider unavailable", http.StatusServiceUnavailable)
			return
		}
		tok, err := v.verifier.Verify(r.Context(), raw)
		if err != nil {
			http.Error(w, "invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), subjectKey, tok.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SubjectFromContext(ctx context.Context) (string, error) {
	v, ok := ctx.Value(subjectKey).(string)
	if !ok || v == "" {
		return "", errors.New("no subject in context")
	}
	return v, nil
}
