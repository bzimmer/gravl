
## Activity planning and analysis

![gopher](./gravl.png)

### Activity
* [Cycling Analytics](https://www.cyclinganalytics.com/)
* [Ride with GPS](https://ridewithgps.com)
* [Strava](https://strava.com)
* [WTA](https://wta.org)

### Geo
* [GNIS](https://geonames.usgs.gov)
* [SRTM](https://github.com/sakisds/go-srtm)

### Weather
* [NOAA](https://weather.gov)
* [OpenWeather API](https://openweathermap.org/api)
* [VisualCrossing](https://visualcrossing.com)

### Examples

```sh
gravl > go run cmd/gravl/* strava --export original 4569050661
"Paris.fit"
gravl > go run cmd/gravl/* ca --upload Paris.fit
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

### Credits

Credits to Renee French for the gophers.
