## WTA Trip Report CLI

This projects provides a structured JSON API for the [WTA](https://www.wta.org/) trip
reports. I was motivated to write it so I could more easily track the upvotes on my
trip reports.

*It is not endorsed by the WTA.*

## Examples

Sum all the upvotes on my trip reports:

```sh
~/Development/src/github.com/bzimmer/wta (master) > wta list bzimmer | jq 'map(.votes) | add | {"votes":.}'

12:04PM INF Please support the WTA build_version=5cd9716 url=https://www.wta.org/
{
  "votes": 94
}
```

Show the top 15 most popular trip reports:

```sh
~/Development/src/github.com/bzimmer/wta (master) > wta l | jq -c 'sort_by(.votes)
    | reverse
    | map({title: .title, date: .hike_date, votes: .votes})
    | .[]' | head -15
11:58AM INF Please support the WTA build_version=7ddc0d9 url=https://www.wta.org/
{"title":"Blue Lake","date":"2020-10-14T00:00:00Z","votes":31}
{"title":"Lake Ingalls","date":"2020-10-15T00:00:00Z","votes":29}
{"title":"Chain Lakes Loop","date":"2020-10-15T00:00:00Z","votes":25}
{"title":"Lake Ingalls","date":"2020-10-15T00:00:00Z","votes":20}
{"title":"Snow Lake, Gem Lake","date":"2020-10-15T00:00:00Z","votes":20}
{"title":"Iron Bear - Teanaway Ridge, Jester Mountain","date":"2020-10-14T00:00:00Z","votes":18}
{"title":"Heather - Maple Pass Loop","date":"2020-10-15T00:00:00Z","votes":18}
{"title":"Blue Lake","date":"2020-10-15T00:00:00Z","votes":18}
{"title":"Lake Valhalla","date":"2020-10-16T00:00:00Z","votes":18}
{"title":"Lake Ingalls","date":"2020-10-15T00:00:00Z","votes":17}
{"title":"Heather - Maple Pass Loop","date":"2020-10-15T00:00:00Z","votes":14}
{"title":"Lake Ingalls","date":"2020-10-17T00:00:00Z","votes":14}
{"title":"Navaho Peak, Navaho Pass","date":"2020-10-15T00:00:00Z","votes":13}
{"title":"Gem Lake","date":"2020-10-15T00:00:00Z","votes":13}
{"title":"Yellow Aster Butte","date":"2020-10-15T00:00:00Z","votes":12}
```

These examples demonstrate combing the WTA trip report data with the very powerful
[jq](https://stedolan.github.io/jq/) to highlight the most popular current hikes. At
the time of writing, it's early fall and the clear winners in the hiking category are
those hikes showing off fall colors!

```sh
~/Development/src/github.com/bzimmer/wta (master) > dev/demo.sh
group by title, sum votes
{"title":"Lake Ingalls","votes":91}
{"title":"Blue Lake","votes":49}
{"title":"Heather - Maple Pass Loop","votes":37}
{"title":"Iron Bear - Teanaway Ridge, Jester Mountain","votes":33}
{"title":"Lake Valhalla","votes":24}
{"title":"Chain Lakes Loop","votes":24}
{"title":"Snow Lake, Gem Lake","votes":19}
{"title":"Kendall Katwalk","votes":19}
{"title":"Carne Mountain","votes":17}
{"title":"Mount Dickerman","votes":16}

sort by title, votes
{"title":"Blue Lake","date":"2020-10-14","votes":31}
{"title":"Lake Ingalls","date":"2020-10-15","votes":29}
{"title":"Chain Lakes Loop","date":"2020-10-15","votes":24}
{"title":"Lake Ingalls","date":"2020-10-15","votes":20}
{"title":"Snow Lake, Gem Lake","date":"2020-10-15","votes":19}
{"title":"Iron Bear - Teanaway Ridge, Jester Mountain","date":"2020-10-14","votes":18}
{"title":"Blue Lake","date":"2020-10-15","votes":18}
{"title":"Heather - Maple Pass Loop","date":"2020-10-15","votes":17}
{"title":"Lake Ingalls","date":"2020-10-15","votes":17}
{"title":"Lake Valhalla","date":"2020-10-16","votes":17}
```
