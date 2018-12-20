# gokit example

This project uses the gokit profilesvc example as a base but will become custom as time goes on. Build with go build, and run with go run cmd/server/main.go

Run the example with the optional port address for the service: 

```bash
$ go run ./cmd/profilesvc/main.go -http.addr :8080
ts=2018-05-01T16:13:12.849086255Z caller=main.go:47 transport=HTTP addr=:8080
```

Create a Profile:

```bash
$ curl -d '{"id":"1234","Name":"Go Kit"}' -H "Content-Type: application/json" -X POST http://localhost:8080/profiles/
{}
```

Get the profile you just created

```bash
$ curl localhost:8080/profiles/1234
{"profile":{"id":"1234","name":"Go Kit"}}
```
