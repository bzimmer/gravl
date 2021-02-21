Query the [GNIS](https://www.usgs.gov/core-science-systems/national-geospatial-program/geographic-names) database for US States.

This functionality was added mainly to get pseudo-accurate coordinates for use in querying weather forecasts.

```sh
$ gravl -c gnis NH | wc -l
   14770
$ gravl -c gnis NH | head -5
{"type":"Feature","id":"205110","geometry":{"type":"Point","coordinates":[-77.0775473,40.3221113,200]},"properties":{"class":"Trail","locale":"PA","name":"North Country National Scenic Trail","source":"https://geonames.usgs.gov"}}
{"type":"Feature","id":"206425","geometry":{"type":"Point","coordinates":[-72.3331382,41.2723203,0]},"properties":{"class":"Stream","locale":"CT","name":"Connecticut River","source":"https://geonames.usgs.gov"}}
{"type":"Feature","id":"561049","geometry":{"type":"Point","coordinates":[-71.0306287,44.9383838,384]},"properties":{"class":"Stream","locale":"ME","name":"Abbott Brook","source":"https://geonames.usgs.gov"}}
{"type":"Feature","id":"561428","geometry":{"type":"Point","coordinates":[-70.7517197,43.0823107,9]},"properties":{"class":"Island","locale":"ME","name":"Badgers Island","source":"https://geonames.usgs.gov"}}
{"type":"Feature","id":"561491","geometry":{"type":"Point","coordinates":[-70.9716939,43.6157566,168]},"properties":{"class":"Reservoir","locale":"ME","name":"Balch Pond","source":"https://geonames.usgs.gov"}}
```

I might typically use it like this:

```sh
$ gravl -c gnis WA | grep "Barlow Pass"
{"type":"Feature","id":"1516141","geometry":{"type":"Point","coordinates":[-121.4440005,48.0264959,721]},"properties":{"class":"Gap","locale":"WA","name":"Barlow Pass","source":"https://geonames.usgs.gov"}}
$ gravl -c gnis WA | grep "Barlow Pass" | jq ".geometry | .coordinates"
[
  -121.4440005,
  48.0264959,
  721
]
```
