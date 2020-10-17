# WTA API

This projects provides a structured JSON API for the WTA trip reports.

The primary motivation was making it easier to track the upvotes on my
trip reports.

# Examples

Show the top 15 most popular trip reports:

```sh
~/Downloads > wta l| jq -c "sort_by(.votes) | reverse | map({title: .title, date: .hike_date, votes: .votes}) | .[]" | head -15
8:35AM INF wta build_version=a9c7936
{"title":"Blue Lake","date":"2020-10-14T00:00:00Z","votes":31}
{"title":"Lake Ingalls","date":"2020-10-15T00:00:00Z","votes":29}
{"title":"Chain Lakes Loop","date":"2020-10-15T00:00:00Z","votes":24}
{"title":"Lake Ingalls","date":"2020-10-15T00:00:00Z","votes":20}
{"title":"Snow Lake, Gem Lake","date":"2020-10-15T00:00:00Z","votes":19}
{"title":"Iron Bear - Teanaway Ridge, Jester Mountain","date":"2020-10-14T00:00:00Z","votes":18}
{"title":"Blue Lake","date":"2020-10-15T00:00:00Z","votes":18}
{"title":"Heather - Maple Pass Loop","date":"2020-10-15T00:00:00Z","votes":17}
{"title":"Lake Ingalls","date":"2020-10-15T00:00:00Z","votes":17}
{"title":"Lake Valhalla","date":"2020-10-16T00:00:00Z","votes":17}
{"title":"Mount Dickerman","date":"2020-10-14T00:00:00Z","votes":16}
{"title":"Bullion Basin","date":"2020-10-14T00:00:00Z","votes":14}
{"title":"Heather - Maple Pass Loop","date":"2020-10-15T00:00:00Z","votes":14}
{"title":"Navaho Peak, Navaho Pass","date":"2020-10-15T00:00:00Z","votes":13}
{"title":"Gem Lake","date":"2020-10-15T00:00:00Z","votes":13}
```
