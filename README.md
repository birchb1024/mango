# mango ![](mango.png)

The fruitful web server.

A very simple web server using the out-of-the box Golang HTTP server. 

# Usage

```
$ mango <port number> <root directory>
```

## Example

```
$ mango 8080 $PWD
```

## CGI

Executable files in the `cgi-bin/` directory are run according to the CGI standard as implemented
by the Golang net/http/cgi module.

Invoke the dumping CGI to see the variables:
```
curl -v http://localhost:7777/cgi-bin/env.html.sh
```

Example: Invoke a CGI to render the date.
```
curl http://localhost:7777/cgi-bin/date.sh
```

Debug: See what curl is sending by targeting nc server
```
nc -l -p 6666 & curl -v --ipv4 --http1.0 -X POST -H 'Content-Type: text/plain' --data-binary @.gitignore -sS http://localhost:6666/cgi-bin/echo.sh 2>&1
```


Icons made by <a href="https://www.flaticon.com/authors/freepik" title="Freepik">Freepik</a> from <a href="https://www.flaticon.com/" title="Flaticon"> www.flaticon.com</a>
