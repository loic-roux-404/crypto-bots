package helpers

import "fmt"

type fnMap map[string](func () (interface{}, error))

// GetInMap secure get in map
func GetInMap(m fnMap, name string) (interface{}, error) {
	if fn, ok := m[name]; ok {
        final, err := fn()

        return final, err
    }

    return nil, fmt.Errorf("No field named %s available in map", name)
}
