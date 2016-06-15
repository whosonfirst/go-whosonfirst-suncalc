# go-whosonfirst-suncalc

## Usage

The short version is "Too soon". But read on if you're feeling adventurous.

## Starting the server

```
$> ./bin/wof-suncalc-server
start and listen for requests at http://localhost:8080
```

## Starting the server with a local WOF data 

_This will become clear in a moment._

```
$> ./bin/wof-suncalc-server -wof-root file:///usr/local/mapzen/whosonfirst-data/data
start and listen for requests at http://localhost:8080
```

## Talking to the server

```
$> curl -s 'http://localhost:8080?latitude=40.677524&longitude=-73.987343&ymd=19700101' | python -mjson.tool
{
    "dawn": "1024-07-02T00:49:20-08:00",
    "dusk": "1024-07-02T17:14:22-08:00",
    "goldenHour": "1024-07-02T15:54:25-08:00",
    "goldenHourEnd": "1024-07-02T02:09:17-08:00",
    "nadir": "1024-07-01T21:01:51-08:00",
    "nauticalDawn": "1024-07-02T00:02:45-08:00",
    "nauticalDusk": "1024-07-02T18:00:58-08:00",
    "night": "1024-07-02T18:58:32-08:00",
    "nightEnd": "1024-07-01T23:05:10-08:00",
    "solarNoon": "1024-07-02T09:01:51-08:00",
    "sunrise": "1024-07-02T01:25:15-08:00",
    "sunriseEnd": "1024-07-02T01:28:48-08:00",
    "sunset": "1024-07-02T16:38:27-08:00",
    "sunsetStart": "1024-07-02T16:34:54-08:00"
}
```

## Talking to the server with Who's On First IDs

```
$> curl -s 'http://localhost:8080?wofid=85784831&ymd=10240703' | python -mjson.tool
{
    "dawn": "1024-07-02T00:49:20-08:00",
    "dusk": "1024-07-02T17:14:22-08:00",
    "goldenHour": "1024-07-02T15:54:25-08:00",
    "goldenHourEnd": "1024-07-02T02:09:17-08:00",
    "nadir": "1024-07-01T21:01:51-08:00",
    "nauticalDawn": "1024-07-02T00:02:45-08:00",
    "nauticalDusk": "1024-07-02T18:00:58-08:00",
    "night": "1024-07-02T18:58:32-08:00",
    "nightEnd": "1024-07-01T23:05:10-08:00",
    "solarNoon": "1024-07-02T09:01:51-08:00",
    "sunrise": "1024-07-02T01:25:15-08:00",
    "sunriseEnd": "1024-07-02T01:28:48-08:00",
    "sunset": "1024-07-02T16:38:27-08:00",
    "sunsetStart": "1024-07-02T16:34:54-08:00"
}
```
