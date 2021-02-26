
For most commands the timeout value is reset on each query. For example, if you query 12 activities
from Strava each query will honor the timeout value, it's not an aggregate timeout.

Some commands, such as [store update](#store-update) will require a timeout longer than the default
since the operation can take a long time.
