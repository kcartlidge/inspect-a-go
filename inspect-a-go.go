package inspectago

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// PropertyList ... All known public properties on a thing.
type PropertyList map[string]propertyType

// propertyType ... A single (public) property on a thing.
type propertyType struct {
	Name  string
	Type  string
	Value string
}

// An in-memory cache of property lists.
var memcache = newMemoryCache()

// GetNames ... Returns all the property names (sorted).
func (p PropertyList) GetNames() []string {
	keys := []string{}
	for _, k := range p {
		keys = append(keys, k.Name)
	}

	// Without this, the ordering is (intentionally) pseudorandom and inconsistent.
	sort.Strings(keys)
	return keys
}

// GetNamesAsCSV ... Returns property names (sorted) as a CSV list.
func (p PropertyList) GetNamesAsCSV() string {
	return strings.Join(p.GetNames(), ",")
}

// GetNamesAsSQL ... Returns property names (sorted) as a CSV list, spaced for clarity.
func (p PropertyList) GetNamesAsSQL() string {
	return strings.Join(p.GetNames(), ", ")
}

// Inspect ... Discover the properties for the given thing, with a from-cache flag.
func Inspect(thing interface{}) (PropertyList, bool) {

	if thing == nil {
		return PropertyList{}, false
	}

	// Only structs are supported.
	if reflect.TypeOf(thing).Kind().String() == "struct" {
		properties := PropertyList{}
		thingValue := reflect.ValueOf(thing)
		typeOfT := thingValue.Type()
		pkgPath := typeOfT.PkgPath()

		// Derive a cache key from the name of a non-anonymous
		// structure and try the cache.
		cacheKey := ""
		cachedValue := *new(interface{})
		cached := false
		if pkgPath != "" {
			cacheKey = typeOfT.Name()
			cachedValue, cached = memcache.Fetch("i__" + cacheKey)
		}

		if cached {
			// Reuse the property definitions but fetch current values.
			cachedProperties := cachedValue.(PropertyList)
			for _, p := range cachedProperties {
				field := thingValue.FieldByName(p.Name)
				p.Value = fmt.Sprintf("%v", field.Interface())
				properties[p.Name] = p
			}
		} else {
			// Derive all from scratch then cache.
			cached = false
			for i := 0; i < thingValue.NumField(); i++ {
				field := thingValue.Field(i)
				name := fmt.Sprintf("%s", typeOfT.Field(i).Name)
				isPublic := name[0:1] == strings.ToUpper(name[0:1])
				if isPublic {
					fieldType := fmt.Sprintf("%s", field.Type())
					value := fmt.Sprintf("%v", field.Interface())
					propType := propertyType{name, fieldType, value}
					properties[name] = propType
				}
			}
			if cacheKey != "" {
				memcache.Set("i__"+cacheKey, properties)
			}
		}
		return properties, cached
	}

	// Default for non-structs.
	return PropertyList{}, false
}
