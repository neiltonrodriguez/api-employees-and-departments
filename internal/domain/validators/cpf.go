package validators

import (
	"regexp"
	"strconv"
)

// ValidateCPF validates a Brazilian CPF number
func ValidateCPF(cpf string) bool {
	// Remove non-numeric characters
	re := regexp.MustCompile(`[^0-9]`)
	cpf = re.ReplaceAllString(cpf, "")

	// CPF must have 11 digits
	if len(cpf) != 11 {
		return false
	}

	// Known invalid CPFs (all digits the same)
	invalidCPFs := []string{
		"00000000000", "11111111111", "22222222222", "33333333333",
		"44444444444", "55555555555", "66666666666", "77777777777",
		"88888888888", "99999999999",
	}

	for _, invalid := range invalidCPFs {
		if cpf == invalid {
			return false
		}
	}

	// Validate first check digit
	sum := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (10 - i)
	}
	remainder := sum % 11
	firstCheckDigit := 0
	if remainder >= 2 {
		firstCheckDigit = 11 - remainder
	}

	digit9, _ := strconv.Atoi(string(cpf[9]))
	if digit9 != firstCheckDigit {
		return false
	}

	// Validate second check digit
	sum = 0
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (11 - i)
	}
	remainder = sum % 11
	secondCheckDigit := 0
	if remainder >= 2 {
		secondCheckDigit = 11 - remainder
	}

	digit10, _ := strconv.Atoi(string(cpf[10]))
	if digit10 != secondCheckDigit {
		return false
	}

	return true
}
