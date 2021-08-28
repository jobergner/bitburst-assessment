# Assessment:

## How to Run:
1. Start Postgres
```bash
docker-compose up
```
2. Start Service
```bash
# uses localhost:9010 as object source
go run . -src=http://localhost:9010
```
3. Start Object Source
```bash
# in /objectsource/
go run .
```
The Services should now be able to communicate with another.

## Testing:
Feel free to mess around in `/integrationtest/main_test.go`.

```bash
go test ./...
```

## flags
| flag | description | default|
|------|------|-----|
|src|endpoint url for object source eg `http://localhost:9010`| uses mock when not specified|
|mock_db|whether to use the postgres db or an in-memory mock|false|
|ol| the object lifespan in seconds| 30 |


## project structure
| name | description |
|------|------|
|`/integrationtest`|contains an integration test (end2end) |
|`/objectsource`|the service I was provided with|
|`/pkg/get`|logic for fetching objects from remote or from a mock|
|`/pkg/handle`|logic for managing objects and their lifespans|
|`/pkg/object`|nothing but object structure|
|`/pkg/persist`|logic for object persistence in db or in-memory mock|
|`/pkg/server`|server start logic|
|`/.env`|environment variables read by postgres and service (committed for convenience)|
|`/docker-compose.yml`|starts dockerized postgres instance (`docker-compose up`)|
|`/init.sql`|describes table for db|

## issues:
When the object source service tried to make it's post requests to `/callback` it would reuse connections which were sometimes already closed by the peer, causing multiple errors:
```
Post http://localhost:9090/callback: read tcp 127.0.0.1:41312->127.0.0.1:9090: read: connection reset by peer
Post http://localhost:9090/callback: EOF
Post http://localhost:9090/callback: http: server closed idle connection
```
The only solution I found was to modify the code and close the request so the connection would not be reused. If there is a better solution please let me know.

# Task
Write a rest-service that listens on localhost:9090 for POST requests on /callback.
Run the go service attached to this task. It will send requests to your service
at a fixed interval of 5 seconds. The request body will look like this:

{
"object_ids": [1,2,3,4,5,6]
}

The amount of IDs varies with each request. Expect up to 200 IDs.

Every ID is linked to an object whose details can be fetched from the provided
service. Our service listens on localhost:9010/objects/:id and returns the
following response:

{
"id": <id>,
"online": true|false
}

Note that this endpoint has an unpredictable response time between 300ms and 4s!

Your task is to request the object information for every incoming object_idand
filter the objects by their "online" status.

Store all objects in a PostgreSQL database along with a timestamp when
the object was last seen.

Let your service delete objects in the database when they have not been
received for more than 30 seconds.

Important: due to business constraints we are not allowed to miss any callback to
our service. Write code in such a way that all errors are properly recovered
and that the endpoint is always available. Optimize for very high throughput
so that this service could work in production.

Bonus:

some comments in the code to explain the more complicated parts are appreciated
it a nice bonus if you provide some way to set up the things needed for us to
Test your code.
