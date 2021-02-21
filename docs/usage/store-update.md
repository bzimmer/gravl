In order to have a more performant experience when running analyzers all Strava activities
are stored locally in a `Store` implementation. The default `Store` is an implementation
using `buntdb` as it allows very simple, fast, and durable local storage though other
implementations exist.

Updates from Strava are incremental and should be run periodically to get the latest activities.

*Note: if the activity already exists locally `gravl` will not update it, it will need to be removed and updated*

```sh
$ gravl store update
2021-02-20T15:59:25-08:00 INF bunt db path="/Users/bzimmer/Library/Application Support/com.github.bzimmer.gravl/gravl.db"
2021-02-20T15:59:26-08:00 INF do all=0 count=100 n=0 start=0 total=0
2021-02-20T15:59:29-08:00 INF do all=100 count=100 n=100 start=1 total=0
2021-02-20T15:59:29-08:00 INF querying activity details ID=4819927284
2021-02-20T15:59:30-08:00 INF saving activity details ID=4819927284 n=1 name="Morning Ride"
2021-02-20T15:59:30-08:00 INF querying activity details ID=4814540574
2021-02-20T15:59:30-08:00 INF saving activity details ID=4814540574 n=2 name="Afternoon Ride"
2021-02-20T15:59:31-08:00 INF do all=200 count=100 n=100 start=2 total=0
2021-02-20T15:59:34-08:00 INF do all=300 count=100 n=100 start=3 total=0
2021-02-20T15:59:36-08:00 INF do all=400 count=100 n=100 start=4 total=0
2021-02-20T15:59:39-08:00 INF do all=500 count=100 n=100 start=5 total=0
2021-02-20T15:59:41-08:00 INF do all=600 count=100 n=100 start=6 total=0
2021-02-20T15:59:43-08:00 INF do all=700 count=100 n=100 start=7 total=0
2021-02-20T15:59:48-08:00 INF do all=800 count=100 n=100 start=8 total=0
2021-02-20T15:59:51-08:00 INF do all=900 count=100 n=100 start=9 total=0
2021-02-20T15:59:54-08:00 INF do all=1000 count=100 n=100 start=10 total=0
2021-02-20T15:59:56-08:00 INF do all=1100 count=100 n=100 start=11 total=0
2021-02-20T15:59:58-08:00 INF do all=1200 count=100 n=100 start=12 total=0
2021-02-20T15:59:59-08:00 INF do all=1200 count=100 n=0 start=13 total=0
{
 "existing": 1198,
 "new": 2,
 "total": 1200
}
```

The results of the command show the number of new, existing, and total activities stored locally.
