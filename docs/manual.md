# gravl - Clients for activty-related services and an extensible analysis framework for activities

### Produce statistics and other interesting artifacts from Strava activities

**Syntax:**

```sh
$ gravl analysis
```

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
 "Run": {
  "totals": {
   "count": 1,
   "distance": 6.367501292452079,
   "elevation": 741.469816272966,
   "calories": 306,
   "movingtime": 3940,
   "centuries": {
    "metric": 0,
    "imperial": 0
   }
  }
 },
 "Snowshoe": {
  "totals": {
   "count": 1,
   "distance": 1.4278488626421697,
   "elevation": 79.3963254593176,
   "calories": 176,
   "movingtime": 1748,
   "centuries": {
    "metric": 0,
    "imperial": 0
   }
  }
 },
 "VirtualRide": {
  "totals": {
   "count": 6,
   "distance": 82.22194881889764,
   "elevation": 5813.648293963255,
   "calories": 3570,
   "movingtime": 15503,
   "centuries": {
    "metric": 0,
    "imperial": 0
   }
  }
 },
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

### Return the list of available analyzers

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
	"benford": {
		"base": false,
		"doc": "benford returns the benford distribution of all the activities",
		"flags": false
	},
	"climbing": {
		"base": true,
		"doc": "climbing returns all activities exceeding the Golden Ratio - https://blog.wahoofitness.com/by-the-numbers-distance-and-elevation/",
		"flags": false
	},
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
	"forecast": {
		"base": false,
		"doc": "forecast the weather for an activity",
		"flags": false
	},
	"hourrecord": {
		"base": true,
		"doc": "hourrecord returns the longest distance traveled (in miles | kilometers) exceeding the average speed (mph | mps)",
		"flags": false
	},
	"koms": {
		"base": true,
		"doc": "koms returns all KOMs for the activities",
		"flags": false
	},
	"pythagorean": {
		"base": true,
		"doc": "pythagorean determines the activity with the highest pythagorean value defined as the sqrt(distance^2 + elevation^2) in meters",
		"flags": false
	},
	"rolling": {
		"base": true,
		"doc": "rolling returns the 'window' of activities with the highest accumulated distance.",
		"flags": true
	},
	"splat": {
		"base": false,
		"doc": "splat returns all activities in the units specified. This analyzer is useful for debugging the filter",
		"flags": false
	},
	"staticmap": {
		"base": false,
		"doc": "staticmap generates a staticmap for every activity",
		"flags": true
	},
	"totals": {
		"base": true,
		"doc": "totals returns the number of centuries (100 mi or 100 km)",
		"flags": false
	}
}
```

### Query CyclingAnalytics

**Syntax:**

```sh
$ gravl cyclinganalytics
```

### Query activities for the authenticated athlete

**Syntax:**

```sh
$ gravl cyclinganalytics activities
```

### Query for the authenticated athlete

**Syntax:**

```sh
$ gravl cyclinganalytics athlete
```

### Authentication endpoints for access and refresh tokens

**Syntax:**

```sh
$ gravl cyclinganalytics oauth
```

### Query an activity for the authenticated athlete

**Syntax:**

```sh
$ gravl cyclinganalytics activity
```

### Upload an activity file

**Syntax:**

```sh
$ gravl cyclinganalytics upload
```

### Query the GNIS database

**Syntax:**

```sh
$ gravl gnis
```

### gpx

**Syntax:**

```sh
$ gravl gpx
```

### Return basic statistics about a GPX file

**Syntax:**

```sh
$ gravl gpx info
```

### Return all possible commands

**Syntax:**

```sh
$ gravl commands
```

### Generate the 'gravl' manual

**Syntax:**

```sh
$ gravl manual
```

### Query NOAA for forecasts

**Syntax:**

```sh
$ gravl noaa
```

### 

**Syntax:**

```sh
$ gravl noaa forecast
```

### Query OpenWeather for forecasts

**Syntax:**

```sh
$ gravl openweather
```

### Query OpenWeather for a forecast

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

### Query RideWithGPS for rides and routes

**Syntax:**

```sh
$ gravl rwgps
```

### Query activities for the authenticated athlete

**Syntax:**

```sh
$ gravl rwgps activities
```

### Query an activity from RideWithGPS

**Syntax:**

```sh
$ gravl rwgps activity
```

### Query for the authenticated athlete

**Syntax:**

```sh
$ gravl rwgps athlete
```

### Query a route from RideWithGPS

**Syntax:**

```sh
$ gravl rwgps route
```

### Query routes for an athlete from RideWithGPS

**Syntax:**

```sh
$ gravl rwgps routes
```

### Query the SRTM database for elevation data

**Syntax:**

```sh
$ gravl srtm
```

### Manage a local store of Strava activities

**Syntax:**

```sh
$ gravl store
```

### Export activities from local storage

**Syntax:**

```sh
$ gravl store export
```

### Remove activities from local storage

**Syntax:**

```sh
$ gravl store remove
```

### Query and update Strava activities to local storage

**Syntax:**

```sh
$ gravl store update
```

### Query Strava for rides and routes

**Syntax:**

```sh
$ gravl strava
```

### Query activities for an athlete from Strava

**Syntax:**

```sh
$ gravl strava activities
```

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

### Query an activity from Strava

**Syntax:**

```sh
$ gravl strava activity
```

### Query an athlete from Strava

**Syntax:**

```sh
$ gravl strava athlete
```

### Export a Strava activity by id

**Syntax:**

```sh
$ gravl strava export
```

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

### Query Strava for training load data

**Syntax:**

```sh
$ gravl strava fitness
```

### Authentication endpoints for access and refresh tokens

**Syntax:**

```sh
$ gravl strava oauth
```

### Acquire a new refresh token

**Syntax:**

```sh
$ gravl strava refresh
```

### Query a route from Strava

**Syntax:**

```sh
$ gravl strava route
```

### Query routes for an athlete from Strava

**Syntax:**

```sh
$ gravl strava routes
```

### Query streams for an activity from Strava

**Syntax:**

```sh
$ gravl strava stream
```

### Upload an activity file

**Syntax:**

```sh
$ gravl strava upload
```

### Manage webhook subscriptions

**Syntax:**

```sh
$ gravl strava webhook
```

### List all active webhook subscriptions

**Syntax:**

```sh
$ gravl strava webhook list
```

### Subscribe for webhook notications

**Syntax:**

```sh
$ gravl strava webhook subscribe
```

### Unsubscribe an active webhook subscription (or all if specified)

**Syntax:**

```sh
$ gravl strava webhook unsubscribe
```

### Version

**Syntax:**

```sh
$ gravl version
```

### Query VisualCrossing for forecasts

**Syntax:**

```sh
$ gravl visualcrossing
```

### 

**Syntax:**

```sh
$ gravl visualcrossing forecast
```

### Query the WTA site for trip reports

**Syntax:**

```sh
$ gravl wta
```

### Query Zwift for activities

**Syntax:**

```sh
$ gravl zwift
```

### Query activities for an athlete from Strava

**Syntax:**

```sh
$ gravl zwift activities
```

### Query an activity from Zwift

**Syntax:**

```sh
$ gravl zwift activity
```

### Query the athlete profile from Zwift

**Syntax:**

```sh
$ gravl zwift athlete
```

### Export a Zwift activity by id

**Syntax:**

```sh
$ gravl zwift export
```

### List all local Zwift files

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

### Acquire a new refresh token

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

### Shows a list of commands or help for one command

**Syntax:**

```sh
$ gravl help
```

