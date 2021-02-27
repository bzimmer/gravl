# Examples

The following examples demonstrate some of the flexibility provided by `gravl` coupled with common
unix tools like `jq`. Enjoy!

### Start the Strava OAuth server

Start the server, navigate to [localhost:9001](http://localhost:9001) and follow the directions.

```sh
$ gravl strava oauth
2021-02-26T08:20:08-08:00 INF serving address=0.0.0.0:9001
2021-02-26T08:20:15-08:00 INF request client_ip=[::1]:63528 elapsed=0.101966 method=GET path=/strava/auth/login status=302 user_agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36"
2021-02-26T08:20:35-08:00 INF request client_ip=[::1]:63528 elapsed=989.706809 method=GET path=/strava/auth/callback?state=yicOcr...BYNrVzg%3D%3D&code=a13a5556bc906d8...42b3e8668&scope=read,activity:write,activity:read_all,profile:read_all,read_all status=200 user_agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36"
2021-02-26T08:20:35-08:00 INF request client_ip=[::1]:63528 elapsed=0.461866 method=GET path=/strava/auth/login status=302 user_agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36"
```

After being authenticated copy the `access_token` and `refresh_token` to the `gravl.yaml` configuration file.

### Download the last three activities from Strava and save them locally

```sh
$ gravl -c strava activities -N 3 | jq -r ".id" | xargs gravl -c strava export -o
2021-02-26T08:09:06-08:00 INF do all=0 count=3 n=0 start=0 total=3
2021-02-26T08:09:07-08:00 INF do all=3 count=3 n=3 start=1 total=3
2021-02-26T08:09:08-08:00 INF export activityID=4844875821 format=original
{"name":"This_and_that.fit","format":"fit","id":4844875821}
2021-02-26T08:09:08-08:00 INF export activityID=4838740537 format=original
{"name":"Blakely_Harbor.fit","format":"fit","id":4838740537}
2021-02-26T08:09:09-08:00 INF export activityID=4832995501 format=original
{"name":"Panda_II.fit","format":"fit","id":4832995501}
$ ls -als *.fit
128 -rw-r--r--  1 bzimmer  staff    61K Feb 26 08:09 Blakely_Harbor.fit
216 -rw-r--r--  1 bzimmer  staff   107K Feb 26 08:09 Panda_II.fit
536 -rw-r--r--  1 bzimmer  staff   268K Feb 26 08:09 This_and_that.fit
```

### Export the last four activities from Zwift and upload them to CyclingAnalytics

*In this example all of the activities already exist, hence the error message, but the mechanics are the important bit*

*Note the use of `.id_str` with `jq` because the Zwift id is [not properly parsed](https://github.com/stedolan/jq/issues/369) by `jq`*

```sh
$ gravl -c zwift activities -N 4 | jq -r ".id_str" | xargs gravl -c qp -e zwift -u ca
2021-02-26T13:12:25-08:00 INF do all=0 count=4 n=0 start=0 total=4
2021-02-26T13:12:25-08:00 INF do all=4 count=4 n=4 start=1 total=4
2021-02-26T13:12:25-08:00 INF exporter provider=zwift
2021-02-26T13:12:26-08:00 INF uploader provider=ca
2021-02-26T13:12:26-08:00 INF export activityID=752165469668926976
2021-02-26T13:12:26-08:00 INF export activityID=736562609413586256
2021-02-26T13:12:26-08:00 INF export activityID=743809551513814416
2021-02-26T13:12:26-08:00 INF export activityID=752890649715818240
2021-02-26T13:12:26-08:00 INF export activityID=752890649715818240 filename=2021-02-18-06-56-16.fit
2021-02-26T13:12:26-08:00 INF upload activityID=752890649715818240
2021-02-26T13:12:26-08:00 INF export activityID=736562609413586256 filename=2021-01-26-18-15-16.fit
2021-02-26T13:12:26-08:00 INF upload activityID=736562609413586256
2021-02-26T13:12:26-08:00 INF export activityID=743809551513814416 filename=2021-02-05-18-13-46.fit
2021-02-26T13:12:26-08:00 INF upload activityID=743809551513814416
2021-02-26T13:12:26-08:00 INF export activityID=752165469668926976 filename=2021-02-17-06-55-29.fit
2021-02-26T13:12:26-08:00 INF upload activityID=752165469668926976
2021-02-26T13:12:27-08:00 INF status uploadID=1289862592
2021-02-26T13:12:27-08:00 INF status uploadID=8033493930
2021-02-26T13:12:27-08:00 INF status uploadID=3117629228
{"upload_id":1289862592,"status":"processing","ride_id":0,"user_id":1603533,"format":"fit","datetime":"2021-02-26T21:12:29","filename":"2021-02-18-06-56-16.fit","size":79806,"error":"","error_code":""}
{"upload_id":8033493930,"status":"processing","ride_id":0,"user_id":1603533,"format":"fit","datetime":"2021-02-26T21:12:29","filename":"2021-02-17-06-55-29.fit","size":93239,"error":"","error_code":""}
{"upload_id":3117629228,"status":"processing","ride_id":0,"user_id":1603533,"format":"fit","datetime":"2021-02-26T21:12:29","filename":"2021-02-05-18-13-46.fit","size":114921,"error":"","error_code":""}
2021-02-26T13:12:27-08:00 INF status uploadID=5906692136
{"upload_id":5906692136,"status":"processing","ride_id":0,"user_id":1603533,"format":"fit","datetime":"2021-02-26T21:12:29","filename":"2021-01-26-18-15-16.fit","size":113287,"error":"","error_code":""}
2021-02-26T13:12:29-08:00 INF status uploadID=1289862592
2021-02-26T13:12:29-08:00 INF status uploadID=8033493930
2021-02-26T13:12:29-08:00 INF status uploadID=3117629228
{"upload_id":1289862592,"status":"error","ride_id":0,"user_id":1603533,"format":"fit","datetime":"2021-02-26T21:12:29","filename":"2021-02-18-06-56-16.fit","size":79806,"error":"The ride already exists: 801167646702","error_code":"duplicate_ride"}
{"upload_id":8033493930,"status":"error","ride_id":0,"user_id":1603533,"format":"fit","datetime":"2021-02-26T21:12:29","filename":"2021-02-17-06-55-29.fit","size":93239,"error":"The ride already exists: 509770415934","error_code":"duplicate_ride"}
{"upload_id":3117629228,"status":"error","ride_id":0,"user_id":1603533,"format":"fit","datetime":"2021-02-26T21:12:29","filename":"2021-02-05-18-13-46.fit","size":114921,"error":"The ride already exists: 336104747302","error_code":"duplicate_ride"}
2021-02-26T13:12:30-08:00 INF status uploadID=5906692136
```

The go implementation of `jq`, [`gojq`](https://github.com/itchyny/gojq), does handle the ids correctly.

```sh
$ gravl -c zwift activities -N 4 | jq -r ".id"
2021-02-26T13:24:44-08:00 INF do all=0 count=4 n=0 start=0 total=4
2021-02-26T13:24:44-08:00 INF do all=4 count=4 n=4 start=1 total=4
752890649715818200
752165469668927000
743809551513814400
736562609413586300
$ gravl -c zwift activities -N 4 | gojq -r ".id"
2021-02-26T13:24:49-08:00 INF do all=0 count=4 n=0 start=0 total=4
2021-02-26T13:24:49-08:00 INF do all=4 count=4 n=4 start=1 total=4
752890649715818240
752165469668926976
743809551513814416
736562609413586256
```
