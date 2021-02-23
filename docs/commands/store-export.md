`gravl` allows flexible exporting of Strava activities from the local store by the `export` command. As an example
of exporting a subset of activities:

```sh
$ gravl -c store export -f ".Type == 'NordicSki'" | wc -l
2021-02-20T19:12:07-08:00 INF bunt db path="/Users/bzimmer/Library/Application Support/com.github.bzimmer.gravl/gravl.db"
2021-02-20T19:12:08-08:00 INF export activities=46 elapsed=678.191786
46
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
