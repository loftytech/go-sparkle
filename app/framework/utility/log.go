package utility

import "fmt"

func LogSuccess(text string) {
	fmt.Println("\033[32m" + text + "\033[0m")
}

func LogError(text string) {
	fmt.Println("\033[31m" + text + " \033[0m\n")
}

func LogWarning(text string) {
	fmt.Println("\033[33m" + text + " \033[0m\n")
}

func LogNeutral(text string) {
	fmt.Println("\033[35m" + text + " \033[0m\n")
}
