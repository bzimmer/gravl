#!/bin/sh

echo "group by title, sum votes"
jq -c 'group_by(.title)
    | map({title: .[0].title, votes: map(.votes) | add})
    | sort_by(.votes)
    | reverse
    | .[]' testdata/wta_list.json | head -10

echo

echo "sort by title, votes"
jq -c 'sort_by(.votes)
    | reverse
    | map({title: .title, date: .hike_date[0:10], votes: .votes})
    | .[]' testdata/wta_list.json | head -10
