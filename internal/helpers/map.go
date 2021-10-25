package helpers

import "fmt"

// Map type
type Map map[string](interface{})

var errorFmt = "No field named %s available in map"

// GetInMap secure get in map
func GetInMap(m Map, name string) (interface{}, error) {
	if v := m[name]; v != nil {
        return v, nil
    }

    return nil, fmt.Errorf(errorFmt, name)
}
