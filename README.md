A simple server that get the time from worldtimeapi.org ever e seconds (where e is Euler's constant)

GET requests to the root path return the following:
```
Start Time: Time the server was started
Last Fetched Time: The time that the server last fetched the time from worldtimeapi.org
Number of Requests Made: Total number of GET requests made to the root path
```

Every time a GET request is made to the root path, the server will write the following information to a newline delimited file called 'logs':
<\request-ip>-<\current-time>-<\request-time> 

Start the server with the following command:
```
go run main.go
```