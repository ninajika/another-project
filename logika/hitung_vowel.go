package main

import (
	"fmt"
)

func hitungHuruf(huruf string) {
	hurufVowel := 0
	hurufConst := 0
	vowels := "aeiou"

	for _, char := range huruf {
		// Convert character to lowercase
		charLower := char | 32 // Convert to lowercase

		// Check if the character is a vowel
		isVowel := false
		for _, vowel := range vowels {
			if charLower == vowel {
				isVowel = true
				break
			}
		}

		// Increment the vowel or consonant count
		if isVowel {
			hurufVowel++
		} else {
			hurufConst++
		}
	}

	fmt.Println("Total huruf Konsonan:", hurufConst)
	fmt.Println("Total Huruf Vowel:", hurufVowel)
}

func main() {
	huruf := "hello"
	hitungHuruf(huruf)
}
