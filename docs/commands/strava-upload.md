An example exporting an activity to a local file and uploading back to Strava. This
example shows how the upload command can poll for status changes, in this case the
file is uploading a duplicate activity (of course!) and completes the polling.

*Note: check the results for any semantic errors*

```sh
$ gravl strava export -o 4838740537
2021-02-24T08:39:19-08:00 INF export activityID=4838740537 format=original
{
 "name": "Blakey_Harbor.fit",
 "format": "fit",
 "id": 4838740537
}
$ gravl strava upload -p Blakey_Harbor.fit
2021-02-24T08:39:36-08:00 INF collecting file=Blakey_Harbor.fit
2021-02-24T08:39:36-08:00 INF uploading file=Blakey_Harbor.fit
2021-02-24T08:39:37-08:00 INF status uploadID=5165766717
{
 "id": 5165766717,
 "id_str": "5165766717",
 "external_id": "",
 "error": "",
 "status": "Your activity is still being processed.",
 "activity_id": 0
}
2021-02-24T08:39:39-08:00 INF status uploadID=5165766717
{
 "id": 5165766717,
 "id_str": "5165766717",
 "external_id": "",
 "error": "Blakey_Harbor.fit duplicate of activity 4838740537",
 "status": "There was an error processing your activity.",
 "activity_id": 0
}
```
