
### Documentation

* [API Documentation](https://www.visualcrossing.com/resources/documentation/weather-api/weather-api-documentation/)

### Surprises

1. Inconsistent use of `datetime` in forecast and current conditions
   1. In forecast conditions `datetime` is an `int64`
      1. The field `datetimeStr` is usable for easy parsing
   2. In current conditions `datetime` is a `string`
2. HTTP status code is always 200, even for an error
   1. The format for an error is not documented but in testing appears to be a simple dictionary

### Conditions

Difference between `current` and `forecast` conditions:

```py
# in forecast, not current
['conditions', 'datetimeStr', 'lw_radiation', 'maxt', 'mint', 'pop', 'snow', 'sunshine', 'sw_radiation']
# in current, not forecast
['icon', 'moonphase', 'stations', 'sunrise', 'sunset']
# in both
['cloudcover', 'datetime', 'dew', 'heatindex', 'humidity', 'precip', 'sealevelpressure', 'snowdepth', 'temp', 'visibility', 'wdir', 'wgust', 'windchill', 'wspd']
```
