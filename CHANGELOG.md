# Changelog

All notable changes to Naboo Email Server will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-01-04

### Added
- Initial production release
- gRPC-based email service with `SendEmail` RPC method
- SMTP client connection pooling with configurable pool size
- TLS 1.2+ encryption for all SMTP connections
- SMTP PLAIN authentication support
- HTML email support with MIME formatting
- Environment-based configuration (.env file)
- Graceful connection lifecycle management (acquire/release pattern)
- Connection health checks before reuse
- 10-second timeout protection for SMTP operations
- Docker multi-stage build with scratch base image
- Multi-platform Docker images (linux/amd64, linux/arm64)
- Automated CI/CD pipeline with GitHub Actions
- Security scanning with Trivy
- Docker Hub integration with automated description sync
- Optimized build caching for faster CI/CD

### Security
- Non-root user execution in Docker (UID 65534)
- Minimal attack surface with scratch-based container
- TLS certificate validation for SMTP connections
- Secure credential management via environment variables

### Performance
- Connection pooling reduces SMTP handshake overhead
- Optimized Docker image size (~10-15MB vs ~800MB)
- Build cache optimization with Go module layer caching
- Multi-platform builds with ARM64 support

[1.0.0]: https://github.com/lodjim/naboo-email/releases/tag/v1.0.0
