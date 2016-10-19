package inspectago

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

// PropertyList ... All known public properties on a thing.
type PropertyList map[string]PropertyType

// PropertyType ... A single (public) property on a thing.
type PropertyType struct {
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
	sort.Strings(keys)
	return keys
}

// GetNamesAsCSV ... Returns all the property names (sorted) as a CSV list.
func (p PropertyList) GetNamesAsCSV() string {
	return strings.Join(p.GetNames(), ",")
}

// Inspect ... Discover the properties for the given thing, with a from-cache flag.
func Inspect(name string, thing interface{}) (PropertyList, bool) {

	if thing == nil {
		return PropertyList{}, false
	}

	if reflect.TypeOf(thing).Kind().String() == "struct" {
		properties := PropertyList{}
		cached := false
		thingValue := reflect.ValueOf(thing)
		typeOfT := thingValue.Type()

		result, ok := memcache.Fetch("i__" + name)
		if ok {
			cached = true
			cachedProperties := result.(PropertyList)
			for _, p := range cachedProperties {
				field := thingValue.FieldByName(p.Name)
				p.Value = fmt.Sprintf("%v", field.Interface())
				properties[p.Name] = p
			}
		} else {
			cached = false
			for i := 0; i < thingValue.NumField(); i++ {
				field := thingValue.Field(i)
				name := fmt.Sprintf("%s", typeOfT.Field(i).Name)
				isPublic := name[0:1] == strings.ToUpper(name[0:1])
				if isPublic {
					fieldType := fmt.Sprintf("%s", field.Type())
					value := fmt.Sprintf("%v", field.Interface())
					propertyType := PropertyType{name, fieldType, value}
					properties[name] = propertyType
				}
			}
		}

		memcache.Set("i__"+name, properties)
		return properties, cached
	}

	return PropertyList{}, false
}
