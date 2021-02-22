
## *ageride*

**Description**

ageride returns all activities whose distance is greater than the athlete's age at the time of the activity

**Flags**

|Flag|Default|Description|
|-|-|-|
|```birthday```|```"0001-01-01"```|the athlete's birthday in `YYYY-MM-DD` format|

## *benford*

**Description**

benford returns the benford distribution of all the activities

## *climbing*

**Description**

climbing returns all activities exceeding the Golden Ratio

https://blog.wahoofitness.com/by-the-numbers-distance-and-elevation/

## *cluster*

**Description**

clusters returns the activities clustered by (distance, elevation) dimensions

**Flags**

|Flag|Default|Description|
|-|-|-|
|```clusters```|```4```|number of clusters|
|```threshold```|```0.01```|threshold (in percent between 0.0 and 0.1) aborts processing if less than n% of data points shifted clusters in the last iteration|

## *eddington*

**Description**

eddington returns the Eddington number for all activities

The Eddington is the largest integer E, where you have cycled at least E miles (or kilometers) on at least E days

## *festive500*

**Description**

festive500 returns the activities and distance ridden during the annual #festive500 challenge

Only the activity types 'Ride', 'VirtualRide', and 'Handcycle' are considered.

Thanks Rapha! https://www.rapha.cc/us/en_US/stories/festive-500

## *forecast*

**Description**

forecast the weather for an activity

## *hourrecord*

**Description**

hourrecord returns the longest distance traveled (in miles | kilometers) exceeding the average speed (mph | mps)

## *koms*

**Description**

koms returns all KOMs for the activities

## *pythagorean*

**Description**

pythagorean determines the activity with the highest pythagorean value defined as the sqrt(distance^2 + elevation^2) in meters

## *rolling*

**Description**

rolling returns the `window` of activities with the highest accumulated distance.

**Flags**

|Flag|Default|Description|
|-|-|-|
|```window```|```7```|the number of days in the window|

## *splat*

**Description**

splat returns all activities in the units specified

This analyzer is useful for debugging the filter

## *staticmap*

**Description**

staticmap generates a staticmap for every activity

**Flags**

|Flag|Default|Description|
|-|-|-|
|```workers```|```15```|number of workers|

## *totals*

**Description**

totals returns the number of centuries (100 mi or 100 km)
