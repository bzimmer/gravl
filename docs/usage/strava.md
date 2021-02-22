The Strava client is comprised of general [API](https://developers.strava.com/) access supporting
activites, routes, and streams as well as some functionality available by scraping the website as
inspired by [stravaweblib](https://github.com/pR0Ps/stravaweblib).

Additionally, there's full support for implementing `webhooks` but only only webhook management is
available via the commandline (eg [`strava webhook list`](#strava-webhook-list),
[`strava webhook subscribe`](#strava-webhook-subscribe), and [`strava webhook unsubscribe`](#strava-webhook-unsubscribe)).

The entire [`analysis`](#strava-analysis) package is built around Strava activities.


