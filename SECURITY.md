# Security Policy

## Supported Versions

Security fixes are applied to the latest released version. Please make sure you
are running the most recent release before reporting an issue.

## Reporting a Vulnerability

Please **do not** open a public GitHub issue for security vulnerabilities.

Instead, report them privately using GitHub's
[private vulnerability reporting](https://github.com/rknightion/transceiver-exporter/security/advisories/new)
("Report a vulnerability" under the repository's **Security** tab).

Please include:

- A description of the vulnerability and its impact
- Steps to reproduce (a proof of concept if possible)
- The affected version(s)

You can expect an initial acknowledgement within a few days. Once a fix is
available, a new release will be published and, where appropriate, a GitHub
Security Advisory issued.

## Scope

This exporter reads transceiver diagnostics from local network interfaces via
`ethtool` ioctls and exposes them over an unauthenticated HTTP `/metrics`
endpoint. Operators are responsible for network-level access control (e.g.
binding to a management interface or firewalling port `9458`). Reports about the
exporter's own handling of untrusted input, privilege usage, or dependency
vulnerabilities are all in scope.
