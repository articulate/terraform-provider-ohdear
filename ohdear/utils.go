package ohdear

func getKeysAsSlice(aMap map[string]interface{}) []string {
	keys := make([]string, 0, len(aMap))
	for k := range aMap {
		keys = append(keys, k)
	}

	return keys
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}

func convertInterfaceToStringArr(purportedList interface{}) []string {
	var arr []string
	rawArr, ok := purportedList.([]interface{})

	if ok {
		arr = make([]string, len(rawArr))
		for i, thing := range rawArr {
			arr[i] = thing.(string)
		}
	}

	return arr
}

// Converts interface to string array, if there are no elements it returns nil to conform with optional properties.
func convertInterfaceToStringArrNullable(purportedList interface{}) []string {
	arr := convertInterfaceToStringArr(purportedList)

	if len(arr) < 1 {
		return nil
	}

	return arr
}
