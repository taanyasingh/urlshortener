# urlshortener
URL shortener service

URL shortner is a golang application that shortens the original url using go fiber  framework with redis memory cache.
The shorten url will redirect to original URL if opened in the browser.

__Prerquisites__
- go version go1.21.3 linux/amd64
- Docker
- Redis image
  
__Steps to setup on local__
- Run redis server using docker<br>
      ```
      $ docker run --name <container-name> -d redis
      ```
- To run the application locally run redis server using above command and the run server with : <br>
       ```
      cd urlshortener/api/
      ``` <br>
      ```
     go run main.go
      ```

__URL shorten API :__
http://127.0.0.1:3000/api/v1  <br>
body : <br>
```javascript
{
    "url" : "input url", 
    "short" : "custom url for shortening",
    "expiry" : time of expiry for shorten url 
}
```

__Redirect API__
Copy the shorten url and run it in browser, link will get redirected to original url.

__Metric API__
http://127.0.0.1:3000/api/v1/metrics

