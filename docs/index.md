# gravl - command line access to activity platforms

## Installation

```shell
$ brew install bzimmer/tap/gravl
```

## Configuration

All necessary credentials are read from environment variables or command-line flags. To
create a template `.env` file, run the `envvars` command:

```sh
~/Development/src/github.com/bzimmer/gravl (oauth) > ./dist/gravl envvars
CYCLINGANALYTICS_ACCESS_TOKEN=
CYCLINGANALYTICS_CLIENT_ID=
CYCLINGANALYTICS_CLIENT_SECRET=
RWGPS_ACCESS_TOKEN=
RWGPS_CLIENT_ID=
STRAVA_CLIENT_ID=
STRAVA_CLIENT_SECRET=
STRAVA_REFRESH_TOKEN=
ZWIFT_PASSWORD=
ZWIFT_USERNAME=
```

Save these to a file, add your own credentials, and then source the file (or use whatever
environment variable mechanism suits your setup).

## Authentication

The package has functionality to generate access and refresh tokens for both
`cyclinganalytics` and `strava` by using the `oauth` command for each after acquiring the
client id from the respective sites.

```sh
~/Development/src/github.com/bzimmer/gravl (oauth) > ./dist/gravl strava oauth
2021-10-22T07:38:38-07:00 INF created strava client
2021-10-22T07:38:38-07:00 INF oauth redirect=http://localhost:9001/strava/auth/callback
2021-10-22T07:38:38-07:00 INF serving address=http://localhost:9001
```

Open a browser to http://localhost:9001 and you will be redirected to, in this case,
Strava. Once you authorize the application the credentials will be provided in a json
document in the browser. Copy the tokens to your `env` configuration file and try some
commands.

_For most commands the timeout value is reset on each query. For example, if you query 12
activities from Strava each query will honor the timeout value, it's not a deadline._

## Usage

See the [manual](https://bzimmer.github.io/gravl/commands) for an overview of all the commands.
