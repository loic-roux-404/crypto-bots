package helpers

import "fmt"

type fnMap map[string](func (args ...interface{}) (interface{}, error))

// GetInMap secure get in map
func GetInMap(m fnMap, name string, args ...interface{}) (interface{}, error) {
	if fn, ok := m[name]; ok {
        final, err := fn(args...)

        return final, err
    }

    return nil, fmt.Errorf("No field named %s available in map", name)
}
