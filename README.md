#Bitly Backend Coding Challenge Submission

An http server implementation of the Bitly REST API
```GET /groups/{groupGuid}/countries/averages```
which returns the average number of clicks of all the Bitlinks in the user's
default group over a time interval (default 30 days) by country.

## Getting Started

Download the project into your desired directory. All remaining instructions
assume that the commands are run from the root directory of the project.  

##Prerequisites

It is recommended that you have Docker installed as well as the latest Go version
(the app was tested against Go 1.13).  The resulting docker image is fairly large
at ~774MB just as a heads up.  

##Running the Tests
All the tests live in the main package so it suffices to execute from the root project directory:
```go test
```

There are 3 levels of tests.  At the highest level are the HTTP tests, which test the response codes
based on various input conditions and mock the handler internal dependencies.
Next are the handler tests (aka integration tests) which exercise the logic that
aggregates the data from the API and does the computation for the response. These tests
mock the API dependencies.  Third are the API tests which exercise the API functionality
and mock the over the wire dependencies.

##Installing

You should first run the automated tests before installing.  See above section.

You will build a Docker image containing the go executable and then use Docker
to run the image in a container.

Navigate into the root directory of the project.  

To build the application, execute the command

```
docker build -t <my_image_name> .
```
In the below examples, the value of ```<my_image_name>``` is ```bitlymetrics```

You should now have a Docker image called ```bitlymetrics``` in your Docker repo.
You can check that it is there with the command

```
docker image ls
```

Now to run it in a container and expose the service over port 8080 run

```
docker run -d -p 8080:8080 bitlymetrics
```

Docker should respond with the ```<containerID>``` if it has been successfully started.

Now the application should be useable.

When you are done running the app you can execute
```
docker container stop <containerID>
```
which will stop the container.  

##App Instructions
Running the application exposes the HTTP service on localhost over port 8080
by default. You can then consume the service with your desired HTTP client.
The service exposes only one RESTful endpoint:

```GET /groups/{groupGuid}/countries/averages```

Note that the endpoint URL takes the user's Group GUID corresponding to the group for which
you want the click metrics.

In order to be able to successfully consume the API you must pass the OAuth2
bearer token issued for your Bitly account in the HTTP request header.  The
key-value pair for this header should be as follows:

```Authorization: Bearer <my_access_token>```

The returned data is the average number of clicks, over a 30 day period, for all
Bitlinks in the group corresponding to the provided groupGuid, by country.

##Performance, Bottlenecks, Future Work, and Optimizations
Pagination naively depends on the default number of results returned by the
Bitly API (default according to the documentation is 50).  It makes a request
for each next link returned, and won't terminate until the next link is empty.

A future optimization for pagination would involve measuring the number of API calls made by
the paginating methods (in this case only GetBitlinksForGroup) and the
latency of those calls, as well as looking at the total number of records that
the API reports in the first response (ie. the "total" field).  Then using this data
we can construct an optimal combination of page size (ie. "size") and number of
pages that we need to request.

Another naïve implementation involves the retrieval of the total clicks per bitlink
by country.  This implementation won't scale to large numbers of bitlinks, since
there is an http request to GetBitlinkClicksByCountry for each bitlink.

The hope here is that a future optimized implementation would not have to rely on an HTTP API
that retrieves the metrics for only one Bitlink.  For example, a Bitly engineer
would probably have access to the backend processes that gather that data for the
external API and thus could engineer a more efficient process perhaps using
a more efficient protocol.  One way to optimize the current code would be to spin
up goroutines for each Bitlink's HTTP request so that the requests could run
concurrently, and then write the results via go channels to a shared data structure.
This solution isn't trivial since it would require synchronization and locking mechanisms
for concurrent writing to shared memory.

Currently the API doesn't support customization of the time range over which the
metrics are collected.  The API as is can be readily extended to receive query
parameters from the user request and propagate them through to the
GetBitlinkClicksByCountry request method.
