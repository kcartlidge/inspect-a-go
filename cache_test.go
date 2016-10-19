package inspectago

import "testing"

type implementation struct {
	Name     string
	Instance cache
}

var implementations = []implementation{
	{"MemoryCache", newMemoryCache()},
}

func TestCache_Fetch_WithNewThing_ReturnsNothing(t *testing.T) {
	for _, imp := range implementations {
		imp.Instance = imp.Instance.AsNew()

		_, ok := imp.Instance.Fetch("unknown")

		if ok {
			t.Errorf("Unexpected entry found in cache (%s)", imp.Name)
		}
	}
}

func TestCache_Set_WithThing_SetsThing(t *testing.T) {
	for _, imp := range implementations {
		imp.Instance = imp.Instance.AsNew()
		thing := struct{ Name string }{"thing"}

		imp.Instance.Set("test", thing)

		_, ok := imp.Instance.Fetch("test")
		if !ok {
			t.Errorf("Expected entry missing from cache (%s)", imp.Name)
		}
	}
}

func TestCache_Clear_RemovesThings(t *testing.T) {
	for _, imp := range implementations {
		thing := struct{ Name string }{"thing"}
		imp.Instance.Set("test", thing)

		imp.Instance = imp.Instance.AsNew()

		_, ok := imp.Instance.Fetch("test")
		if ok {
			t.Error("Cache should be empty, but is not.")
		}
	}
}
