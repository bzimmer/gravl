List all the available analyzers.

```sh
$ gravl pass list
{
	"ageride": {
		"base": false,
		"doc": "ageride returns all activities whose distance is greater than the athlete's age at the time of the activity",
		"flags": true
	},
	...,
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
	...,
	"totals": {
		"base": true,
		"doc": "totals returns the number of centuries (100 mi or 100 km)",
		"flags": false
	}
}
```
