# Examples

The following examples demonstrate some of the flexibility provided by `gravl` coupled with common
unix tools like `jq`. Enjoy!

### Start the Strava OAuth server

Start the server, navigate to [localhost:9001](http://localhost:9001) and follow the directions.

```sh
$ gravl strava oauth
2021-02-26T08:20:08-08:00 INF serving address=0.0.0.0:9001
2021-02-26T08:20:15-08:00 INF request client_ip=[::1]:63528 elapsed=0.101966 method=GET path=/strava/auth/login status=302 user_agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36"
2021-02-26T08:20:35-08:00 INF request client_ip=[::1]:63528 elapsed=989.706809 method=GET path=/strava/auth/callback?state=yicOcreW0lRHMZrBYNrVzg%3D%3D&code=a13a5556bc906d8dde45a9717a87e4142b3e8668&scope=read,activity:write,activity:read_all,profile:read_all,read_all status=200 user_agent="Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36"
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

### Export the last three activities from Strava and upload them to CyclingAnalytics

*In this example all of the activities already exist, hence the error message, but the mechanics are the important bit*

```sh
$ gravl -c strava activities -N 3 | jq -r ".id" | xargs gravl -c qp -e strava -u ca
2021-02-26T08:12:01-08:00 INF do all=0 count=3 n=0 start=0 total=3
2021-02-26T08:12:03-08:00 INF do all=3 count=3 n=3 start=1 total=3
2021-02-26T08:12:03-08:00 INF exporter provider=strava
2021-02-26T08:12:04-08:00 INF uploader provider=ca
2021-02-26T08:12:04-08:00 INF export activityID=4844875821
2021-02-26T08:12:04-08:00 INF export activityID=4844875821 format=original
2021-02-26T08:12:05-08:00 INF upload activityID=4844875821
2021-02-26T08:12:07-08:00 INF status uploadID=2053926290
{"upload_id":2053926290,"status":"processing","ride_id":0,"user_id":1603533,"format":"fit","datetime":"2021-02-26T16:12:08","filename":"This_and_that.fit","size":274315,"error":"","error_code":""}
2021-02-26T08:12:09-08:00 INF status uploadID=2053926290
{"upload_id":2053926290,"status":"error","ride_id":0,"user_id":1603533,"format":"fit","datetime":"2021-02-26T16:12:08","filename":"This_and_that.fit","size":274315,"error":"The ride already exists: 230371166219","error_code":"duplicate_ride"}
2021-02-26T08:12:09-08:00 INF export activityID=4838740537
2021-02-26T08:12:09-08:00 INF export activityID=4838740537 format=original
2021-02-26T08:12:10-08:00 INF upload activityID=4838740537
2021-02-26T08:12:10-08:00 INF status uploadID=4434808581
{"upload_id":4434808581,"status":"processing","ride_id":0,"user_id":1603533,"format":"fit","datetime":"2021-02-26T16:12:12","filename":"Blakely_Harbor.fit","size":62824,"error":"","error_code":""}
2021-02-26T08:12:13-08:00 INF status uploadID=4434808581
{"upload_id":4434808581,"status":"error","ride_id":0,"user_id":1603533,"format":"fit","datetime":"2021-02-26T16:12:12","filename":"Blakely_Harbor.fit","size":62824,"error":"The ride already exists: 582750551527","error_code":"duplicate_ride"}
2021-02-26T08:12:13-08:00 INF export activityID=4832995501
2021-02-26T08:12:13-08:00 INF export activityID=4832995501 format=original
2021-02-26T08:12:14-08:00 INF upload activityID=4832995501
2021-02-26T08:12:15-08:00 INF status uploadID=6654126725
{"upload_id":6654126725,"status":"processing","ride_id":0,"user_id":1603533,"format":"fit","datetime":"2021-02-26T16:12:16","filename":"Panda_II.fit","size":109764,"error":"","error_code":""}
2021-02-26T08:12:17-08:00 INF status uploadID=6654126725
{"upload_id":6654126725,"status":"error","ride_id":0,"user_id":1603533,"format":"fit","datetime":"2021-02-26T16:12:16","filename":"Panda_II.fit","size":109764,"error":"The ride already exists: 686878103073","error_code":"duplicate_ride"}
```
