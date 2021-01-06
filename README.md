
## Activity planning and analysis

### Activity
* [Cycling Analytics](https://cyclinganalytics.com)
* [Ride with GPS](https://ridewithgps.com)
* [Strava](https://strava.com)
* [WTA](https://wta.org)

### Geo
* [GNIS](https://geonames.usgs.gov)
* [SRTM](https://github.com/sakisds/go-srtm)

### Weather
* [NOAA](https://api.weather.gov)
* [OpenWeather API](https://openweathermap.org/api)
* [VisualCrossing](https://weather.visualcrossing.com)

### Examples

```sh
~/Development/src/github.com/bzimmer/gravl (cli) > go run cmd/gravl/* strava --export original 4569050661
"Paris.fit"
~/Development/src/github.com/bzimmer/gravl (cli) > go run cmd/gravl/* ca --upload Paris.fit
2021-01-04T20:38:12-08:00 INF uploading file=Paris.fit size=67629
{
 "status": "processing",
 "ride_id": 0,
 "user_id": 1603533,
 "format": "fit",
 "datetime": "2021-01-05T04:38:14",
 "upload_id": 7899891711,
 "filename": "Paris.fit",
 "size": 67629,
 "error": "",
 "error_code": ""
}
```
