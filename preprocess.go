package vgui

import (
    "fmt"
    "strings"
)

// Preprocess processes base file imports for a Vgui file
// It returns the string with all of the #base import directives replaced with the contents of those files
// Or it will return an error if there is an issue reading source file
func Preprocess(provider FileSourceProvider, file string) (string, error) {
    input, err := provider.ReadAllLines(file)
    if err != nil {
        return "", err
    }
    sb := strings.Builder{}
    for _, line := range input {
        trimmed := strings.TrimSpace(line)
        if strings.HasPrefix(trimmed, "#base ") {
            baseFile := strings.TrimSpace(trimmed[len("#base "):])
            processed, err := Preprocess(provider, baseFile)
            if err != nil {
                return "", err
            }
            sb.WriteString(fmt.Sprintln(processed))
            continue
        }

        sb.WriteString(fmt.Sprintln(line))
    }
    return sb.String(), nil
}