package merge

import (
	"fmt"
	"strconv"
	"strings"
)

type Info struct {
	Errors []error
}

func (info *Info) mergeValue(path []string, patch map[string]interface{}, key string, value interface{}) interface{} {
	patchValue, patchHasValue := patch[key]

	if !patchHasValue {
		return value
	}

	_, patchValueIsObject := patchValue.(map[string]interface{})

	path = append(path, key)
	pathStr := strings.Join(path, ".")

	if _, ok := value.(map[string]interface{}); ok {
		if !patchValueIsObject {
			err := fmt.Errorf("patch value must be object for key \"%v\"", pathStr)
			info.Errors = append(info.Errors, err)
			return value
		}

		return info.mergeObjects(value, patchValue, path)
	}

	if _, ok := value.([]interface{}); ok && patchValueIsObject {
		return info.mergeObjects(value, patchValue, path)
	}

	return patchValue
}

func (info *Info) mergeObjects(data, patch interface{}, path []string) interface{} {
	if patchObject, ok := patch.(map[string]interface{}); ok {
		if dataArray, ok := data.([]interface{}); ok {
			ret := make([]interface{}, len(dataArray))

			for i, val := range dataArray {
				ret[i] = info.mergeValue(path, patchObject, strconv.Itoa(i), val)
			}

			return ret
		} else if dataObject, ok := data.(map[string]interface{}); ok {
			ret := make(map[string]interface{})

			for k, v := range dataObject {
				ret[k] = info.mergeValue(path, patchObject, k, v)
			}

			return ret
		}
	}

	return data
}

func Merge(data, patch map[string]any) (map[string]any, error) {
	info := &Info{}
	ret := info.mergeObjects(data, patch, nil)
	if len(info.Errors) > 0 {
		return nil, info.Errors[0]
	}

	return ret.(map[string]any), nil
}
