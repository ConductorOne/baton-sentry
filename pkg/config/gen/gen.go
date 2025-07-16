package main

import (
	"github.com/conductorone/baton-sdk/pkg/config"
	cfg "github.com/conductorone/baton-sentry/pkg/config"
)

func main() {
	config.Generate("sentry", cfg.Config)
}
