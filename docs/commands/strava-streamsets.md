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
