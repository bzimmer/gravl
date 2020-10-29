## Adventure Planning

## Examples

### WTA
This projects provides a structured JSON API for the [WTA](https://www.wta.org/) trip
reports. I was motivated to write it so I could more easily track the upvotes on my
trip reports.

*It is not endorsed by the WTA.*

Sum all the upvotes on my trip reports:

```sh
~/Development/src/github.com/bzimmer/gravl (master) > wta list bzimmer | jq 'map(.votes) | add | {"votes":.}'

12:04PM INF Please support the WTA build_version=5cd9716 url=https://www.wta.org/
{
  "votes": 94
}
```
