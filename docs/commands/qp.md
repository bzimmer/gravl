```sh
$ gravl qp -e strava -u ca 4838740537
2021-02-25T20:53:13-08:00 INF exporter provider=strava
2021-02-25T20:53:14-08:00 INF uploader provider=ca
2021-02-25T20:53:14-08:00 INF export activityID=4838740537
2021-02-25T20:53:14-08:00 INF export activityID=4838740537 format=original
2021-02-25T20:53:15-08:00 INF upload activityID=4838740537
2021-02-25T20:53:16-08:00 INF status uploadID=8145126587
{
 "upload_id": 8145126587,
 "status": "processing",
 "ride_id": 0,
 "user_id": 1603533,
 "format": "fit",
 "datetime": "2021-02-26T04:53:18",
 "filename": "Blakely_Harbor.fit",
 "size": 62824,
 "error": "",
 "error_code": ""
}
2021-02-25T20:53:18-08:00 INF status uploadID=8145126587
{
 "upload_id": 8145126587,
 "status": "error",
 "ride_id": 0,
 "user_id": 1603533,
 "format": "fit",
 "datetime": "2021-02-26T04:53:18",
 "filename": "Blakely_Harbor.fit",
 "size": 62824,
 "error": "The ride already exists: 582750551527",
 "error_code": "duplicate_ride"
}
```
