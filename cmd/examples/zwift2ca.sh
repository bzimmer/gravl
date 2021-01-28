#!/bin/bash
set -e

num="${NUM_ACTIVITIES:-50}"

# If no activity is provided display the most recent rides
if [[ $# -eq 0 ]]
then
    gravl --timeout 1m strava activities -N $num | jq -s -c '.[] | select(.type == "VirtualRide") | [.id, .name, .start_date_local]'
    exit 0
fi

for arg in "$@"
do
gravl strava export -o -T "$arg.fit" -F fit $arg
gravl ca upload -p "$arg.fit"
rm -f "$arg.fit"
done
