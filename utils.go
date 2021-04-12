package main

import "strconv"

func CastInt(i interface{}) int {
	s, ok := i.(string)
	if !ok {
		return -1
	}

	if i64, err := strconv.Atoi(s); err == nil {
		return i64
	}

	return -1
}
