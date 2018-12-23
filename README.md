# gokit example

Build : go build
Run   : go run main.go

This service reads pokemon cards from https://api.pokemontcg.io/v1/cards and saves them to memory. Will be exposing endpoints to get the cards from memory. This should result in a service that is much faster than the existing one. It will not pick up any new cards as there are no events being fired to say when to clear out the memory or that something new has been added. This is just for fun and to explain the concepts of memory and micro services.

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
