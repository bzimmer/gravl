# Commands



## Global Flags
|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|verbosity|v||Log level (trace, debug, info, warn, error, fatal, panic)|
|monochrome|m||Use monochrome logging, color enabled by default|
|json|j||Emit all results as JSON and print to stdout|
|http-tracing|||Log all http calls (warning: no effort is made to mask log ids, keys, and other sensitive information)|
|timeout|t||Timeout duration (eg, 1ms, 2s, 5m, 3h)|
|help|h||show help|

## Commands
* [cyclinganalytics](#cyclinganalytics)
* [cyclinganalytics activities](#cyclinganalytics-activities)
* [cyclinganalytics activity](#cyclinganalytics-activity)
* [cyclinganalytics athlete](#cyclinganalytics-athlete)
* [cyclinganalytics oauth](#cyclinganalytics-oauth)
* [cyclinganalytics streamsets](#cyclinganalytics-streamsets)
* [envvars](#envvars)
* [help](#help)
* [qp](#qp)
* [qp copy](#qp-copy)
* [qp export](#qp-export)
* [qp list](#qp-list)
* [qp providers](#qp-providers)
* [qp status](#qp-status)
* [qp upload](#qp-upload)
* [rwgps](#rwgps)
* [rwgps activities](#rwgps-activities)
* [rwgps activity](#rwgps-activity)
* [rwgps athlete](#rwgps-athlete)
* [rwgps route](#rwgps-route)
* [rwgps routes](#rwgps-routes)
* [strava](#strava)
* [strava activities](#strava-activities)
* [strava activity](#strava-activity)
* [strava athlete](#strava-athlete)
* [strava oauth](#strava-oauth)
* [strava photos](#strava-photos)
* [strava refresh](#strava-refresh)
* [strava route](#strava-route)
* [strava routes](#strava-routes)
* [strava streams](#strava-streams)
* [strava streamsets](#strava-streamsets)
* [strava update](#strava-update)
* [strava webhook](#strava-webhook)
* [strava webhook list](#strava-webhook-list)
* [strava webhook subscribe](#strava-webhook-subscribe)
* [strava webhook unsubscribe](#strava-webhook-unsubscribe)
* [version](#version)
* [zwift](#zwift)
* [zwift activities](#zwift-activities)
* [zwift activity](#zwift-activity)
* [zwift athlete](#zwift-athlete)
* [zwift files](#zwift-files)
* [zwift refresh](#zwift-refresh)

### *cyclinganalytics*

**Description**

Operations supported by the CyclingAnalytics API



**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|cyclinganalytics-client-id||CYCLINGANALYTICS_CLIENT_ID|CyclingAnalytics client id|
|cyclinganalytics-client-secret||CYCLINGANALYTICS_CLIENT_SECRET|CyclingAnalytics client secret|
|cyclinganalytics-access-token||CYCLINGANALYTICS_ACCESS_TOKEN|CyclingAnalytics access token|
|rate-limit|||Minimum time interval between API request events (eg, 1ms, 2s, 5m, 3h)|
|rate-burst|||Maximum burst size for API request events|
|concurrency|||Maximum concurrent API queries|


### *cyclinganalytics activities*

**Description**

Query activities for the authenticated athlete



**Syntax**

```sh
$ gravl cyclinganalytics activities [flags]
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|count|N||The number of activities to query from CA (the number returned will be <= N)|


### *cyclinganalytics activity*

**Description**

Query an activity for the authenticated athlete



**Syntax**

```sh
$ gravl cyclinganalytics activity [flags]
```



### *cyclinganalytics athlete*

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


### *cyclinganalytics oauth*

**Description**

Authentication endpoints for access and refresh tokens



**Syntax**

```sh
$ gravl cyclinganalytics oauth [flags]
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|origin|||Callback origin|
|port|||Port on which to listen|


### *cyclinganalytics streamsets*

**Description**

Return the set of available streams for query



**Syntax**

```sh
$ gravl cyclinganalytics streamsets [flags]
```



### *envvars*

**Description**

Useful for creating a .env file for all possible environment variables



**Syntax**

```sh
$ gravl envvars [flags]
```



### *help*

**Description**

Shows a list of commands or help for one command



**Syntax**

```sh
$ gravl help [flags] [command]
```



### *qp*

**Description**

Manage the flow of activity between different platforms



**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|cyclinganalytics-client-id||CYCLINGANALYTICS_CLIENT_ID|CyclingAnalytics client id|
|cyclinganalytics-client-secret||CYCLINGANALYTICS_CLIENT_SECRET|CyclingAnalytics client secret|
|cyclinganalytics-access-token||CYCLINGANALYTICS_ACCESS_TOKEN|CyclingAnalytics access token|
|rwgps-client-id||RWGPS_CLIENT_ID|RideWithGPS client id|
|rwgps-access-token||RWGPS_ACCESS_TOKEN|RideWithGPS access token|
|strava-client-id||STRAVA_CLIENT_ID|Strava client id|
|strava-client-secret||STRAVA_CLIENT_SECRET|Strava client secret|
|strava-refresh-token||STRAVA_REFRESH_TOKEN|Strava refresh token|
|zwift-username||ZWIFT_USERNAME|Zwift username|
|zwift-password||ZWIFT_PASSWORD|Zwift password|


### *qp copy*

**Description**

Copy an activity from a source to a destination



**Syntax**

```sh
$ gravl qp copy [flags] --from <exporter> --to <uploader> id [id, ...]
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|rate-limit|||Minimum time interval between API request events (eg, 1ms, 2s, 5m, 3h)|
|rate-burst|||Maximum burst size for API request events|
|concurrency|||Maximum concurrent API queries|
|from|||Source data provider|
|to|||Sink data provider|
|overwrite|o||Overwrite the file if it exists; fail otherwise|
|output|O||The filename to use for writing the contents of the export,
if not specified the contents are streamed to stdout|
|poll|||Continually check the status of the request until it is completed|
|interval|||The amount of time to wait between polling for an updated status|
|iterations|N||The max number of polling iterations to perform|

**Example**

```sh
$ gravl qp copy --from zwift --to cyclinganalytics 934398333398662432
2021-10-28T06:25:55-07:00 INF created zwift client
2021-10-28T06:25:55-07:00 INF created cyclinganalytics client
2021-10-28T06:25:56-07:00 INF export exp=2021-10-26-18-19-39.fit id=934398333398662432
2021-10-28T06:25:57-07:00 INF poll id=1067278367 iteration=0
2021-10-28T06:26:00-07:00 INF poll id=1067278367 iteration=1
2021-10-28T06:26:00-07:00 INF counters count=1 metric=gravl.upload.file.success
2021-10-28T06:26:00-07:00 INF counters count=2 metric=gravl.upload.poll
2021-10-28T06:26:00-07:00 INF samples count=1 max=4.181048393249512 mean=4.181048393249512 metric=gravl.runtime min=4.181048393249512 stddev=0
```

```sh
$ gravl qp copy --from zwift --to cyclinganalytics 934398333398662432
2021-10-28T06:34:06-07:00 INF created zwift client
2021-10-28T06:34:06-07:00 INF created cyclinganalytics client
2021-10-28T06:34:07-07:00 INF export exp=2021-10-26-18-19-39.fit id=934398333398662432
2021-10-28T06:34:08-07:00 INF poll id=9819686356 iteration=0
{
 "upload_id": 9819686356,
 "status": "processing",
 "ride_id": 0,
 "user_id": 1603544,
 "format": "fit",
 "datetime": "2021-10-28T13:34:09",
 "filename": "2021-10-26-18-19-39.fit",
 "size": 113324,
 "error": "",
 "error_code": ""
}
2021-10-28T06:34:10-07:00 INF poll id=9819686356 iteration=1
{
 "upload_id": 9819686356,
 "status": "error",
 "ride_id": 0,
 "user_id": 1603544,
 "format": "fit",
 "datetime": "2021-10-28T13:34:09",
 "filename": "2021-10-26-18-19-39.fit",
 "size": 113324,
 "error": "The ride already exists: 894091983723",
 "error_code": "duplicate_ride"
}
2021-10-28T06:34:10-07:00 INF counters count=2 metric=gravl.upload.poll
2021-10-28T06:34:10-07:00 INF counters count=1 metric=gravl.upload.file.success
2021-10-28T06:34:10-07:00 INF samples count=1 max=3.807884454727173 mean=3.807884454727173 metric=gravl.runtime min=3.807884454727173 stddev=0
```


### *qp export*

**Description**

Export an activity from the source



**Syntax**

```sh
$ gravl qp export [flags]
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|rate-limit|||Minimum time interval between API request events (eg, 1ms, 2s, 5m, 3h)|
|rate-burst|||Maximum burst size for API request events|
|concurrency|||Maximum concurrent API queries|
|from|||Source data provider|
|overwrite|o||Overwrite the file if it exists; fail otherwise|
|output|O||The filename to use for writing the contents of the export,
if not specified the contents are streamed to stdout|

**Example**

If neither `-o` or `-O` are specified the contents of the file are written to stdout.
If `-o` is specified, the file will be written to disk using the name provided by Strava, even if it already exists locally.
If `-O` is specified, the file will be written to disk using the name provided by the flag. It will not overwrite an existing
file unless `-o` was also specified.


### *qp list*

**Description**

List the files suitable for uploading



**Syntax**

```sh
$ gravl qp list [flags] {FILE | DIRECTORY} (...)
```



### *qp providers*

**Description**

Return the set of active exporters and uploaders



**Syntax**

```sh
$ gravl qp providers [flags]
```



### *qp status*

**Description**

Check the status of an upload



**Syntax**

```sh
$ gravl qp status [flags] UPLOAD_ID (...)
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|rate-limit|||Minimum time interval between API request events (eg, 1ms, 2s, 5m, 3h)|
|rate-burst|||Maximum burst size for API request events|
|concurrency|||Maximum concurrent API queries|
|to|||Sink data provider|
|poll|||Continually check the status of the request until it is completed|
|interval|||The amount of time to wait between polling for an updated status|
|iterations|N||The max number of polling iterations to perform|


### *qp upload*

**Description**





**Syntax**

```sh
$ gravl qp upload [flags] {FILE | DIRECTORY} (...)
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|rate-limit|||Minimum time interval between API request events (eg, 1ms, 2s, 5m, 3h)|
|rate-burst|||Maximum burst size for API request events|
|concurrency|||Maximum concurrent API queries|
|to|||Sink data provider|
|poll|||Continually check the status of the request until it is completed|
|interval|||The amount of time to wait between polling for an updated status|
|iterations|N||The max number of polling iterations to perform|


### *rwgps*

**Description**

Operations supported by the RideWithGPS API



**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|rwgps-client-id||RWGPS_CLIENT_ID|RideWithGPS client id|
|rwgps-access-token||RWGPS_ACCESS_TOKEN|RideWithGPS access token|
|rate-limit|||Minimum time interval between API request events (eg, 1ms, 2s, 5m, 3h)|
|rate-burst|||Maximum burst size for API request events|
|concurrency|||Maximum concurrent API queries|


### *rwgps activities*

**Description**

Query activities for the authenticated athlete



**Syntax**

```sh
$ gravl rwgps activities [flags]
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|count|N||The number of activities to query from RideWithGPS|

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


### *rwgps activity*

**Description**

Query an activity from RideWithGPS



**Syntax**

```sh
$ gravl rwgps activity [flags]
```



### *rwgps athlete*

**Description**

Query for the authenticated athlete



**Syntax**

```sh
$ gravl rwgps athlete [flags]
```



### *rwgps route*

**Description**

Query a route from RideWithGPS



**Syntax**

```sh
$ gravl rwgps route [flags]
```



### *rwgps routes*

**Description**

Query routes for an athlete from RideWithGPS



**Syntax**

```sh
$ gravl rwgps routes [flags]
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|count|N||The number of routes to query from RideWithGPS|


### *strava*

**Description**

Operations supported by the Strava API



**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|strava-client-id||STRAVA_CLIENT_ID|Strava client id|
|strava-client-secret||STRAVA_CLIENT_SECRET|Strava client secret|
|strava-refresh-token||STRAVA_REFRESH_TOKEN|Strava refresh token|
|rate-limit|||Minimum time interval between API request events (eg, 1ms, 2s, 5m, 3h)|
|rate-burst|||Maximum burst size for API request events|
|concurrency|||Maximum concurrent API queries|

**Overview**

The Strava client is comprised of general [API](https://developers.strava.com/) access supporting
activities, routes, and streams.

Additionally, there's full support for implementing `webhooks` but only only webhook management is
available via the commandline (eg [`strava webhook list`](#strava-webhook-list),
[`strava webhook subscribe`](#strava-webhook-subscribe), and [`strava webhook unsubscribe`](#strava-webhook-unsubscribe)).


### *strava activities*

**Description**

Query activities for an athlete from Strava



**Syntax**

```sh
$ gravl strava activities [flags]
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|count|N||The number of activities to query from Strava (the number returned will be <= N)|
|filter|f||Expression for filtering activities to remove|
|attribute|B||Evaluate the expression on an activity and return only those results|
|after|since||Return results after the time specified|
|before|||Return results before the time specified|

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


### *strava activity*

**Description**

Query an activity from Strava



**Syntax**

```sh
$ gravl strava activity [flags] ACTIVITY_ID (...)
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|stream|s||Streams to include in the activity|

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


### *strava athlete*

**Description**

Query an athlete from Strava



**Syntax**

```sh
$ gravl strava athlete [flags]
```



### *strava oauth*

**Description**

Authentication endpoints for access and refresh tokens



**Syntax**

```sh
$ gravl strava oauth [flags]
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|origin|||Callback origin|
|port|||Port on which to listen|


### *strava photos*

**Description**

Query photos from Strava



**Syntax**

```sh
$ gravl strava photos [flags] ACTIVITY_ID (...)
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|size|s|||


### *strava refresh*

**Description**

Acquire a new refresh token



**Syntax**

```sh
$ gravl strava refresh [flags]
```



### *strava route*

**Description**

Query a route from Strava



**Syntax**

```sh
$ gravl strava route [flags] ROUTE_ID (...)
```



### *strava routes*

**Description**

Query routes for an athlete from Strava



**Syntax**

```sh
$ gravl strava routes [flags]
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|count|N||The number of routes to query from Strava (the number returned will be <= N)|


### *strava streams*

**Description**

Query streams for an activity from Strava



**Syntax**

```sh
$ gravl strava streams [flags] ACTIVITY_ID (...)
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|stream|s||Streams to include in the activity|


### *strava streamsets*

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


### *strava update*

**Description**





**Syntax**

```sh
$ gravl strava update [flags] ACTIVITY_ID (...)
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|name|||Set the name for the activity|
|gear|||Set the gear id for the activity|
|sport|||Set the sport for the activity|
|description|||Set the description for the activity|
|hidden|||Hide the activity from the home dashboard|
|no-hidden|||Display the activity on the home dashboard|
|commute|||The activity is a commute|
|no-commute|||The activity is not a commute|
|trainer|||The activity was completed on a trainer|
|no-trainer|||The activity was not completed on a trainer|


### *strava webhook*

**Description**

Manage webhook subscriptions




### *strava webhook list*

**Description**

List all active webhook subscriptions



**Syntax**

```sh
$ gravl strava webhook list [flags]
```



### *strava webhook subscribe*

**Description**

Subscribe for webhook notifications



**Syntax**

```sh
$ gravl strava webhook subscribe [flags]
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|url|||Address where webhook events will be sent (max length 255 characters|
|verify|||String chosen by the application owner for client security|


### *strava webhook unsubscribe*

**Description**

Unsubscribe an active webhook subscription (or all if specified)



**Syntax**

```sh
$ gravl strava webhook unsubscribe [flags]
```



### *version*

**Description**

Version



**Syntax**

```sh
$ gravl version [flags]
```



### *zwift*

**Description**

Operations supported by the Zwift API



**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|zwift-username||ZWIFT_USERNAME|Zwift username|
|zwift-password||ZWIFT_PASSWORD|Zwift password|
|rate-limit|||Minimum time interval between API request events (eg, 1ms, 2s, 5m, 3h)|
|rate-burst|||Maximum burst size for API request events|
|concurrency|||Maximum concurrent API queries|

**Overview**

The Zwift client was heavily inspired by [zwift-client](https://github.com/jsmits/zwift-client).

The command [`files`](#zwift-files) is useful for uploading local files to services without direct
integration such as [Cycling Analytics](#cyclinganalytics).


### *zwift activities*

**Description**

Query activities for an athlete from Zwift



**Syntax**

```sh
$ gravl zwift activities [flags]
```


**Flags**

|Name|Aliases|EnvVars|Description|
|-|-|-|-|
|count|N||The number of activities to query from Zwift (the number returned will be <= N)|


### *zwift activity*

**Description**

Query an activity from Zwift



**Syntax**

```sh
$ gravl zwift activity [flags] ACTIVITY_ID (...)
```



### *zwift athlete*

**Description**

Query the athlete profile from Zwift



**Syntax**

```sh
$ gravl zwift athlete [flags]
```



### *zwift files*

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


### *zwift refresh*

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

