package main

func CountingSort(contacts []string) map[string]int {
	CountingSort := make(map[string]int)
	num := string
	for num = range contacts {
		CountingSort[num]++
	}

	return CountingSort
}
