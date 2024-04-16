package main

import "fmt"

func hitungCharacter(char string) []map[string]interface{} {
	hurufDaftar := []map[string]interface{}{}

	for _, i := range char {
		hitungan := 0

		for _, j := range char {
			if j == i {
				hitungan++
			}
		}

		addTotal := map[string]interface{}{
			"char":  string(i),
			"total": hitungan,
		}
		hurufDaftar = append(hurufDaftar, addTotal)
	}

	hurufDuplicate := []map[string]interface{}{}
	for _, item := range hurufDaftar {
		found := false
		for _, dup := range hurufDuplicate {
			if dup["char"] == item["char"] {
				found = true
				break
			}
		}
		if !found {
			hurufDuplicate = append(hurufDuplicate, item)
		}
	}

	return hurufDuplicate
}

func main() {
	char := "hello"
	daftarHuruf := hitungCharacter(char)
	fmt.Println(daftarHuruf)
}
