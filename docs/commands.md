# gravl - Activity related analysis, exploration, & planning


The top level `gravl` command has some useful flags:

* Use `--http-tracing` to enable verbose logging of http requests for debugging.
* Use `--timeout DURATION` to specify a timeout duration (eg `1m`, `10s`)

For most commands the timeout value is reset on each query. For example, if you query 12 activities
from Strava each query will honor the timeout value, it's not an aggregate timeout.

Some commands, such as [store update](#store-update) will require a timeout longer than the default
since the operation can take a long time.


## Commands
* [analysis](#analysis)
* [analysis list](#analysis-list)
* [commands](#commands)
* [cyclinganalytics](#cyclinganalytics)
* [cyclinganalytics activities](#cyclinganalytics-activities)
* [cyclinganalytics activity](#cyclinganalytics-activity)
* [cyclinganalytics athlete](#cyclinganalytics-athlete)
* [cyclinganalytics oauth](#cyclinganalytics-oauth)
* [cyclinganalytics streamsets](#cyclinganalytics-streamsets)
* [cyclinganalytics upload](#cyclinganalytics-upload)
* [gnis](#gnis)
* [gpx](#gpx)
* [gpx info](#gpx-info)
* [help](#help)
* [noaa](#noaa)
* [noaa forecast](#noaa-forecast)
* [openweather](#openweather)
* [openweather forecast](#openweather-forecast)
* [rwgps](#rwgps)
* [rwgps activities](#rwgps-activities)
* [rwgps activity](#rwgps-activity)
* [rwgps athlete](#rwgps-athlete)
* [rwgps route](#rwgps-route)
* [rwgps routes](#rwgps-routes)
* [srtm](#srtm)
* [store](#store)
* [store export](#store-export)
* [store remove](#store-remove)
* [store update](#store-update)
* [strava](#strava)
* [strava activities](#strava-activities)
* [strava activity](#strava-activity)
* [strava athlete](#strava-athlete)
* [strava export](#strava-export)
* [strava fitness](#strava-fitness)
* [strava oauth](#strava-oauth)
* [strava refresh](#strava-refresh)
* [strava route](#strava-route)
* [strava routes](#strava-routes)
* [strava streams](#strava-streams)
* [strava streamsets](#strava-streamsets)
* [strava upload](#strava-upload)
* [strava webhook](#strava-webhook)
* [strava webhook list](#strava-webhook-list)
* [strava webhook subscribe](#strava-webhook-subscribe)
* [strava webhook unsubscribe](#strava-webhook-unsubscribe)
* [version](#version)
* [visualcrossing](#visualcrossing)
* [visualcrossing forecast](#visualcrossing-forecast)
* [wta](#wta)
* [zwift](#zwift)
* [zwift activities](#zwift-activities)
* [zwift activity](#zwift-activity)
* [zwift athlete](#zwift-athlete)
* [zwift export](#zwift-export)
* [zwift files](#zwift-files)
* [zwift refresh](#zwift-refresh)

## *analysis*

**Description**

Produce statistics and other interesting artifacts from Strava activities


**Syntax**

```sh
$ gravl analysis [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```units```|```u```|Units|
|```filter```|```f```|Expression for filtering activities to remove|
|```group```|```g```|Expressions for grouping activities|
|```analyzer```|```a```|Analyzers to include (if none specified, default set is used)|
|```input```|```i```|Input data store|

**Example**

Run the analysis on the Strava activities. Learn more about analyzers [here](analyzers.md).

```sh
$ gravl pass -a totals -f ".StartDate.Year() == 2021"
{
 "": {
  "totals": {
   "count": 50,
   "distance": 592.2288211842838,
   "elevation": 47769.68503937008,
   "calories": 34076.5,
   "movingtime": 192936,
   "centuries": {
    "metric": 1,
    "imperial": 0
   }
  }
 }
}
```

In addition to filtering, it's often useful to group activities and perform analysis on sub-groups.
In this example the year is filtered and then totals are computed per type.

```sh
$ gravl pass -a totals -f ".StartDate.Year() == 2021" -g ".Type"
{
 "NordicSki": {
  "totals": {
   "count": 4,
   "distance": 27.716013481269382,
   "elevation": 2030.8398950131236,
   "calories": 1984,
   "movingtime": 19512,
   "centuries": {
    "metric": 0,
    "imperial": 0
   }
  }
 },
 "Ride": {
  "totals": {
   "count": 19,
   "distance": 433.0711768273283,
   "elevation": 34855.64304461942,
   "calories": 23870.5,
   "movingtime": 105859,
   "centuries": {
    "metric": 1,
    "imperial": 0
   }
  }
 },
 ...,
 "Walk": {
  "totals": {
   "count": 19,
   "distance": 41.42433190169411,
   "elevation": 4248.687664041995,
   "calories": 4170,
   "movingtime": 46374,
   "centuries": {
    "metric": 0,
    "imperial": 0
   }
  }
 }
}
```


## *analysis list*

**Description**

Return the list of available analyzers


**Syntax**

```sh
$ gravl analysis list [flags]
```


**Example**

List all the available analyzers.

```sh
$ gravl pass list
{
	"ageride": {
		"base": false,
		"doc": "ageride returns all activities whose distance is greater than the athlete's age at the time of the activity",
		"flags": true
	},
	...,
	"cluster": {
		"base": false,
		"doc": "clusters returns the activities clustered by (distance, elevation) dimensions",
		"flags": true
	},
	"eddington": {
		"base": true,
		"doc": "eddington returns the Eddington number for all activities - The Eddington is the largest integer E, where you have cycled at least E miles (or kilometers) on at least E days",
		"flags": false
	},
	"festive500": {
		"base": true,
		"doc": "festive500 returns the activities and distance ridden during the annual #festive500 challenge - Thanks Rapha! https://www.rapha.cc/us/en_US/stories/festive-500",
		"flags": false
	},
	...,
	"totals": {
		"base": true,
		"doc": "totals returns the number of centuries (100 mi or 100 km)",
		"flags": false
	}
}
```


## *commands*

**Description**

Return all possible commands


**Syntax**

```sh
$ gravl commands [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```relative```|```r```|Specify the command relative to the current working directory|


## *cyclinganalytics*

**Description**

Operations supported by the Cycling Analytics website


**Flags**

|Name|Aliases|Description|
|-|-|-|
|```cyclinganalytics.client-id```||API key for Cycling Analytics API|
|```cyclinganalytics.client-secret```||API secret for Cycling Analytics API|
|```cyclinganalytics.access-token```||Access token for Cycling Analytics API|


## *cyclinganalytics activities*

**Description**

Query activities for the authenticated athlete


**Syntax**

```sh
$ gravl cyclinganalytics activities [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```count```|```N```|The number of activities to query from CA (the number returned will be <= N)|


## *cyclinganalytics activity*

**Description**

Query an activity for the authenticated athlete


**Syntax**

```sh
$ gravl cyclinganalytics activity [flags]
```



## *cyclinganalytics athlete*

**Description**

Query for the authenticated athlete


**Syntax**

```sh
$ gravl cyclinganalytics athlete [flags]
```


**Example**

```sh
$ gravl ca t
{
 "email": "me@example.com",
 "id": 88827722,
 "name": "That Guy",
 "sex": "male",
 "timezone": "America/Los_Angeles",
 "units": "us"
}
```


## *cyclinganalytics oauth*

**Description**

Authentication endpoints for access and refresh tokens


**Syntax**

```sh
$ gravl cyclinganalytics oauth [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```origin```||Callback origin|
|```port```||Port on which to listen|


## *cyclinganalytics streamsets*

**Description**

Return the set of available streams for query


**Syntax**

```sh
$ gravl cyclinganalytics streamsets [flags]
```



## *cyclinganalytics upload*

**Description**

Upload an activity file


**Syntax**

```sh
$ gravl cyclinganalytics upload [flags] {FILE | DIRECTORY}
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```status```|```s```|Check the status of the upload|
|```poll```|```p```|Continually check the status of the request until it is completed|
|```dryrun```|```n```|Show the files which would be uploaded but do not upload them|


## *gnis*

**Description**

Query the GNIS database


**Syntax**

```sh
$ gravl gnis [flags] US-STATE-TWO-LETTER-ABBREVIATION
```


**Example**

Query the [GNIS](https://www.usgs.gov/core-science-systems/national-geospatial-program/geographic-names) database for US States.

This functionality was added mainly to get pseudo-accurate coordinates for use in querying weather forecasts.

```sh
$ gravl -c gnis NH | wc -l
   14770
$ gravl -c gnis NH | head -5
{"type":"Feature","id":"205110","geometry":{"type":"Point","coordinates":[-77.0775473,40.3221113,200]},"properties":{"class":"Trail","locale":"PA","name":"North Country National Scenic Trail","source":"https://geonames.usgs.gov"}}
{"type":"Feature","id":"206425","geometry":{"type":"Point","coordinates":[-72.3331382,41.2723203,0]},"properties":{"class":"Stream","locale":"CT","name":"Connecticut River","source":"https://geonames.usgs.gov"}}
{"type":"Feature","id":"561049","geometry":{"type":"Point","coordinates":[-71.0306287,44.9383838,384]},"properties":{"class":"Stream","locale":"ME","name":"Abbott Brook","source":"https://geonames.usgs.gov"}}
{"type":"Feature","id":"561428","geometry":{"type":"Point","coordinates":[-70.7517197,43.0823107,9]},"properties":{"class":"Island","locale":"ME","name":"Badgers Island","source":"https://geonames.usgs.gov"}}
{"type":"Feature","id":"561491","geometry":{"type":"Point","coordinates":[-70.9716939,43.6157566,168]},"properties":{"class":"Reservoir","locale":"ME","name":"Balch Pond","source":"https://geonames.usgs.gov"}}
```

I might typically use it like this:

```sh
$ gravl -c gnis WA | grep "Barlow Pass"
{"type":"Feature","id":"1516141","geometry":{"type":"Point","coordinates":[-121.4440005,48.0264959,721]},"properties":{"class":"Gap","locale":"WA","name":"Barlow Pass","source":"https://geonames.usgs.gov"}}
$ gravl -c gnis WA | grep "Barlow Pass" | jq ".geometry | .coordinates"
[
  -121.4440005,
  48.0264959,
  721
]
```


## *gpx*

**Description**

gpx




## *gpx info*

**Description**

Return basic statistics about a GPX file


**Syntax**

```sh
$ gravl gpx info [flags] GPX_FILE (...)
```


**Example**

A simple utility for summarizing gpx data files.

_Distance units are metric, time is seconds_

```sh
$ gravl gpx info pkg/commands/geo/gpx/testdata/2017-07-13-TdF-Stage18.gpx
{
 "filename": "pkg/commands/geo/gpx/testdata/2017-07-13-TdF-Stage18.gpx",
 "tracks": 1,
 "segments": 1,
 "points": 396,
 "distance2d": 180993.07498903852,
 "distance3d": 181154.40869605148,
 "ascent": 2310.6341268663095,
 "descent": 1881.945114464303,
 "start_time": "2017-07-13T06:00:00Z",
 "moving_time": 23583,
 "stopped_time": 12525
}
```


## *help*

**Description**

Shows a list of commands or help for one command


**Syntax**

```sh
$ gravl help [flags] [command]
```



## *noaa*

**Description**

Query NOAA for forecasts




## *noaa forecast*

**Description**

Query NOAA for a forecast


**Syntax**

```sh
$ gravl noaa forecast [flags] [--] LATITUDE LONGITUDE
```



## *openweather*

**Description**

Query OpenWeather for forecasts


**Flags**

|Name|Aliases|Description|
|-|-|-|
|```openweather.access-token```||API key for OpenWeather API|


## *openweather forecast*

**Description**

Query OpenWeather for a forecast


**Syntax**

```sh
$ gravl openweather forecast [flags] [--] LATITUDE LONGITUDE
```


**Example**

Query [OpenWeather](https://openweathermap.org/) for a forecast

```sh
$ gravl openweather forecast -- 48.8 -128.0
{
 "lat": 48.8,
 "lon": -128,
 "timezone": "Etc/GMT+9",
 "timezone_offset": -32400,
 "current": {
  "dt": 1613843684,
  "sunrise": 1613835032,
  "sunset": 1613872869,
  "temp": 7.19,
  "feels_like": 1.02,
  "pressure": 1023,
  "humidity": 70,
  "dew_point": 2.09,
  "uvi": 0.68,
  "clouds": 91,
  "visibility": 10000,
  "wind_speed": 6.45,
  "wind_deg": 252,
  "wind_gust": 0,
  "weather": [
   {
    "id": 804,
    "main": "Clouds",
    "description": "overcast clouds",
    "icon": "04d"
   }
  ]
 },
 ...
}
```


## *rwgps*

**Description**

Query RideWithGPS for rides and routes


**Flags**

|Name|Aliases|Description|
|-|-|-|
|```rwgps.client-id```||Client ID for RideWithGPS API|
|```rwgps.access-token```||Access token for RideWithGPS API|


## *rwgps activities*

**Description**

Query activities for the authenticated athlete


**Syntax**

```sh
$ gravl rwgps activities [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```count```|```N```|The number of activities to query from RideWithGPS (the number returned will be <= N)|

**Example**


Query RideWithGPS activities with an optional count. The order of the activities is not guaranteed but generally they are returned most recent first.

```sh
$ gravl rwgps activities -N 10
2021-02-20T14:31:23-08:00 INF do all=0 count=10 n=0 start=0 total=10
2021-02-20T14:31:23-08:00 INF do all=10 count=10 n=10 start=1 total=10
{
 "created_at": "2021-02-20T20:12:49Z",
 "departed_at": "2021-02-20T18:39:22Z",
 "description": "",
 "distance": 30741.9,
 "duration": 5593,
 "elevation_gain": 516.281,
 "elevation_loss": 521.453,
 "id": 62829442,
 "name": "2021/02/20",
 "type": "",
 "track_id": "60316d406b34d7693a3f015b",
 "updated_at": "2021-02-20T20:12:49Z",
 "user_id": 836,
 "visibility": 1,
 "first_lat": 47.502655,
 "first_lng": -122.602798,
 "last_lat": 47.502655,
 "last_lng": -122.602798
}
...
{
 "created_at": "2021-01-26T00:38:58Z",
 "departed_at": "2021-01-25T23:34:32Z",
 "description": "",
 "distance": 25356,
 "duration": 3851,
 "elevation_gain": 403.445,
 "elevation_loss": 400.121,
 "id": 61904776,
 "name": "2021/01/25",
 "type": "",
 "track_id": "600f64a2fa348e8e23e2d546",
 "updated_at": "2021-01-26T00:38:58Z",
 "user_id": 836,
 "visibility": 1,
 "first_lat": 47.502655,
 "first_lng": -122.602798,
 "last_lat": 47.502655,
 "last_lng": -122.602798
}
```


## *rwgps activity*

**Description**

Query an activity from RideWithGPS


**Syntax**

```sh
$ gravl rwgps activity [flags]
```



## *rwgps athlete*

**Description**

Query for the authenticated athlete


**Syntax**

```sh
$ gravl rwgps athlete [flags]
```



## *rwgps route*

**Description**

Query a route from RideWithGPS


**Syntax**

```sh
$ gravl rwgps route [flags]
```



## *rwgps routes*

**Description**

Query routes for an athlete from RideWithGPS


**Syntax**

```sh
$ gravl rwgps routes [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```count```|```N```|The number of routes to query from RideWithGPS (the number returned will be <= N)|


## *srtm*

**Description**

Query the SRTM database for elevation data


**Syntax**

```sh
$ gravl srtm [flags]
```



## *store*

**Description**

Manage a local store of Strava activities




## *store export*

**Description**

Export activities from local storage


**Syntax**

```sh
$ gravl store export [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```input```|```i```|Input data store|
|```filter```|```f```|Expression for filtering activities|
|```attribute```|```B```|Evaluate the expression on an activity and return only those results|

**Example**

`gravl` allows flexible exporting of Strava activities from the local store by the `export` command. As an example
of exporting a subset of activities:

```sh
$ gravl -c store export -f ".Type == 'NordicSki'" | wc -l
2021-02-20T19:12:07-08:00 INF bunt db path="/Users/bzimmer/Library/Application Support/com.github.bzimmer.gravl/gravl.db"
2021-02-20T19:12:08-08:00 INF export activities=46 elapsed=678.191786
46
```

It's also possible to use the attribute functionality by specifying one or more attributes using the `-B` flag. In this
example we export only those activities of type `Ride`, extract their distance in miles, and use standard unix tools to
create the top 10 rides by distance.

```sh
$ gravl -c store export -f ".Type == 'Ride'" -B ".Distance.Miles()" | jq ".[]" | sort -nr | head -10
2021-02-20T18:56:27-08:00 INF bunt db path="/Users/bzimmer/Library/Application Support/com.github.bzimmer.gravl/gravl.db"
2021-02-20T18:56:28-08:00 INF export activities=618 elapsed=682.669242
161.20357114451602
107.90359301678198
101.80421339378032
99.08571442774199
97.57764654418197
90.31630279169649
84.45552970651396
83.14567923327766
79.88223773164718
76.653593016782
```


## *store remove*

**Description**

Remove activities from local storage


**Syntax**

```sh
$ gravl store remove [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```input```|```i```|Input data store|
|```filter```|```f```|Expression for filtering activities|
|```dryrun```|```n```|Don't actually remove anything, just show what would be done|


## *store update*

**Description**

Query and update Strava activities to local storage


**Syntax**

```sh
$ gravl store update [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```input```|```i```|Input data store|
|```output```|```o```|Output data store|
|```strava.client-id```||API key for Strava API|
|```strava.client-secret```||API secret for Strava API|
|```strava.access-token```||Access token for Strava API|
|```strava.refresh-token```||Refresh token for Strava API|
|```strava.username```||Username for the Strava website|
|```strava.password```||Password for the Strava website|

**Example**

In order to have a more performant experience when running analyzers all Strava activities
are stored locally in a `Store` implementation. The default `Store` is an implementation
using `buntdb` as it allows very simple, fast, and durable local storage though other
implementations exist.

Updates from Strava are incremental and should be run periodically to get the latest activities.

*Note: if the activity already exists locally `gravl` will not update it, it will need to be removed and updated*

```sh
$ gravl store update
2021-02-20T15:59:25-08:00 INF bunt db path="/Users/bzimmer/Library/Application Support/com.github.bzimmer.gravl/gravl.db"
2021-02-20T15:59:26-08:00 INF do all=0 count=100 n=0 start=0 total=0
2021-02-20T15:59:29-08:00 INF do all=100 count=100 n=100 start=1 total=0
2021-02-20T15:59:29-08:00 INF querying activity details ID=4819927284
2021-02-20T15:59:30-08:00 INF saving activity details ID=4819927284 n=1 name="Morning Ride"
2021-02-20T15:59:30-08:00 INF querying activity details ID=4814540574
2021-02-20T15:59:30-08:00 INF saving activity details ID=4814540574 n=2 name="Afternoon Ride"
2021-02-20T15:59:31-08:00 INF do all=200 count=100 n=100 start=2 total=0
2021-02-20T15:59:34-08:00 INF do all=300 count=100 n=100 start=3 total=0
2021-02-20T15:59:36-08:00 INF do all=400 count=100 n=100 start=4 total=0
2021-02-20T15:59:39-08:00 INF do all=500 count=100 n=100 start=5 total=0
2021-02-20T15:59:41-08:00 INF do all=600 count=100 n=100 start=6 total=0
2021-02-20T15:59:43-08:00 INF do all=700 count=100 n=100 start=7 total=0
2021-02-20T15:59:48-08:00 INF do all=800 count=100 n=100 start=8 total=0
2021-02-20T15:59:51-08:00 INF do all=900 count=100 n=100 start=9 total=0
2021-02-20T15:59:54-08:00 INF do all=1000 count=100 n=100 start=10 total=0
2021-02-20T15:59:56-08:00 INF do all=1100 count=100 n=100 start=11 total=0
2021-02-20T15:59:58-08:00 INF do all=1200 count=100 n=100 start=12 total=0
2021-02-20T15:59:59-08:00 INF do all=1200 count=100 n=0 start=13 total=0
{
 "existing": 1198,
 "new": 2,
 "total": 1200
}
```

The results of the command show the number of new, existing, and total activities stored locally.


## *strava*

**Description**

Query Strava for rides and routes


**Flags**

|Name|Aliases|Description|
|-|-|-|
|```strava.client-id```||API key for Strava API|
|```strava.client-secret```||API secret for Strava API|
|```strava.access-token```||Access token for Strava API|
|```strava.refresh-token```||Refresh token for Strava API|
|```strava.username```||Username for the Strava website|
|```strava.password```||Password for the Strava website|

**Overview**

The Strava client is comprised of general [API](https://developers.strava.com/) access supporting
activites, routes, and streams as well as some functionality available by scraping the website as
inspired by [stravaweblib](https://github.com/pR0Ps/stravaweblib).

Additionally, there's full support for implementing `webhooks` but only only webhook management is
available via the commandline (eg [`strava webhook list`](#strava-webhook-list),
[`strava webhook subscribe`](#strava-webhook-subscribe), and [`strava webhook unsubscribe`](#strava-webhook-unsubscribe)).

The entire [`analysis`](#analysis) package is built around Strava activities.


## *strava activities*

**Description**

Query activities for an athlete from Strava


**Syntax**

```sh
$ gravl strava activities [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```count```|```N```|The number of activities to query from Strava (the number returned will be <= N)|
|```filter```|```f```|Expression for filtering activities to remove|
|```attribute```|```B```|Evaluate the expression on an activity and return only those results|

**Example**

List all `VirtualRides` in the last 20 activities and display their `ID`, `Name`, `StartDate`, and their `Distance` in miles

```sh
$ gravl -c --timeout 1m strava activities -N 20 -f ".Type == 'VirtualRide'" -B ".ID, .Name, .StartDateLocal, .Distance.Miles()"
2021-02-20T08:50:32-08:00 INF do all=0 count=20 n=0 start=0 total=20
2021-02-20T08:50:34-08:00 INF do all=20 count=20 n=20 start=1 total=20
[4807285657,"Yorkshire - Jon's Short Mix","2021-02-18T06:56:20Z",10.702124592380498]
[4802094087,"London","2021-02-17T06:55:39Z",12.95105334844508]
[4741552384,"2004","2021-02-05T18:15:27Z",17.51514902966675]
```


## *strava activity*

**Description**

Query an activity from Strava


**Syntax**

```sh
$ gravl strava activity [flags] ACTIVITY_ID (...)
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```stream```|```s```|Streams to include in the activity|

**Example**

To query a specific activity:

```sh
$ gravl strava activity 4802094087
{
 "id": 4802094087,
 "resource_state": 3,
 "external_id": "zwift-activity-752165469668926976.fit",
 "upload_id": 5123870639,
...
 "gear": {
  "id": "b6713218",
  "primary": false,
  "name": "Smart Trainer",
  "resource_state": 2,
  "distance": 711026,
  "athlete_id": 0
 },
 ...
 "device_name": "Zwift",
 "segment_leaderboard_opt_out": false,
 "leaderboard_opt_out": false,
 "perceived_exertion": 6,
 "prefer_perceived_exertion": false
}
```

To include stream data use the `-s` flag:

```sh
$ gravl strava activity -s watts 4802094087
{
 "id": 4802094087,
 "resource_state": 3,
 "external_id": "zwift-activity-752165469668926976.fit",
 "upload_id": 5123870639,
 ...
 "name": "London",
 "distance": 20842.7,
 "moving_time": 2463,
 "elapsed_time": 2463,
 "total_elevation_gain": 216,
 "type": "VirtualRide",
 "start_date": "2021-02-17T14:55:39Z",
 "start_date_local": "2021-02-17T06:55:39Z",
 ...
 "streams": {
  "activity_id": 4802094087,
  "distance": {
   "original_size": 2464,
   "resolution": "high",
   "series_type": "distance",
   "data": [
    2.5,
    4.3,
    6.6,
    9.2,
    ...
    20829.9,
    20835.3,
    20840.2,
    20845.2
   ]
  },
  "watts": {
   "original_size": 2464,
   "resolution": "high",
   "series_type": "distance",
   "data": [
    89,
    105,
    105,
    105,
    ...
    367,
    368,
    376,
    408,
    406,
    396,
    396,
    408,
    412,
    361,
    400,
    ...
   ]
  }
 }
}
```


## *strava athlete*

**Description**

Query an athlete from Strava


**Syntax**

```sh
$ gravl strava athlete [flags]
```



## *strava export*

**Description**

Export a Strava activity by id


**Syntax**

```sh
$ gravl strava export [flags] ACTIVITY_ID (...)
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```overwrite```|```o```|Overwrite the file if it exists; fail otherwise|
|```output```|```O```|The filename to use for writing the contents of the export, if not specified the contents are streamed to stdout|

**Example**

Strava export uses the website and therefore requires a username and password instead of the usual oauth credentials.

If neither `-o` or `-O` are specified the contents of the file are written to stdout.
If `-o` is specified, the file will be written to disk using the name provided by Strava, even if it already exists locally.
If `-O` is specified, the file will be written to disk using the name provided by the flag. It will not overwrite an existing
file unless `-o` was also specified.

```sh
$ gravl strava export -o 4814450574
2021-02-20T09:20:29-08:00 INF export activityID=4814540574 format=original
{
 "id": 4814450574,
 "name": "Friday.fit",
 "format": "fit"
}
$ ls -las Friday.fit
56 -rw-r--r--  1 bzimmer  staff    25K Feb 20 09:20 Friday.fit
```

An example of the overwrite logic.

```sh
$ gravl strava export -O Friday.fit 4814540547
2021-02-20T09:24:44-08:00 INF export activityID=4814540547 format=original
2021-02-20T09:24:45-08:00 ERR file exists and -o flag not specified filename=Friday.fit
2021-02-20T09:24:45-08:00 ERR gravl strava error="file already exists"
```


## *strava fitness*

**Description**

Query Strava for training load data for the authenticated user


**Syntax**

```sh
$ gravl strava fitness [flags]
```


**Example**

Queries the fitness and freshness data for the authenticated user.

```sh
$ gravl strava fitness
[
 {
  "date": {
   "year": 2020,
   "month": 8,
   "day": 22
  },
  "fitness_profile": {
   "fitness": 107.39025265712681,
   "impulse": 83,
   "relative_effort": 64,
   "fatigue": 86.97505302175927,
   "form": 20.415199635367543
  },
  "activities": [
   {
    "id": 3951687537,
    "impulse": 83,
    "relative_effort": 64
   }
  ]
 },
 ...,
 {
  "date": {
   "year": 2021,
   "month": 3,
   "day": 8
  },
  "fitness_profile": {
   "fitness": 56.367803371875894,
   "impulse": 0,
   "relative_effort": 0,
   "fatigue": 5.510876434392897,
   "form": 50.856926937482996
  },
  "activities": []
 }
]


## *strava oauth*

**Description**

Authentication endpoints for access and refresh tokens


**Syntax**

```sh
$ gravl strava oauth [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```origin```||Callback origin|
|```port```||Port on which to listen|


## *strava refresh*

**Description**

Acquire a new refresh token


**Syntax**

```sh
$ gravl strava refresh [flags]
```



## *strava route*

**Description**

Query a route from Strava


**Syntax**

```sh
$ gravl strava route [flags] ROUTE_ID (...)
```



## *strava routes*

**Description**

Query routes for an athlete from Strava


**Syntax**

```sh
$ gravl strava routes [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```count```|```N```|The number of routes to query from Strava (the number returned will be <= N)|


## *strava streams*

**Description**

Query streams for an activity from Strava


**Syntax**

```sh
$ gravl strava streams [flags] ACTIVITY_ID (...)
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```stream```|```s```|Streams to include in the activity|


## *strava streamsets*

**Description**

Return the set of available streams for query


**Syntax**

```sh
$ gravl strava streamsets [flags]
```


**Example**

Query the available streams for query with `activity`.

```sh
$ gravl strava streamsets
{
 "altitude": "The sequence of altitude values for this stream, in meters [float]",
 "cadence": "The sequence of cadence values for this stream, in rotations per minute [integer]",
 "distance": "The sequence of distance values for this stream, in meters [float]",
 "grade_smooth": "The sequence of grade values for this stream, as percents of a grade [float]",
 "heartrate": "The sequence of heart rate values for this stream, in beats per minute [integer]",
 "latlng": "The sequence of lat/long values for this stream [float, float]",
 "moving": "The sequence of moving values for this stream, as boolean values [boolean]",
 "temp": "The sequence of temperature values for this stream, in celsius degrees [float]",
 "time": "The sequence of time values for this stream, in seconds [integer]",
 "velocity_smooth": "The sequence of velocity values for this stream, in meters per second [float]",
 "watts": "The sequence of power values for this stream, in watts [integer]"
}
```

To get just the stream names:

```sh
$ gravl strava streamsets | jq -r "keys | .[]"
altitude
cadence
distance
grade_smooth
heartrate
latlng
moving
temp
time
velocity_smooth
watts
```


## *strava upload*

**Description**

Upload an activity file


**Syntax**

```sh
$ gravl strava upload [flags] {{FILE | DIRECTORY} | UPLOAD_ID (...)}
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```status```|```s```|Check the status of the upload|
|```poll```|```p```|Continually check the status of the request until it is completed|
|```dryrun```|```n```|Show the files which would be uploaded but do not upload them|


## *strava webhook*

**Description**

Manage webhook subscriptions




## *strava webhook list*

**Description**

List all active webhook subscriptions


**Syntax**

```sh
$ gravl strava webhook list [flags]
```



## *strava webhook subscribe*

**Description**

Subscribe for webhook notications


**Syntax**

```sh
$ gravl strava webhook subscribe [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```url```||Address where webhook events will be sent (max length 255 characters|
|```verify```||String chosen by the application owner for client security|


## *strava webhook unsubscribe*

**Description**

Unsubscribe an active webhook subscription (or all if specified)


**Syntax**

```sh
$ gravl strava webhook unsubscribe [flags]
```



## *version*

**Description**

Version


**Syntax**

```sh
$ gravl version [flags]
```



## *visualcrossing*

**Description**

Query VisualCrossing for forecasts


**Flags**

|Name|Aliases|Description|
|-|-|-|
|```visualcrossing.access-token```||API key for Visual Crossing API|


## *visualcrossing forecast*

**Description**

Query VisualCrossing for a forecast


**Syntax**

```sh
$ gravl visualcrossing forecast [flags] [--] LATITUDE LONGITUDE
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```interval```|```i```|Forecast interval (eg 1, 12, 24)|


## *wta*

**Description**

Query the WTA site for trip reports, if no reporter is specified the most recent reports are returned


**Syntax**

```sh
$ gravl wta [flags] [REPORTER_NAME ...]
```


**Example**

Please support the [Washington Trails Association](https://wta.org), thanks!


## *zwift*

**Description**

Query Zwift for activities


**Flags**

|Name|Aliases|Description|
|-|-|-|
|```zwift.username```||Username for the Zwift website|
|```zwift.password```||Password for the Zwift website|

**Overview**

The Zwift client was heavily inspired by [zwift-client](https://github.com/jsmits/zwift-client).

The command [`files`](#zwift-files) is useful for uploading local files to services without direct
integration such as [Cycling Analytics](#cyclinganalytics).


## *zwift activities*

**Description**

Query activities for an athlete from Zwift


**Syntax**

```sh
$ gravl zwift activities [flags]
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```count```|```N```|The number of activities to query from Zwift (the number returned will be <= N)|


## *zwift activity*

**Description**

Query an activity from Zwift


**Syntax**

```sh
$ gravl zwift activity [flags] ACTIVITY_ID (...)
```



## *zwift athlete*

**Description**

Query the athlete profile from Zwift


**Syntax**

```sh
$ gravl zwift athlete [flags]
```



## *zwift export*

**Description**

Export a Zwift activity by id


**Syntax**

```sh
$ gravl zwift export [flags] ACTIVITY_ID (...)
```

**Flags**

|Name|Aliases|Description|
|-|-|-|
|```overwrite```|```o```|Overwrite the file if it exists; fail otherwise|
|```output```|```O```|The filename to use for writing the contents of the export, if not specified the contents are streamed to stdout|


## *zwift files*

**Description**

List all local Zwift files


**Syntax**

```sh
$ gravl zwift files [flags]
```


**Example**

List all local files from the Zwift app's directory. Any files less than 1K in size or named `inProgressActivity.fit` will be ignored.

```sh
$ gravl zwift files
2021-02-19T19:39:25-08:00 WRN skipping, too small file=/Users/bzimmer/Documents/Zwift/Activities/2021-01-12-18-39-52.fit size=584
2021-02-19T19:39:25-08:00 WRN skipping, too small file=/Users/bzimmer/Documents/Zwift/Activities/2021-01-26-18-14-13.fit size=584
2021-02-19T19:39:25-08:00 WRN skipping, activity in progress file=/Users/bzimmer/Documents/Zwift/Activities/inProgressActivity.fit
[
 "/Users/bzimmer/Documents/Zwift/Activities/2021-01-12-18-40-53.fit",
 "/Users/bzimmer/Documents/Zwift/Activities/2021-01-26-18-15-16.fit"
]
```

This command can be easily combined with `jq` to upload files to CyclingAnalytics, Strava, or any other site.

```sh
$ gravl zwift files | jq -r ".[]" | xargs gravl strava upload -n
2021-02-19T19:41:50-08:00 WRN skipping, too small file=/Users/bzimmer/Documents/Zwift/Activities/2021-01-12-18-39-52.fit size=584
2021-02-19T19:41:50-08:00 WRN skipping, too small file=/Users/bzimmer/Documents/Zwift/Activities/2021-01-26-18-14-13.fit size=584
2021-02-19T19:41:50-08:00 WRN skipping, activity in progress file=/Users/bzimmer/Documents/Zwift/Activities/inProgressActivity.fit
2021-02-19T19:41:50-08:00 INF collecting file=/Users/bzimmer/Documents/Zwift/Activities/2021-01-12-18-40-53.fit
2021-02-19T19:41:50-08:00 INF uploading dryrun=true file=2021-01-12-18-40-53.fit
2021-02-19T19:41:50-08:00 INF collecting file=/Users/bzimmer/Documents/Zwift/Activities/2021-01-26-18-15-16.fit
2021-02-19T19:41:50-08:00 INF uploading dryrun=true file=2021-01-26-18-15-16.fit
```


## *zwift refresh*

**Description**

Acquire a new refresh token


**Syntax**

```sh
$ gravl zwift refresh [flags]
```


**Example**

Query for a new refresh token from Zwift.

```sh
$ gravl zwift refresh
{
	"access_token": "12345",
	"token_type": "bearer",
	"refresh_token": "67890",
	"expiry": "2021-02-20T01:29:05.964572-08:00"
}
```
