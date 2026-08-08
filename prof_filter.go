package main

import (
	"fmt"
	"strings"
)

func profane_replacement(body string) string {

	profane_words := []string{"kerfuffle", "sharbert", "fornax"}

	body_word_list := strings.Split(body, " ")
	check_logs := []string{}

	for _, word := range profane_words {
		for i, body_word := range body_word_list {
			check_logs = append(check_logs, fmt.Sprintf("Checking %s against %s", word, body_word))
			if strings.ToLower(body_word) == word {
				body_word_list[i] = "****"
			}
		}
	}

	cleaned_body := strings.Join(body_word_list, " ")
	return cleaned_body
}
