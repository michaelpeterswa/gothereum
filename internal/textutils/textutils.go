package textutils

import "fmt"

func InsertSymbolAndColor(text string, status string) string {
	switch status {
	case "more":
		return fmt.Sprintf("\033[32m\ufa34 %s\033[0m", text)
	case "equal":
		return fmt.Sprintf("\033[37m\ufa33 %s\033[0m", text)
	case "less":
		return fmt.Sprintf("\033[31m\ufa32 %s\033[0m", text)
	default:
		return text
	}
}
