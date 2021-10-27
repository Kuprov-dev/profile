package db

import "fmt"

func contains(s []string, str string) bool {
	for _, v := range s {
		fmt.Println(v, str)
		if v == str {
			return true
		}
	}

	return false
}
