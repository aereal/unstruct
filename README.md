[![status][ci-status-badge]][ci-status]
[![PkgGoDev][pkg-go-dev-badge]][pkg-go-dev]

# unstruct

**unstruct** is a library to read external values and decode them into the structs.

unstruct reads values from external sources using **`Source`** interface.

It is very simple interface, so you can easily implement custom `Source`.

## Synopsis

```go
import (
	"fmt"
	"os"

	"github.com/aereal/unstruct"
)

func main() {
	decoder := unstruct.NewDecoder(unstruct.NewEnvironmentSource())
	var val struct {
		ExampleString string
		ExampleInt    int
		ExampleBool   bool
	}
	_ = decoder.Decode(&val)
}
```

See examples on [pkg.go.dev][pkg-go-dev].

## Installation

```sh
go get github.com/aereal/unstruct
```

## Motivation

The practical application may have various and complicated configuration.

If you build the configuration from external data sources, you may face issues like below:

- _fetch from remote data sources_
  - the application must read secret data that like credentials from remote data sources using secure channels.
- _initialization_
  - complicated configuration may have many fields; you must initialize all of them.
- _validation_:
  - you have to check the existence of all of the mandatory fields.
- _type conversions_
  - the application configuration may include some values that differs from strings, such as numbers, booleans, or anything else.
  - external data sources may support just limited scalar types, so you have to convert/parse external values into the application configuration's values.

unstruct works as simple decoder that decodes external data into Go's values like encoding/json.

## License

See LICENSE file.

[pkg-go-dev]: https://pkg.go.dev/github.com/aereal/unstruct
[pkg-go-dev-badge]: https://pkg.go.dev/badge/aereal/unstruct
[ci-status-badge]: https://github.com/aereal/unstruct/workflows/CI/badge.svg?branch=main
[ci-status]: https://github.com/aereal/unstruct/actions/workflows/CI
[12-factor app]: https://12factor.net/config
