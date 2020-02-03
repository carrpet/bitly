Installation instructions:

You should first run the automated tests.  To do so, execute

`go test` in the root directory.  They should all pass.  


Performance and Bottlenecks:
Pagination naively depends on the default number of results returned by the
Bitly API (default according to the documentation is 50).  It makes a request
for each next link returned, and won't terminate until the next link is empty.
A future optimization would involve measuring the number of API calls made by
the paginating methods (in this case only GetBitlinksForGroup) and the
latency of those calls, as well as looking at the total number of records that
the API reports in the first response (ie. the "total" field).  Then using this data
we can construct an optimal combination of page size (ie. "size") and number of
pages that we need to request.

Another naive implementation involves the retrieval of the total clicks per bitlink
by country.  This implementation won't scale to large numbers of bitlinks, since
there is an http request to GetBitlinkClicksByCountry for each bitlink.  The hope
here is that a future optimized implementation would not have to rely on an HTTP API
that retrieves the metrics for only one Bitlink.  For example, a Bitly engineer
would probably have access to the backend processes that gather that data for the
external API and thus could engineer a more efficient process perhaps using
a more efficient protocol.  One way to optimize the current code would be to spin
up goroutines for each Bitlink's HTTP request so that the requests could run
concurrently, and then write the results via go channels to a shared data structure.
This solution isn't trivial since it would require synchronization, locking mechanisms
for concurrent writing to shared memory.

Future Work and Optimization:

Currently the API doesn't support customizing the time range over which the
metrics are collected.  The internals
