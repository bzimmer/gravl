If neither `-o` or `-O` are specified the contents of the file are written to stdout.
If `-o` is specified, the file will be written to disk using the name provided by Strava, even if it already exists locally.
If `-O` is specified, the file will be written to disk using the name provided by the flag. It will not overwrite an existing
file unless `-o` was also specified.
