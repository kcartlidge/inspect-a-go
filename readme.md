# inspect-a-go

*Automated inspection of a struct instance's properties and values, returning property lists, CSV or SQL.*

Version: 0.1.1, Licence: MIT

Pass in any 'thing' whose type is a ```struct``` and it will be inspected for public properties and their values. It will also provide both *CSV* and *SQL* lists of those properties or their values, quoted where necessary.

## Status

Alpha. In progress. Test coverage via TDD. See **Coming next** section below. *Signatures will change*.

## Ideal for ...

* Creating CSV exports for records stored as structs
* Automatically generating display columns on-screen
* Using as the basis of an auto-mapper
* Generation of SQL on the fly (parameterised)

## Limitations

The code maintains an in-memory cache of assessed properties. In order to do this you must provide the name of the type of thing assessed (to be used internally as a key). *This may soon switch to auto-discover the type name where the struct is not anonymous*.

Currently all public properties are located (private ones ignored) regardless of their types, however the returned property collection only contains the string representation of the actual values found not the underlying thing.

## Coming next

* Changes to visibility and return types of main ```Inspect``` function.
* Values of the underlying type not just their string version
* Multiple return values to include private fields too
* Return of SQL CSV for field names and quoted field values
* Understanding of attributes (e.g. for autogenerating form labels)
* Possible consideration of nested structs (undecided)

## Code serves as an example for

* Tests and code created via TDD
* Reflection to discover properties and their values
* Table-driven test cases
* Testing multiple implementations with one test via the interface

## Example usage:

The following code declares and populates an anonymous struct with sample data.
It then calls the ```Inspect``` function from *inspect-a-go* which returns a map of the public properties (only), which are printed out.

The ```password``` field is private, so there is no entry in the resulting property map so the sample below will result in no password being shown.

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
	props, _ := inspectago.Inspect("anon", anonymousStruct)

	fmt.Println("  Name:", props["Name"].Value)
	fmt.Println(" Email:", props["Email"].Value)
	fmt.Println("   Age:", props["Age"].Value)
	fmt.Println("  Pass:", props["password"].Value)
}
```

Output:
```
  Name: Karl
 Email: karl@younger.days
   Age: 30
  Pass:
```
