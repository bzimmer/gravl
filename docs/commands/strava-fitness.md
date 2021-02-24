Queries the fitness and freshness data for the authenticated user.

```sh
$ gravl strava fitness
[
 {
  "date": {
   "year": 2020,
   "month": 8,
   "day": 22
  },
  "fitness_profile": {
   "fitness": 107.39025265712681,
   "impulse": 83,
   "relative_effort": 64,
   "fatigue": 86.97505302175927,
   "form": 20.415199635367543
  },
  "activities": [
   {
    "id": 3951687537,
    "impulse": 83,
    "relative_effort": 64
   }
  ]
 },
 ...,
 {
  "date": {
   "year": 2021,
   "month": 3,
   "day": 8
  },
  "fitness_profile": {
   "fitness": 56.367803371875894,
   "impulse": 0,
   "relative_effort": 0,
   "fatigue": 5.510876434392897,
   "form": 50.856926937482996
  },
  "activities": []
 }
]
