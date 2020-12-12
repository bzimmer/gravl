/*
Package wta provides a client to access trip reports submitted to the Washington Trails Association.

The (WTA) https://wta.org website is an amazing repository of trip reports in the Washington State
area. This package provides programmatic query access for those reports. It was originally written
to enable summing the upvotes across all submitted reports for a user.

Note, it is *not* endorsed by the WTA and because it is implemented by scraping the results of a request
it's prone to breaking.

Sum all the upvotes on my trip reports:

```sh
$ gravl wta bzimmer | jq -c 'map(.votes) | add | {"votes":.}'
{"votes": 94}
```
*/
package wta
