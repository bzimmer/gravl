#!/bin/bash
set -e

for arg in "$@"
do
gravl strava export -t "$arg.fit" -f fit $arg
gravl ca upload -p "$arg.fit"
rm -f "$arg.fit"
done
