# whereis

This is what that powers `whereis.tjhorner.nyc`. The code is not very good.

## Build

### Bare Metal

Requirements:
- Go 1.13+
- Some recent version of Node

Then run:

```shell
cd frontend && npm install && cd ..
make build
```

### Docker

There's a Dockerfile. No build args are required.

## Running

### Environment Variables

- `WHEREIS_SHARED_KEY`: Shared key between server and client. Make it random!
- `WHEREIS_LISTEN_ADDRESS`: Where to listen (default: `:3000`)
- `WHEREIS_DB_PATH`: Where the SQLite database should be (default: `./whereis.db3`)
- `WHEREIS_LIQ_KEY`: API key provided by [LocationIQ](https://locationiq.com)

### Setup

You need some sort of client that can send your location as a POST request to a server periodically.

If you're on Android, there's [this app](https://play.google.com/store/apps/details?id=com.mendhak.gpslogger) that can do exactly that. To set it up, follow these instructions:

1. Set up the server somewhere.
2. Go to the hamburger menu > Logging details > Log to custom URL
3. Set these variables accordingly:
  - **URL:** `https://your.host.name/api/v1/location`
  - **HTTP Body:** `lat=%LAT&lon=%LON&batt=%BATT&acc=%ACC`
  - **HTTP Headers:**
    ```
    Authorization: Bearer [YOUR_SHARED_KEY]
    Content-Type: application/x-www-form-urlencoded
    ```
  - **HTTP Method:** `POST`

Turn on logging and see if it works!

If you're on iOS, you're on your own... sorry. But the above should help you get started.

## License

MIT