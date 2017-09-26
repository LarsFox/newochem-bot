package core

func stringInArray(s string, arr []string) int {
	for i, item := range arr {
		if s == item {
			return i
		}
	}
	return -1
}

func intInArray(s int, arr []int) int {
	for i, item := range arr {
		if s == item {
			return i
		}
	}
	return -1
}

func chunksString(array []string, length int) [][]string {
	var result [][]string
	var i = 0
	for ; i < len(array)-length; i += length {
		result = append(result, array[i:i+length])
	}
	result = append(result, array[i:])
	return result
}
