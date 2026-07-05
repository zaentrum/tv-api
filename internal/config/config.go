package config

import (
	"os"
	"strings"
)

type Config struct {
	Addr         string
	OIDCIssuer   string
	OIDCAudience string
	OIDCEnabled  bool
}

func Load() Config {
	c := Config{
		Addr:         envDefault("ADDR", ":8080"),
		OIDCIssuer:   envDefault("OIDC_ISSUER", ""),
		OIDCAudience: envDefault("OIDC_AUDIENCE", "tv"),
		OIDCEnabled:  envDefault("OIDC_ENABLED", "true") != "false",
	}
	return c
}

func envDefault(key, def string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	return v
}
