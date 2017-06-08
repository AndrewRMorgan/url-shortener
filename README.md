# FreeCodeCamp API Project: URL Shortener Microservice
## User stories:
1. I can pass a URL as a parameter and I will receive a shortened URL in the JSON response.
2. If I pass an invalid URL that doesn't follow the valid http://www.example.com format, the JSON response will contain an error instead.
3. When I visit that shortened URL, it will redirect me to my original link.

## Example creation usage:

```js
https://morning-retreat-24523.herokuapp.com/new/https://www.google.com 
```

## Example creation output:

```js
{ "original_url":"https://www.google.com", "short_url":"https://morning-retreat-24523.herokuapp.com/get/3578" }
```

## Usage:

`https://morning-retreat-24523.herokuapp.com/get/3578`

### Will redirect to:

`https://www.google.com/`
