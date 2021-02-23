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
