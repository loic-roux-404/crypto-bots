package helpers

import "fmt"

// SimpleMap type
type SimpleMap map[string](interface{})

var errorFmt = "No field named %s available in map"

// GetInMap secure get in map
func GetInMap(m SimpleMap, name string) (interface{}, error) {
	if v := m[name]; v != nil {
		return v, nil
	}

	return nil, fmt.Errorf(errorFmt, name)
}
