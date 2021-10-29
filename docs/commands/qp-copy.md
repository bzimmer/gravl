```sh
$ gravl qp copy --from zwift --to cyclinganalytics 934398333398662432
2021-10-28T06:25:55-07:00 INF created zwift client
2021-10-28T06:25:55-07:00 INF created cyclinganalytics client
2021-10-28T06:25:56-07:00 INF export exp=2021-10-26-18-19-39.fit id=934398333398662432
2021-10-28T06:25:57-07:00 INF poll id=1067278367 iteration=0
2021-10-28T06:26:00-07:00 INF poll id=1067278367 iteration=1
2021-10-28T06:26:00-07:00 INF counters count=1 metric=gravl.upload.file.success
2021-10-28T06:26:00-07:00 INF counters count=2 metric=gravl.upload.poll
2021-10-28T06:26:00-07:00 INF samples count=1 max=4.181048393249512 mean=4.181048393249512 metric=gravl.runtime min=4.181048393249512 stddev=0
```

```sh
$ gravl qp copy --from zwift --to cyclinganalytics 934398333398662432
2021-10-28T06:34:06-07:00 INF created zwift client
2021-10-28T06:34:06-07:00 INF created cyclinganalytics client
2021-10-28T06:34:07-07:00 INF export exp=2021-10-26-18-19-39.fit id=934398333398662432
2021-10-28T06:34:08-07:00 INF poll id=9819686356 iteration=0
{
 "upload_id": 9819686356,
 "status": "processing",
 "ride_id": 0,
 "user_id": 1603544,
 "format": "fit",
 "datetime": "2021-10-28T13:34:09",
 "filename": "2021-10-26-18-19-39.fit",
 "size": 113324,
 "error": "",
 "error_code": ""
}
2021-10-28T06:34:10-07:00 INF poll id=9819686356 iteration=1
{
 "upload_id": 9819686356,
 "status": "error",
 "ride_id": 0,
 "user_id": 1603544,
 "format": "fit",
 "datetime": "2021-10-28T13:34:09",
 "filename": "2021-10-26-18-19-39.fit",
 "size": 113324,
 "error": "The ride already exists: 894091983723",
 "error_code": "duplicate_ride"
}
2021-10-28T06:34:10-07:00 INF counters count=2 metric=gravl.upload.poll
2021-10-28T06:34:10-07:00 INF counters count=1 metric=gravl.upload.file.success
2021-10-28T06:34:10-07:00 INF samples count=1 max=3.807884454727173 mean=3.807884454727173 metric=gravl.runtime min=3.807884454727173 stddev=0
```
