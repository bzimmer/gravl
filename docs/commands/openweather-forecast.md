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
