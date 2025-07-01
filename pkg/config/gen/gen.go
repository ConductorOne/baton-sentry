package main

import (
	cfg "github.com/conductorone/baton-sentry/pkg/config"
	"github.com/conductorone/baton-sdk/pkg/config"
)

func main() {
	config.Generate("sentry", cfg.Config)
}
