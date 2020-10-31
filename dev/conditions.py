forecast = {
    "wdir": 11.5,
    "temp": 39.2,
    "maxt": 45.6,
    "visibility": 15.0,
    "sunshine": 2.9,
    "wspd": 4.0,
    "datetimeStr": "2020-11-12T00:00:00-08:00",
    "heatindex": None,
    "cloudcover": 12.0,
    "pop": 33.2,
    "mint": 35.1,
    "datetime": 1605139200000,
    "precip": 0.0,
    "snowdepth": 0.0,
    "sealevelpressure": 1004.9,
    "sw_radiation": 94.6,
    "snow": 0.0,
    "dew": 32.8,
    "humidity": 78.6,
    "wgust": 12.3,
    "lw_radiation": 249.6,
    "conditions": "Clear",
    "windchill": 32.1
}

current = {
    "wdir": 130.0,
    "temp": 51.0,
    "sunrise": "2020-10-28T07:49:07-07:00",
    "visibility": None,
    "wspd": 0.9,
    "icon": "clear-day",
    "stations": "",
    "heatindex": None,
    "cloudcover": None,
    "datetime": "2020-10-28T12:45:10-07:00",
    "precip": 0.0,
    "moonphase": 0.46,
    "snowdepth": None,
    "sealevelpressure": None,
    "dew": 48.5,
    "sunset": "2020-10-28T17:57:49-07:00",
    "humidity": 91.0,
    "wgust": 4.0,
    "windchill": None
}

print("conditions")

fk = set(forecast.keys())
ck = set(current.keys())

print(sorted(fk-ck))
print(sorted(ck-fk))
print(sorted(fk.intersection(ck)))
