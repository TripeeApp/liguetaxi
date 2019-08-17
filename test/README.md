### liguetax tests ###

This directory ctonains additional test suites beyonde unit tests already in [../liguetaxi](../liguetaxi). WHereas the unit tests run very quickly (since they don't make any network calls), the tests in this directory hit directly the [Ligue Taxi API]().

The test package are:

## integration ##

This will exercice the entire liguetaxi library against the live [Ligue Taxi API][]. These tests will verify that the library is properly coded against the actual behaviour of the API, and will (hopefully) fail upon any incompatible change in the API.

Because these tests are running using live data, there is a much higher propability of false positives in test failures due to network issues, test data having been changed, etc.

Additionally, in order to test the methods that modify data, a real API credential need to be present. While tests will try to be well-behaved in terms of what data they modify, it is **strongly** recommended that these tests only be run using a dedicated test account.

Run tests using:

    $ LIGUETAXI_HOST='<LT_HOST>' LIGUETAXI_TOKEN='<LT_TOKEN>' go test -v ./integration

Additionally there is a flag `log` that will make tests log the requests:

    $ LIGUETAXI_HOST='<LT_HOST>' LIGUETAXI_TOKEN='<LT_TOKEN>' go test -v -log ./integration

[Ligue Taxi API]: https://portal.taxidigital.net/suporte/php/API_TD/
