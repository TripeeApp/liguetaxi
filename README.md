# liguetaxi #

liguetaxi is a Go client library for acessing the [Ligue Taxi API][]

## Usage ##

```go
import "bitbicket.org/mobilitee/liguetaxi"
```

Construct a new LigueTaxi client, then use the various services on the client to
access different parts of the LigueTaxi API. For example:

```go
host, _ := url.Parse("https://portal.taxidigital.net/suporte/php/API_TD/api/")

ligtaxi := liguetaxi.New(host, "token", nil)

// Read user info.
user, _ := ligtaxi.User.Read(context.Background(), "00115422321", "João da Silva")

```

The services of a client divide the API into logical chunks and correspond to
the struct of the [Ligue Taxi API][] documentation.

For more sample code snippets, head over to the `_test.go` files.

### Authentication ###

The liguetaxi library handles the Authorization header with a custom Transport.

### Creating and Updating Resources ###

Resources that can be created or updated in the [Ligue Taxi API][] are exposed
by the liguetaxi library through structs.

```go
newUser := &liguetaxi.User{
        Name: "João da Silva",
        Email: "test@gmail.com",
        Phone: "11986548744",
        Password: "1234321",
        Classifier1: "Classifier Field",
}
ligtaxi.User.Create(context.Background(), newUser)
```

## Roadmap ##

This library is being initally developed for an internal application at
Mobilitee, so API methods will likely be implemented in the order they are
needed by the application.

## Versioning ##

liguetax folows [semver](https://semver.org/) as closely as we
can for tagging releases of the package. Because liguetaxi is a client
library for the [Ligue Taxi API][], which itself can change behaviour,
we've adopted the following policy:

* Increment the **major version** with any incompatible change to the
API functionality.
* Increment the **minor version** with any backward-compatible changes to the
API functionality.
* Increment the **patch version** with any backwards-compatible bug fixes.

### TODO ###
- Implement the ride methods
- Implement XML requests

[Ligue Taxi API]: https://portal.taxidigital.net/suporte/php/API_TD/
