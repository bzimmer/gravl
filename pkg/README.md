# pkg

Package pkg provides access to a number of different clients capable of querying
data useful in the reviewing and planning of activities.

The general pattern is the same for all client implementations with small variances
due to the authentication requirements of each (some require no auth, some require
tokens).

They generally fall into three categories:

*Activities*

* [Strava](./strava/strava.go)

* [RideWithGPS](./rwgps/rwgps.go)

* [CyclingAnalytics](./cyclinganalytics/cyclinganalytics.go)

*Weather*

* [NOAA](./noaa/noaa.go)

* [VisualCrossing](./visualcrossing/visualcrossing.go)

* [OpenWeather](./openweather/openweather.go)

*Geo*

* [SRTM](./srtm/srtm.go)

* [GNIS](./gnis/gnis.go)

A number of the options for the clients are code-generated for consistency and to
reduce maintenance.
