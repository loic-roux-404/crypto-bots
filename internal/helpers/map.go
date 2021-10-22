package helpers

import "fmt"

// FnMap type
type FnMap map[string](interface{})

// GetInMap secure get in map
func GetInMap(m FnMap, name string) (interface{}, error) {
	if v := m[name]; v != nil {
        println(v)
        return v, nil
    }

    return nil, fmt.Errorf("No field named %s available in map", name)
}
