//go:generate go run ../../cmd/genreadme/genreadme.go

/*
Package visualcrossing provides a client to access VisualCrossing weather forecasts.

Documentation

* (API documentation) https://www.visualcrossing.com/resources/documentation/weather-api/weather-api-documentation/

Notes

* Inconsistent use of `datetime` in forecast and current conditions

In forecast conditions `datetime` is an `int64`, the field `datetimeStr` is usable
for easy parsing. In current conditions `datetime` is a `string`.

* HTTP status code is always 200, even for an error

The format for an error is not documented but in testing appears to be a simple dictionary.

Conditions

Difference between `current` and `forecast` conditions:

```py
# in forecast, not current
['conditions', 'datetimeStr', 'lw_radiation', 'maxt', 'mint', 'pop', 'snow', 'sunshine', 'sw_radiation']
# in current, not forecast
['icon', 'moonphase', 'stations', 'sunrise', 'sunset']
# in both
['cloudcover', 'datetime', 'dew', 'heatindex', 'humidity', 'precip', 'sealevelpressure', 'snowdepth', 'temp', 'visibility', 'wdir', 'wgust', 'windchill', 'wspd']
```
*/
package visualcrossing
