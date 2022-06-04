package merge

import (
	"fmt"
	"strconv"
	"strings"
)

type Info struct {
	Errors []error
}

func (info *Info) mergeValue(path []string, patch map[string]any, key string, value any) any {
	patchValue, patchHasValue := patch[key]

	if !patchHasValue {
		return value
	}

	_, patchValueIsObject := patchValue.(map[string]any)

	path = append(path, key)
	pathStr := strings.Join(path, ".")

	if _, ok := value.(map[string]any); ok {
		if !patchValueIsObject {
			err := fmt.Errorf("patch value must be object for key \"%v\"", pathStr)
			info.Errors = append(info.Errors, err)
			return value
		}

		return info.mergeObjects(value, patchValue, path)
	}

	if _, ok := value.([]any); ok && patchValueIsObject {
		return info.mergeObjects(value, patchValue, path)
	}

	return patchValue
}

func (info *Info) mergeObjects(data, patch any, path []string) any {
	if patchObject, ok := patch.(map[string]any); ok {
		if dataArray, ok := data.([]any); ok {
			ret := make([]any, len(dataArray))

			for i, val := range dataArray {
				ret[i] = info.mergeValue(path, patchObject, strconv.Itoa(i), val)
			}

			return ret
		} else if dataObject, ok := data.(map[string]any); ok {
			ret := make(map[string]any)

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
