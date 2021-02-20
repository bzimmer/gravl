# gravl - Activity related analysis, exploration, & planning

### *analysis* - Produce statistics and other interesting artifacts from Strava activities

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

### *analysis list* - Return the list of available analyzers

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


### *cyclinganalytics activities* - Query activities for the authenticated athlete

**Syntax:**

```sh
$ gravl cyclinganalytics activities
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```count```|```N```|Count|


### *cyclinganalytics athlete* - Query for the authenticated athlete

**Syntax:**

```sh
$ gravl cyclinganalytics athlete
```

### *cyclinganalytics oauth* - Authentication endpoints for access and refresh tokens

**Syntax:**

```sh
$ gravl cyclinganalytics oauth
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```origin```||Callback origin|
|```port```||Port on which to listen|


### *cyclinganalytics activity* - Query an activity for the authenticated athlete

**Syntax:**

```sh
$ gravl cyclinganalytics activity
```

### *cyclinganalytics upload* - Upload an activity file

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


### *gnis* - Query the GNIS database

**Syntax:**

```sh
$ gravl gnis
```


### *gpx info* - Return basic statistics about a GPX file

**Syntax:**

```sh
$ gravl gpx info
```

### *commands* - Return all possible commands

**Syntax:**

```sh
$ gravl commands
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```relative```|```r```|Specify the command relative to the current working directory|


### *manual* - Generate the 'gravl' manual

**Syntax:**

```sh
$ gravl manual
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```help```|```h```|show help|



### *noaa forecast* - 

**Syntax:**

```sh
$ gravl noaa forecast
```


### *openweather forecast* - Query OpenWeather for a forecast

**Syntax:**

```sh
$ gravl openweather forecast
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


### *rwgps activities* - Query activities for the authenticated athlete

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

### *rwgps activity* - Query an activity from RideWithGPS

**Syntax:**

```sh
$ gravl rwgps activity
```

### *rwgps athlete* - Query for the authenticated athlete

**Syntax:**

```sh
$ gravl rwgps athlete
```

### *rwgps route* - Query a route from RideWithGPS

**Syntax:**

```sh
$ gravl rwgps route
```

### *rwgps routes* - Query routes for an athlete from RideWithGPS

**Syntax:**

```sh
$ gravl rwgps routes
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```count```|```N```|Count|


### *srtm* - Query the SRTM database for elevation data

**Syntax:**

```sh
$ gravl srtm
```


### *store export* - Export activities from local storage

**Syntax:**

```sh
$ gravl store export
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```input```|```i```|Input data store|
|```filter```|```f```|Expression for filtering activities|


### *store remove* - Remove activities from local storage

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


### *store update* - Query and update Strava activities to local storage

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



### *strava activities* - Query activities for an athlete from Strava

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

### *strava activity* - Query an activity from Strava

**Syntax:**

```sh
$ gravl strava activity
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```stream```|```s```|Streams to include in the activity|


### *strava athlete* - Query an athlete from Strava

**Syntax:**

```sh
$ gravl strava athlete
```

### *strava export* - Export a Strava activity by id

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

### *strava fitness* - Query Strava for training load data

**Syntax:**

```sh
$ gravl strava fitness
```

### *strava oauth* - Authentication endpoints for access and refresh tokens

**Syntax:**

```sh
$ gravl strava oauth
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```origin```||Callback origin|
|```port```||Port on which to listen|


### *strava refresh* - Acquire a new refresh token

**Syntax:**

```sh
$ gravl strava refresh
```

### *strava route* - Query a route from Strava

**Syntax:**

```sh
$ gravl strava route
```

### *strava routes* - Query routes for an athlete from Strava

**Syntax:**

```sh
$ gravl strava routes
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```count```|```N```|Count|


### *strava stream* - Query streams for an activity from Strava

**Syntax:**

```sh
$ gravl strava stream
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```stream```|```s```|Streams to include in the activity|


### *strava upload* - Upload an activity file

**Syntax:**

```sh
$ gravl strava upload
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```status```|```s```|Check the status of the upload|
|```poll```|```p```|Continually check the status of the request until it is completed|
|```dryrun```|```n```|Show the files which would be uploaded but do not upload them|



### *strava webhook list* - List all active webhook subscriptions

**Syntax:**

```sh
$ gravl strava webhook list
```

### *strava webhook subscribe* - Subscribe for webhook notications

**Syntax:**

```sh
$ gravl strava webhook subscribe
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```url```||Address where webhook events will be sent (max length 255 characters|
|```verify```||String chosen by the application owner for client security|


### *strava webhook unsubscribe* - Unsubscribe an active webhook subscription (or all if specified)

**Syntax:**

```sh
$ gravl strava webhook unsubscribe
```

### *version* - Version

**Syntax:**

```sh
$ gravl version
```


### *visualcrossing forecast* - 

**Syntax:**

```sh
$ gravl visualcrossing forecast
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```interval```|```i```|Forecast interval (eg 1, 12, 24)|


### *wta* - Query the WTA site for trip reports

**Syntax:**

```sh
$ gravl wta
```


### *zwift activities* - Query activities for an athlete from Strava

**Syntax:**

```sh
$ gravl zwift activities
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```count```|```N```|Count|


### *zwift activity* - Query an activity from Zwift

**Syntax:**

```sh
$ gravl zwift activity
```

### *zwift athlete* - Query the athlete profile from Zwift

**Syntax:**

```sh
$ gravl zwift athlete
```

### *zwift export* - Export a Zwift activity by id

**Syntax:**

```sh
$ gravl zwift export
```
**Flags:**

|Flag|Short|Description|
|-|-|-|
|```overwrite```|```o```|Overwrite the file if it exists; fail otherwise|
|```output```|```O```|The filename to use for writing the contents of the export, if not specified the contents are streamed to stdout|


### *zwift files* - List all local Zwift files

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

### *zwift refresh* - Acquire a new refresh token

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

### *help* - Shows a list of commands or help for one command

**Syntax:**

```sh
$ gravl help
```

