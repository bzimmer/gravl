# Overview

<img src="docs/images/gravl.png" width="150" alt="gravl logo" align="right">

**gravl** package provides clients for activty-related services and an extensible analysis framework for activities.

The purpose of the package is to provide easy access to activity, weather, and geo services useful for either planning or analyzing activities. The package is split into a few top level compenents:
* `providers`
  * a library for communicating with services and aims to use a consistent approach to APIs and models
* `analysis`
  * a library for running [analyzers](docs/analyzers.md) on Strava activities
* `store`
  * a library for storing Strava activity data, generally locally through `buntdb` or a file but also capable of interacting with Strava directly
* `eval`
  * a flexible evaluation library useful for dynamic filtering, grouping, and evaluating of Strava activities
* `commands`
  * the commands used by the cli

More documentation and numerous examples can be found in the [manual](docs/manual.md).

## Activity clients
* [Strava](https://strava.com)
* [Cycling Analytics](https://www.cyclinganalytics.com/)
* [Ride with GPS](https://ridewithgps.com)
* [WTA](https://wta.org)
* [Zwift](https://zwift.com)

## Weather
* [NOAA](https://weather.gov)
* [OpenWeather API](https://openweathermap.org/api)
* [VisualCrossing](https://visualcrossing.com)

## Geo
* [GNIS](https://geonames.usgs.gov)
* [SRTM](https://github.com/sakisds/go-srtm)

# Documentation

## Configuration

The configuration file contains the credentials for all services.

```yaml
origin: http://localhost

visualcrossing:
  access-token:       {{access-token}}

rwgps:
  client-id:          {{client-id}}
  access-token:       {{access-token}}

strava:
  client-id:          {{client-id}}
  client-secret:      {{client-secret}}

  access-token:       {{access-token}}
  refresh-token:      {{refresh-token}}

  username:           {{email}}
  password:           {{password}}

openweather:
  access-token:       {{access-token}}

cyclinganalytics:
  client-id:          {{client-id}}
  client-secret:      {{client-secret}}
  access-token:       {{access-token}}
```

## Authentication

The package has functionality to generate access and refresh tokens for both `cyclinganalytics` and `strava` by using the `oauth` command for each after acquiring the client id from the respective sites.

```sh
~ > gravl strava oauth
2021-01-20T18:38:15-08:00 INF serving address=0.0.0.0:9001
```

Open a browser to http://localhost:9001 and you will be redirected to, in this case, Strava. Once you authorize the application the credentials will be provided in a json document in the browser. Copy the tokens to the configuration file.

```json
{
  "access_token": "{{access-token}}",
  "token_type": "Bearer",
  "refresh_token": "{{refresh-token}}",
  "expiry": "2021-01-20T22:16:56.243892-08:00"
}
```

## Analysis

The analysis command supports flexible filtering and grouping of activities using the [expr](https://github.com/antonmedv/expr) package to evaluate the Strava Activity model.


* Filters

As an example, to filter only start dates for 2021:

```sh
~ > gravl pass -a totals -f ".StartDate.Year() == 2021"
```

* Groups

Each of the group expressions will result in a new level in the output of the analysis:

```sh
~ > gravl pass -a totals -f ".StartDate.Year() == 2021" -g "isoweek(.StartDate)" -g ".Type"
```

## Examples

The following bash script will export an activity data file from Strava using the web client, upload it to [Cycling Analytics](https://www.cyclinganalytics.com/), and poll the status until processing is completed.

```sh
#!/bin/bash
set -e

num="${NUM_ACTIVITIES:-50}"

function _jq() {
    if ! command -v jq &> /dev/null
    then
        cat "$@"
        exit
    fi
    jq -sc ".[]" "$@"
}

# If no activity is provided display the most recent rides
if [[ $# -eq 0 ]]
then
    gravl -c --timeout 1m strava activities -N $num -f ".Type == 'VirtualRide'" -B ".ID, .Name, .StartDateLocal, .Distance.Miles()" | _jq
    exit 0
fi

for arg in "$@"
do
    gravl strava export -o -T "$arg.fit" -F fit $arg
    gravl cyclinganalytics upload -p "$arg.fit"
    rm -f "$arg.fit"
done
```

When executed the command output will look something like this:

```sh
~ > qp
2021-01-28T20:39:50-08:00 INF do all=50 count=50 n=50 start=1 total=50
[4687554641,"Innsbruck","2021-01-26T18:15:29Z"]
[4612178259,"Innsbruck","2021-01-12T18:40:56Z"]
[4569050661,"Paris","2021-01-04T19:21:36Z"]
[4481763454,"Watopia","2020-12-16T19:04:41Z"]

~ > qp 4334103705
2021-01-28T20:37:16-08:00 INF export activityID=4334103705 format=original
{
 "id": 4334103705,
 "name": "4334103705.fit",
 "format": "original",
 "ext": "fit"
}
2021-01-28T20:37:17-08:00 INF uploading file=4334103705.fit size=101003
{
 "status": "processing",
 "ride_id": 0,
 "user_id": 1603533,
 "format": "fit",
 "datetime": "2021-01-29T04:37:20",
 "upload_id": 1394198469,
 "filename": "4334103705.fit",
 "size": 101003,
 "error": "",
 "error_code": ""
}
{
 "status": "error",
 "ride_id": 0,
 "user_id": 1603533,
 "format": "fit",
 "datetime": "2021-01-29T04:37:20",
 "upload_id": 1394198469,
 "filename": "4334103705.fit",
 "size": 101003,
 "error": "The ride already exists: 260297069518",
 "error_code": "duplicate_ride"
}
```

Run the `totals` analyzer for the year 2021 by specifying a filter.

```sh
~ > gravl pass -a totals -f ".StartDate.Year() == 2021"
{
 "": {
  "totals": {
   "count": 21,
   "distance": 272.9107636403404,
   "elevation": 21794.619422572177,
   "movingtime": 88083,
   "centuries": {
    "metric": 0,
    "imperial": 0
   }
  }
 }
}
```
