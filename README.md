# Overview

[![build](https://github.com/bzimmer/gravl/actions/workflows/build.yaml/badge.svg)](https://github.com/bzimmer/gravl)
[![codecov](https://codecov.io/gh/bzimmer/gravl/branch/master/graph/badge.svg?token=KIPOKXLNFM)](https://codecov.io/gh/bzimmer/gravl)

<img src="docs/images/gravl.png" width="150" alt="gravl logo" align="right">

**gravl** package provides clients for activity-related services.

## Activity clients
* [Strava](https://strava.com)
* [Cycling Analytics](https://www.cyclinganalytics.com/)
* [Ride with GPS](https://ridewithgps.com)
* [Zwift](https://zwift.com)

# Documentation

* [manual](docs/commands.md)

## Authentication

The package has functionality to generate access and refresh tokens for both `cyclinganalytics` and `strava` by using the `oauth` command for each after acquiring the client id from the respective sites.

```sh
$ gravl strava oauth
2021-01-20T18:38:15-08:00 INF serving address=0.0.0.0:9001
```

Open a browser to http://localhost:9001 and you will be redirected to, in this case, Strava. Once you authorize the application the credentials will be provided in a json document in the browser. Copy the tokens to `env` configuration file.
