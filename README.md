# smoothbrain-test
A very simple library that helps with assertions in unit tests.

### Use `go get` With Your Project

```
go get github.com/barbell-math/smoothbrain-test
```

Then import and use the library as desired in your code.

This package has [zero dependencies](./go.mod).

### Documentation

The `Test.go` file is heavily annotated. Please refer to that file.

### Example

This package allows for writing unit tests that look like this:

```go
import "github.com/barbell-math/smoothbrain-test"

func TestCombine(t *testing.T) {
	h1 := 69
	h2 := 420
	h3 := 5280

	smoothbraintest.Eq(t, h1, h1)
	smoothbraintest.Neq(t, h2, h1)
	smoothbraintest.Neq(t, h3, h1)
	smoothbraintest.Neq(t, h1, h2)
	smoothbraintest.Neq(t, h3, h2)
	smoothbraintest.Neq(t, h1, h3)
	smoothbraintest.Neq(t, h2, h3)
}
```
