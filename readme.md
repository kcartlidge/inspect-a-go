# inspect-a-go

*Automated inspection of a struct instance's properties and values, returning property lists, CSV or SQL.*

Version: 0.2.0, Licence: MIT

Pass in any 'thing' whose type is a ```struct``` and it will be inspected for public properties and their values. It will also provide both *CSV* and *SQL* lists of those properties or their values, quoted where necessary. Internally it will also cache struct analysis for structs that are not anonymous.

## Status

Alpha. In progress. Test coverage via TDD. See **Coming next** section below.

## Building the package and running tests

**This is not necessary if all you want to do is make use of the package.** See the section labelled *Example usage* below if you simply want to consume it.

**Building**
```
cd <wherever-you-cloned-it-to>
go build
go install
```

**Running tests**
```
cd <wherever-you-cloned-it-to>
go test
```

## Ideal for ...

* Creating CSV exports for records stored as structs
* Automatically generating display columns on-screen
* Using as the basis of an auto-mapper
* Generation of SQL on the fly (parameterised)

## Limitations

Currently all public properties are located (private ones ignored) regardless of their types, however the returned property collection only contains the string representation of the actual values found not the underlying thing.

## Coming next

* Values of the underlying type not just their string version
* Multiple return values to include private fields too
* Return of SQL and CSV field values (quoted where relevant)
* Understanding of attributes (e.g. for autogenerating form labels)
* Possible consideration of nested structs (undecided)

## Code serves as an example for

* Tests and code created via TDD
* Reflection to discover properties and their values
* Table-driven test cases
* Testing multiple implementations with one test via the interface

## Example usage:

The following code declares and populates an anonymous struct with sample data. If the struct had been declared as a named type first then internally it's discovered properties would be cached for future hits.

It then calls the ```Inspect``` function from *inspect-a-go* which returns a map of the (public only) properties. These are then displayed, followed by sample SQL.

The ```password``` field is private, so there is no entry in the resulting property map and the sample code will show nothing.

``` go
package main

import (
	"fmt"

	"github.com/kcartlidge/inspect-a-go"
)

var anonymousStruct = struct {
	Name, Email string
	Age         int
	password    string
}{"Karl", "karl@younger.days", 30, "secret"}

func main() {
	props, _ := inspectago.Inspect(anonymousStruct)

	fmt.Println("      Name:", props["Name"].Value)
	fmt.Println("     Email:", props["Email"].Value)
	fmt.Println("       Age:", props["Age"].Value)
	fmt.Println("  Password:", props["password"].Value)

	fmt.Println()
	fmt.Println("       SQL: select", props.GetNamesAsSQL(), "from accounts")
}
```

Output:
```
      Name: Karl
     Email: karl@younger.days
       Age: 30
      Pass:

       SQL: select Age,Email,Name from accounts
```
