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
