package main

import (
	"fmt"
	"strings"
)

func iniParindom(kalimat string) string {
	tulisan := []rune{}
	tulisan2 := []rune{}
	kalimat = strings.ToLower(kalimat)

	for _, huruf := range kalimat {
		if string(huruf) != " " {
			tulisan = append(tulisan, huruf)
		}
	}

	for i := len(tulisan) - 1; i >= 0; i-- {
		tulisan2 = append(tulisan2, tulisan[i])
	}

	if string(tulisan2) == string(tulisan) {
		return "Ini adalah kalimat parindom"
	} else {
		return "Ini bukan kalimat parindom"
	}
}

func main() {
	kalimat := "Live on time, emit no evil"
	fmt.Println(iniParindom(kalimat))
}
