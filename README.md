# StoikUrlShortenerAPI

This URL shortener is a micro-service to shorten long URLs and to handle the redirection by generated short URLs.

It exposes an API with 2 endpoints : one to generate a short URL from an orignal one, the domain is short.io .
The second one handles the redirection from the generated short URL.

The service requires Postgres database connection.

## Request for short URL

URL: <host>[:<port>]/api/url/shorten

Method: POST

Request body: JSON with following parameters:

    long_url: string, URL to shorten, mandatory

Success response: HTTP 200 OK with body containing JSON with following parameters:

    short_url: string, short URL

## Redirect to long URL

URL: <host>[:<port>]/<short_url> - URL from response on request for short URL

Method: GET

Response contain the redirection to long URL (response code: HTTP 301 = long URL in response header)

This can be tested via browser with for example : http://localhost:8080/short.io/22e4xIF

## Run the demo

We should be able to simply run the docker-compose file given you have docker + docker-compose intalled on your machine.

There are unit tests in the project, they can be run independantly, they are also run during the docker image build.

- docker-compose up --build

## Improvements

What we can do to further improve this project before it goes to production :

- Add a check to see if the original URL already exists, if yes, return the corresponding entry
- Put the requests to init the database tables in dedicated scripts (like Sqitch)
- Add monitoring
- Write unit tests for the repository part