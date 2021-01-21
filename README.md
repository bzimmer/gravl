# Overview

<img src="docs/images/gravl.png" width="150" alt="gravl logo" align="right">

**gravl** package provides clients for activty-related services and an extensible analysis framework for activities.

The purpose of the package is to provide easy access to activity, weather, and geo services useful for either planning or analyzing activities. The package is split into a few top level compenents: `providers`, `analysis`, and `commands`.

The `providers` package is responsible for communicating with services and aims to use a consistent approach to APIs and models.

The `analysis` package is responsible for running analyzers on Strava data.

The `commands` package contains all the commands for the cli.

## Activity clients
* [Strava](https://strava.com)
* [Cycling Analytics](https://www.cyclinganalytics.com/)
* [Ride with GPS](https://ridewithgps.com)
* [WTA](https://wta.org)

## Geo
* [GNIS](https://geonames.usgs.gov)
* [SRTM](https://github.com/sakisds/go-srtm)

## Weather
* [NOAA](https://weather.gov)
* [OpenWeather API](https://openweathermap.org/api)
* [VisualCrossing](https://visualcrossing.com)

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

The analysis command supports flexible filtering and grouping of activities using the [expr](https://github.com/antonmedv/expr) package to evaluate the Strava [Activity](https://github.com/bzimmer/gravl/blob/master/pkg/providers/activity/strava/model.go#L333) model.

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

Export an activity data file from Strava using the web client, upload it to [Cycling Analytics](https://www.cyclinganalytics.com/), and poll the status until processing is completed.

```sh
~ > gravl strava export 4612178259
2021-01-12T20:22:13-08:00 INF export activityID=4612178259 format=original
"Innsbruck.fit"
~ > gravl ca upload -p Innsbruck.fit
2021-01-12T20:23:12-08:00 INF uploading file=Innsbruck.fit size=112732
{
 "status": "processing",
 "ride_id": 0,
 "user_id": 1603533,
 "format": "fit",
 "datetime": "2021-01-13T04:23:15",
 "upload_id": 4775060590,
 "filename": "Innsbruck.fit",
 "size": 112732,
 "error": "",
 "error_code": ""
}
{
 "status": "done",
 "ride_id": 382207409453,
 "user_id": 1603533,
 "format": "fit",
 "datetime": "2021-01-13T04:23:15",
 "upload_id": 4775060590,
 "filename": "Innsbruck.fit",
 "size": 112732,
 "error": "",
 "error_code": ""
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
