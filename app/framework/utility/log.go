package utility

import "fmt"

func LogSuccess(text any) {
	fmt.Printf("\033[32m %s \033[0m\n", text)
}

func LogError(text any) {
	fmt.Printf("\033[31m %s \033[0m\n", text)
}

func LogWarning(text any) {
	fmt.Printf("\033[33m %s \033[0m\n", text)
}

func LogNeutral(text any) {
	fmt.Printf("\033[35m %s \033[0m\n", text)
}


func LogBlue(text any) {
	fmt.Printf("\033[34m %s \033[0m\n", text)
}
