# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | ✅ Yes             |
| < 1.0   | ❌ No              |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability in TormentNexus, please report it responsibly.

### How to Report

1. **DO NOT** open a public GitHub issue
2. **Email** <security@tormentnexus.org> with:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### What to Expect

- **Acknowledgment** within 24 hours
- **Status update** within 72 hours
- **Fix timeline** within 1 week
- **Credit** in the security advisory (unless you prefer anonymity)

### Bug Bounty

We're working on establishing a bug bounty program. In the meantime, we'll acknowledge security researchers in our release notes and send swag for critical findings.

## Security Best Practices

### For Users

1. **Keep TormentNexus updated** to the latest version
2. **Use strong passwords** for your accounts
3. **Enable 2FA** where available
4. **Don't expose** the dashboard to the public internet without authentication
5. **Use HTTPS** in production
6. **Regular backups** of your data

### For Developers

1. **Never commit** secrets, API keys, or passwords
2. **Use environment variables** for sensitive configuration
3. **Validate all inputs** from users
4. **Use parameterized queries** to prevent SQL injection
5. **Keep dependencies updated** and audit for vulnerabilities
6. **Follow secure coding practices**

## Security Features

### Current

- **Local-first** — No data leaves your machine unless you configure it
- **SQLite** — Database is local, not in the cloud
- **No telemetry** — We don't track you
- **Open source** — Code is auditable

### Planned

- **SSO/SAML** — Enterprise authentication
- **RBAC** — Role-based access control
- **Audit logs** — Track all actions
- **Encryption at rest** — Encrypt stored data
- **Encryption in transit** — TLS for all connections
- **Security scanning** — Automated vulnerability detection

## Known Issues

No known security issues at this time.

## Security Advisories

Security advisories will be published at:
<https://github.com/MDMAtk/TormentNexus/security/advisories>

## Contact

- **Security email:** <security@tormentnexus.org>
- **PGP key:** (coming soon)

## Acknowledgments

We thank the following security researchers for their responsible disclosures:

(None yet — be the first!)
