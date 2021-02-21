Strava export uses the website and therefore requires a username and password instead of the usual oauth credentials.

If neither `-o` or `-O` are specified the contents of the file are written to stdout.
If `-o` is specified, the file will be written to disk using the name provided by Strava, even if it already exists locally.
If `-O` is specified, the file will be written to disk using the name provided by the flag. It will not overwrite an existing
file unless `-o` was also specified.

```sh
$ gravl strava export -o 4814450574
2021-02-20T09:20:29-08:00 INF export activityID=4814540574 format=original
{
 "id": 4814450574,
 "name": "Friday.fit",
 "format": "fit"
}
$ ls -las Friday.fit
56 -rw-r--r--  1 bzimmer  staff    25K Feb 20 09:20 Friday.fit
```

An example of the overwrite logic.

```sh
$ gravl strava export -O Friday.fit 4814540547
2021-02-20T09:24:44-08:00 INF export activityID=4814540547 format=original
2021-02-20T09:24:45-08:00 ERR file exists and -o flag not specified filename=Friday.fit
2021-02-20T09:24:45-08:00 ERR gravl strava error="file already exists"
```

It's also possible to use the attribute functionality by specifying one or more attributes using the `-B` flag. In this
example we export only those activities of type `Ride`, extract their distance in miles, and use standard unix tools to
create the top 10 rides by distance.

```sh
$ gravl -c store export -f ".Type == 'Ride'" -B ".Distance.Miles()" | jq ".[]" | sort -nr | head -10
2021-02-20T18:56:27-08:00 INF bunt db path="/Users/bzimmer/Library/Application Support/com.github.bzimmer.gravl/gravl.db"
2021-02-20T18:56:28-08:00 INF export activities=618 elapsed=682.669242
161.20357114451602
107.90359301678198
101.80421339378032
99.08571442774199
97.57764654418197
90.31630279169649
84.45552970651396
83.14567923327766
79.88223773164718
76.653593016782
```
