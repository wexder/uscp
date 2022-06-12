package merge

func iteratePatch(data map[string]any, patch map[string]any, path []string) {
	for key, value := range patch {
		if innerPatch, ok := value.(map[string]any); ok {
			if _, ok := data[key].(map[string]any); !ok {
				data[key] = map[string]any{}
			}
			iteratePatch(data[key].(map[string]any), innerPatch, append(path, key))
		} else {
			data[key] = value
		}
	}
}

func Merge(data map[string]any, patch map[string]any) (map[string]any, error) {
	iteratePatch(data, patch, []string{})
	return data, nil
}
