package inspectago

import "testing"

var simpleValueCases = []struct {
	value         interface{}
	propertyCount int
}{
	{"", 0},
	{0, 0},
	{true, 0},
	{nil, 0},
}

// NamedStruct ... Sample type for autogeneration of type name.
type NamedStruct struct {
	Name, Email string
	Age         int
}

func Test_GetNames_WithNoProperties_ReturnsNoProperties(t *testing.T) {
	memcache = newMemoryCache()
	properties := PropertyList{}

	names := properties.GetNames()

	if len(names) != 0 {
		t.Errorf("Expected count of names to be 0 but got %d", len(names))
	}
}

func Test_GetNames_WithProperties_ReturnsProperties(t *testing.T) {
	memcache = newMemoryCache()
	properties := PropertyList{}
	properties["A"] = propertyType{"A", "", "0"}
	properties["B"] = propertyType{"B", "", "0"}

	names := properties.GetNames()

	if names[0] != "A" || names[1] != "B" {
		t.Errorf("Expected name to be A and B but got %q", names)
	}
}

func Test_GetNamesAsCSV_WithNoProperties_ReturnsEmptyString(t *testing.T) {
	memcache = newMemoryCache()
	properties := PropertyList{}

	names := properties.GetNamesAsCSV()

	if names != "" {
		t.Errorf("Expected names to be empty but got '%s'", names)
	}
}

func Test_GetNamesAsCSV_WithProperties_ReturnsPropertiesAsCSV(t *testing.T) {
	memcache = newMemoryCache()
	properties := PropertyList{}
	properties["A"] = propertyType{"A", "", "0"}
	properties["B"] = propertyType{"B", "", "0"}

	names := properties.GetNamesAsCSV()

	if names != "A,B" {
		t.Errorf("Expected names to be 'A,B' but got '%s'", names)
	}
}

func Test_GetNamesAsSQL_WithNoProperties_ReturnsEmptyString(t *testing.T) {
	memcache = newMemoryCache()
	properties := PropertyList{}

	names := properties.GetNamesAsSQL()

	if names != "" {
		t.Errorf("Expected names to be empty but got '%s'", names)
	}
}

func Test_GetNamesAsSQL_WithProperties_ReturnsPropertiesAsSQL(t *testing.T) {
	memcache = newMemoryCache()
	properties := PropertyList{}
	properties["A"] = propertyType{"A", "", "0"}
	properties["B"] = propertyType{"B", "", "0"}

	names := properties.GetNamesAsSQL()

	if names != "A, B" {
		t.Errorf("Expected names to be 'A, B' but got '%s'", names)
	}
}

func Test_Inspect_WithSimpleValueType_FindsNoProperties(t *testing.T) {
	for _, i := range simpleValueCases {
		memcache = newMemoryCache()

		insp, _ := Inspect(i.value)
		result := len(insp)

		if result != i.propertyCount {
			t.Errorf("Expected property count for %q to be %d but got %d", i.value, i.propertyCount, result)
		}
	}
}

func Test_Inspect_WithStruct_FindsOnlyPublicProperties(t *testing.T) {
	memcache = newMemoryCache()
	var thing = struct{ A, b, C, d int }{1, 2, 3, 4}

	insp, _ := Inspect(thing)
	result := insp

	if result.GetNamesAsCSV() != "A,C" {
		t.Errorf("Expected 'A,C' but got '%s'", result.GetNamesAsCSV())
	}
}

func Test_Inspect_WithStruct_ExtractsValues(t *testing.T) {
	memcache = newMemoryCache()
	var thing = struct{ A, b, C, d int }{1, 2, 3, 4}

	insp, _ := Inspect(thing)
	result := insp

	if result["A"].Value != "1" {
		t.Errorf("Expected 1 but got %q", result["A"].Value)
	}
	if result["C"].Value != "3" {
		t.Errorf("Expected 3 but got %q", result["A"].Value)
	}
}

func Test_Inspect_OnFirstRequest_DoesNotReturnFromCache(t *testing.T) {
	memcache = newMemoryCache()
	var thing = struct{ A int }{1}

	_, fromCache := Inspect(thing)

	if fromCache {
		t.Error("Expected fresh but got from cache.")
	}
}

func Test_Inspect_OnSubsequentRequest_ReturnsFromCache(t *testing.T) {
	memcache = newMemoryCache()
	var thing = NamedStruct{"Karl", "karl@younger.days", 30}

	_, _ = Inspect(thing)
	_, fromCache := Inspect(thing)

	if !fromCache {
		t.Error("Expected from cache but got afresh.")
	}
}

func Test_Inspect_WhenFetchedFromCache_DoesNotCacheValues(t *testing.T) {
	memcache = newMemoryCache()
	var thing1 = NamedStruct{"Karl", "karl@younger.days", 25}
	var thing2 = NamedStruct{"Karl", "karl@younger.days", 30}

	_, _ = Inspect(thing1)
	result, _ := Inspect(thing2)

	if result["Age"].Value == "25" {
		t.Error("Cached value was not updated.")
	}
}

func Test_Inspect_WithTypeStruct_DerivesTypeNameAutomatically(t *testing.T) {
	memcache = newMemoryCache()
	var thing = NamedStruct{"Karl", "karl@younger.days", 30}

	insp, _ := Inspect(thing)
	result := insp

	if result.GetNamesAsCSV() != "Age,Email,Name" {
		t.Errorf("Expected 'Age,Email,Name' but got '%s'", result.GetNamesAsCSV())
	}
}
