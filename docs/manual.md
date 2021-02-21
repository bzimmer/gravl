
# gravl - Activity related analysis, exploration, & planning

## *analysis*

**Description**

Produce statistics and other interesting artifacts from Strava activities

**Syntax:**

```sh
$ gravl analysis
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```units```|```u```|Units|
|```filter```|```f```|Expression for filtering activities to remove|
|```group```|```g```|Expressions for grouping activities|
|```analyzer```|```a```|Analyzers to include (if none specified, default set is used)|
|```input```|```i```|Input data store|

**Example:**

Run the analysis on the Strava activities

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

**Syntax:**

```sh
$ gravl analysis list
```
**Example:**

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

## *analysis manual*

**Description**

Print the manual for the available analyzers

**Syntax:**

```sh
$ gravl analysis manual
```


## *cyclinganalytics activities*

**Description**

Query activities for the authenticated athlete

**Syntax:**

```sh
$ gravl cyclinganalytics activities
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```count```|```N```|Count|


## *cyclinganalytics athlete*

**Description**

Query for the authenticated athlete

**Syntax:**

```sh
$ gravl cyclinganalytics athlete
```

## *cyclinganalytics oauth*

**Description**

Authentication endpoints for access and refresh tokens

**Syntax:**

```sh
$ gravl cyclinganalytics oauth
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```origin```||Callback origin|
|```port```||Port on which to listen|


## *cyclinganalytics activity*

**Description**

Query an activity for the authenticated athlete

**Syntax:**

```sh
$ gravl cyclinganalytics activity
```

## *cyclinganalytics upload*

**Description**

Upload an activity file

**Syntax:**

```sh
$ gravl cyclinganalytics upload
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```status```|```s```|Check the status of the upload|
|```poll```|```p```|Continually check the status of the request until it is completed|
|```dryrun```|```n```|Show the files which would be uploaded but do not upload them|


## *gnis*

**Description**

Query the GNIS database

**Syntax:**

```sh
$ gravl gnis <US STATE ABBREVIATION>
```
**Example:**

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


## *gpx info*

**Description**

Return basic statistics about a GPX file

**Syntax:**

```sh
$ gravl gpx info
```

## *commands*

**Description**

Return all possible commands

**Syntax:**

```sh
$ gravl commands
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```relative```|```r```|Specify the command relative to the current working directory|


## *manual*

**Description**

Generate the `gravl` manual

**Syntax:**

```sh
$ gravl manual
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```help```|```h```|show help|



## *noaa forecast*

**Description**

Query NOAA for a forecast

**Syntax:**

```sh
$ gravl noaa forecast [--] <LATITUDE> <LONGITUDE>
```


## *openweather forecast*

**Description**

Query OpenWeather for a forecast

**Syntax:**

```sh
$ gravl openweather forecast [--] <LATITUDE> <LONGITUDE>
```
**Example:**

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


## *rwgps activities*

**Description**

Query activities for the authenticated athlete

**Syntax:**

```sh
$ gravl rwgps activities
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```count```|```N```|Count|

**Example:**

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

**Syntax:**

```sh
$ gravl rwgps activity
```

## *rwgps athlete*

**Description**

Query for the authenticated athlete

**Syntax:**

```sh
$ gravl rwgps athlete
```

## *rwgps route*

**Description**

Query a route from RideWithGPS

**Syntax:**

```sh
$ gravl rwgps route
```

## *rwgps routes*

**Description**

Query routes for an athlete from RideWithGPS

**Syntax:**

```sh
$ gravl rwgps routes
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```count```|```N```|Count|


## *srtm*

**Description**

Query the SRTM database for elevation data

**Syntax:**

```sh
$ gravl srtm
```


## *store export*

**Description**

Export activities from local storage

**Syntax:**

```sh
$ gravl store export
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```input```|```i```|Input data store|
|```filter```|```f```|Expression for filtering activities|
|```attribute```|```B```|Evaluate the expression on an activity and return only those results|


## *store remove*

**Description**

Remove activities from local storage

**Syntax:**

```sh
$ gravl store remove
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```input```|```i```|Input data store|
|```filter```|```f```|Expression for filtering activities|
|```dryrun```|```n```|Don't actually remove anything, just show what would be done|


## *store update*

**Description**

Query and update Strava activities to local storage

**Syntax:**

```sh
$ gravl store update
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```input```|```i```|Input data store|
|```output```|```o```|Output data store|
|```strava.client-id```||API key for Strava API|
|```strava.client-secret```||API secret for Strava API|
|```strava.access-token```||Access token for Strava API|
|```strava.refresh-token```||Refresh token for Strava API|
|```strava.username```||Username for the Strava website|
|```strava.password```||Password for the Strava website|

**Example:**

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


## *strava activities*

**Description**

Query activities for an athlete from Strava

**Syntax:**

```sh
$ gravl strava activities
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```count```|```N```|Count|
|```filter```|```f```|Expression for filtering activities to remove|
|```attribute```|```B```|Evaluate the expression on an activity and return only those results|

**Example:**

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

**Syntax:**

```sh
$ gravl strava activity
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```stream```|```s```|Streams to include in the activity|


## *strava athlete*

**Description**

Query an athlete from Strava

**Syntax:**

```sh
$ gravl strava athlete
```

## *strava export*

**Description**

Export a Strava activity by id

**Syntax:**

```sh
$ gravl strava export
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```format```|```F```|Export data file in the specified format|
|```overwrite```|```o```|Overwrite the file if it exists; fail otherwise|
|```output```|```O```|The filename to use for writing the contents of the export, if not specified the contents are streamed to stdout|

**Example:**

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

## *strava fitness*

**Description**

Query Strava for training load data

**Syntax:**

```sh
$ gravl strava fitness
```

## *strava oauth*

**Description**

Authentication endpoints for access and refresh tokens

**Syntax:**

```sh
$ gravl strava oauth
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```origin```||Callback origin|
|```port```||Port on which to listen|


## *strava refresh*

**Description**

Acquire a new refresh token

**Syntax:**

```sh
$ gravl strava refresh
```

## *strava route*

**Description**

Query a route from Strava

**Syntax:**

```sh
$ gravl strava route
```

## *strava routes*

**Description**

Query routes for an athlete from Strava

**Syntax:**

```sh
$ gravl strava routes
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```count```|```N```|Count|


## *strava stream*

**Description**

Query streams for an activity from Strava

**Syntax:**

```sh
$ gravl strava stream
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```stream```|```s```|Streams to include in the activity|


## *strava upload*

**Description**

Upload an activity file

**Syntax:**

```sh
$ gravl strava upload <FILE or DIRECTORY> | <UPLOAD ID>
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```status```|```s```|Check the status of the upload|
|```poll```|```p```|Continually check the status of the request until it is completed|
|```dryrun```|```n```|Show the files which would be uploaded but do not upload them|



## *strava webhook list*

**Description**

List all active webhook subscriptions

**Syntax:**

```sh
$ gravl strava webhook list
```

## *strava webhook subscribe*

**Description**

Subscribe for webhook notications

**Syntax:**

```sh
$ gravl strava webhook subscribe
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```url```||Address where webhook events will be sent (max length 255 characters|
|```verify```||String chosen by the application owner for client security|


## *strava webhook unsubscribe*

**Description**

Unsubscribe an active webhook subscription (or all if specified)

**Syntax:**

```sh
$ gravl strava webhook unsubscribe
```

## *version*

**Description**

Version

**Syntax:**

```sh
$ gravl version
```


## *visualcrossing forecast*

**Description**

Query VisualCrossing for a forecast

**Syntax:**

```sh
$ gravl visualcrossing forecast [--] <LATITUDE> <LONGITUDE>
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```interval```|```i```|Forecast interval (eg 1, 12, 24)|


## *wta*

**Description**

Query the WTA site for trip reports

**Syntax:**

```sh
$ gravl wta
```


## *zwift activities*

**Description**

Query activities for an athlete from Strava

**Syntax:**

```sh
$ gravl zwift activities
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```count```|```N```|Count|


## *zwift activity*

**Description**

Query an activity from Zwift

**Syntax:**

```sh
$ gravl zwift activity
```

## *zwift athlete*

**Description**

Query the athlete profile from Zwift

**Syntax:**

```sh
$ gravl zwift athlete
```

## *zwift export*

**Description**

Export a Zwift activity by id

**Syntax:**

```sh
$ gravl zwift export
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```overwrite```|```o```|Overwrite the file if it exists; fail otherwise|
|```output```|```O```|The filename to use for writing the contents of the export, if not specified the contents are streamed to stdout|


## *zwift files*

**Description**

List all local Zwift files

**Syntax:**

```sh
$ gravl zwift files
```
**Example:**

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

**Syntax:**

```sh
$ gravl zwift refresh
```
**Example:**

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

## *help*

**Description**

Shows a list of commands or help for one command

**Syntax:**

```sh
$ gravl help [command]
```

