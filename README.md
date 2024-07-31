# GoAQI

CLI program to display the current AQI (Air quality index) for a given city. Go to https://feeds.airnowapi.org/ to find your city and use the current air quality rss feed.

## Build

`go build -o bin/aqi`

## Run

`bin/aqi -url=https://feeds.airnowapi.org/rss/realtime/24.xml`

Example Results

``` sh
Cincinnati, OH - Current Air Quality
Moderate        62 AQI  PM2.5
```
