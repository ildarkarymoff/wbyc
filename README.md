# wbyc
Weather by City API

## Usage

`http://<host>:8080/api/weather/current/{city}`

`city` parameter can be passed in cyrillic.

If service couldn't find city it returns 404 (both in response status and in response body).

## Deployment

Provide environment variables in `.env` file first:

```dotenv
YANDEX_API_KEY=<YOUR_YANDEX_GEOCODER_API_KEY>
APIXU_API_KEY=<YOUR_APIXU_API_KEY>
REDIS_HOST=db
```

Then run following commands:

```bash
docker-compose build
docker-compose up -d
``` 
