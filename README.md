# Overview

[![build](https://github.com/bzimmer/gravl/actions/workflows/build.yaml/badge.svg)](https://github.com/bzimmer/gravl)
[![codecov](https://codecov.io/gh/bzimmer/gravl/branch/master/graph/badge.svg?token=KIPOKXLNFM)](https://codecov.io/gh/bzimmer/gravl)

<img src="docs/images/gravl.png" width="150" alt="gravl logo" align="right">

**gravl** command line clients for activity-related services

## Activity clients
* [Strava](https://strava.com)
* [Cycling Analytics](https://www.cyclinganalytics.com/)
* [Ride with GPS](https://ridewithgps.com)
* [Zwift](https://zwift.com)

# Installation

## Docker

gravl is available as a Docker image on GitHub Container Registry:

```dockerfile
# Minimal image (just the binary)
FROM ghcr.io/bzimmer/gravl:latest AS gravl
FROM scratch
COPY --from=gravl /usr/bin/gravl /usr/bin/gravl
ENTRYPOINT ["/usr/bin/gravl"]

# Or with a base distro if you need shell access
FROM ghcr.io/bzimmer/gravl:latest AS gravl
FROM debian:bookworm-slim
COPY --from=gravl /usr/bin/gravl /usr/bin/gravl
```

Available image tags:
- `ghcr.io/bzimmer/gravl:latest` - latest stable release
- `ghcr.io/bzimmer/gravl:v1.2.3` - specific version

## Homebrew

```bash
brew install bzimmer/tap/gravl
```

# Documentation

* [manual](https://bzimmer.github.io/gravl/)
