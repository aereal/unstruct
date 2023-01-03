package unstruct_test

import (
	"fmt"
	"os"

	"github.com/aereal/unstruct"
)

func Example() {
	decoder := unstruct.NewDecoder(unstruct.NewEnvironmentSource())
	var val struct {
		ExampleString string
		ExampleInt    int
		ExampleBool   bool
	}
	os.Setenv("EXAMPLE_STRING", "str")
	os.Setenv("EXAMPLE_INT", "42")
	os.Setenv("EXAMPLE_BOOL", "true")
	if err := decoder.Decode(&val); err != nil {
		panic(err)
	}
	fmt.Printf("%#v", val)
	// Output:
	// struct { ExampleString string; ExampleInt int; ExampleBool bool }{ExampleString:"str", ExampleInt:42, ExampleBool:true}
}
