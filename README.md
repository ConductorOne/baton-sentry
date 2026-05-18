![Baton Logo](./baton-logo.png)

# `baton-sentry` [![Go Reference](https://pkg.go.dev/badge/github.com/conductorone/baton-sentry.svg)](https://pkg.go.dev/github.com/conductorone/baton-sentry) ![ci](https://github.com/conductorone/baton-sentry/actions/workflows/ci.yaml/badge.svg) ![verify](https://github.com/conductorone/baton-sentry/actions/workflows/verify.yaml/badge.svg)

`baton-sentry` is a connector for built using the [Baton SDK](https://github.com/conductorone/baton-sdk).

Check out [Baton](https://github.com/conductorone/baton) to learn more the project in general.

# Prerequisites
No prerequisites were specified for `baton-sentry`

# Getting Started

## brew

```
brew install conductorone/baton/baton conductorone/baton/baton-sentry
baton-sentry
baton resources
```

## docker

```
docker run --rm -v $(pwd):/out -e BATON_API_TOKEN=apiToken ghcr.io/conductorone/baton-sentry:latest -f "/out/sync.c1z"
docker run --rm -v $(pwd):/out ghcr.io/conductorone/baton:latest -f "/out/sync.c1z" resources
```

## source

```
go install github.com/conductorone/baton/cmd/baton@main
go install github.com/conductorone/baton-sentry/cmd/baton-sentry@main

baton-sentry

baton resources
```

# Data Model

`baton-sentry` will pull down information about the following resources:

| Resource Type  | Description |
|---------------|-------------|
| Organization  | Sentry organizations. Each organization exposes role-based entitlements (billing, member, admin, manager, owner). |
| User          | Organization members, including pending invites. |
| Team          | Teams within an organization. Exposes a `member` assignment entitlement. |
| Project       | Projects within an organization. Exposes an `assigned` entitlement granted to teams. |

## Provisioning

`baton-sentry` supports the following provisioning capabilities:

| Action | Resource | Description |
|--------|----------|-------------|
| Create Account | User | Invite a new member to an organization. |
| Delete Account | User | Remove a member from an organization. |
| Grant / Revoke | Organization Role | Change a member's organization role (billing, member, admin, manager, owner). Revoking a non-default role downgrades the member to the `member` role. |
| Grant / Revoke | Team Membership | Add or remove a member from a team. |
| Grant / Revoke | Project Assignment | Add or remove a team from a project. |

# Contributing, Support and Issues

We started Baton because we were tired of taking screenshots and manually
building spreadsheets. We welcome contributions, and ideas, no matter how
small&mdash;our goal is to make identity and permissions sprawl less painful for
everyone. If you have questions, problems, or ideas: Please open a GitHub Issue!

See [CONTRIBUTING.md](https://github.com/ConductorOne/baton/blob/main/CONTRIBUTING.md) for more details.

# `baton-sentry` Command Line Usage

```
baton-sentry

Usage:
  baton-sentry [flags]
  baton-sentry [command]

Available Commands:
  capabilities       Get connector capabilities
  completion         Generate the autocompletion script for the specified shell
  help               Help about any command

Flags:
      --api-token string             API Token for Sentry ($BATON_API_TOKEN)
      --client-id string             The client ID used to authenticate with ConductorOne ($BATON_CLIENT_ID)
      --client-secret string         The client secret used to authenticate with ConductorOne ($BATON_CLIENT_SECRET)
  -f, --file string                  The path to the c1z file to sync with ($BATON_FILE) (default "sync.c1z")
  -h, --help                         help for baton-sentry
      --log-format string            The output format for logs: json, console ($BATON_LOG_FORMAT) (default "json")
      --log-level string             The log level: debug, info, warn, error ($BATON_LOG_LEVEL) (default "info")
      --org-ids strings              Limit syncing to the specified Sentry organization IDs or slugs. If empty, all organizations are synced. ($BATON_ORG_IDS)
  -p, --provisioning                 If this connector supports provisioning, this must be set in order for provisioning actions to be enabled ($BATON_PROVISIONING)
      --ticketing                    This must be set to enable ticketing support ($BATON_TICKETING)
  -v, --version                      version for baton-sentry

Use "baton-sentry [command] --help" for more information about a command.
```
