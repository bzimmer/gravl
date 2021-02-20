List all the available analyzers.

```sh
$ gravl pass list
{
	"ageride": {
		"base": false,
		"doc": "ageride returns all activities whose distance is greater than the athlete's age at the time of the activity",
		"flags": true
	},
	"benford": {
		"base": false,
		"doc": "benford returns the benford distribution of all the activities",
		"flags": false
	},
	"climbing": {
		"base": true,
		"doc": "climbing returns all activities exceeding the Golden Ratio - https://blog.wahoofitness.com/by-the-numbers-distance-and-elevation/",
		"flags": false
	},
	"cluster": {
		"base": false,
		"doc": "clusters returns the activities clustered by (distance, elevation) dimensions",
		"flags": true
	},
	"eddington": {
		"base": true,
		"doc": "eddington returns the Eddington number for all activities - The Eddington is the largest integer E, where you have cycled at least E miles (or kilometers) on at least E days",
		"flags": false
	},
	"festive500": {
		"base": true,
		"doc": "festive500 returns the activities and distance ridden during the annual #festive500 challenge - Thanks Rapha! https://www.rapha.cc/us/en_US/stories/festive-500",
		"flags": false
	},
	"forecast": {
		"base": false,
		"doc": "forecast the weather for an activity",
		"flags": false
	},
	"hourrecord": {
		"base": true,
		"doc": "hourrecord returns the longest distance traveled (in miles | kilometers) exceeding the average speed (mph | mps)",
		"flags": false
	},
	"koms": {
		"base": true,
		"doc": "koms returns all KOMs for the activities",
		"flags": false
	},
	"pythagorean": {
		"base": true,
		"doc": "pythagorean determines the activity with the highest pythagorean value defined as the sqrt(distance^2 + elevation^2) in meters",
		"flags": false
	},
	"rolling": {
		"base": true,
		"doc": "rolling returns the 'window' of activities with the highest accumulated distance.",
		"flags": true
	},
	"splat": {
		"base": false,
		"doc": "splat returns all activities in the units specified. This analyzer is useful for debugging the filter",
		"flags": false
	},
	"staticmap": {
		"base": false,
		"doc": "staticmap generates a staticmap for every activity",
		"flags": true
	},
	"totals": {
		"base": true,
		"doc": "totals returns the number of centuries (100 mi or 100 km)",
		"flags": false
	}
}
```
