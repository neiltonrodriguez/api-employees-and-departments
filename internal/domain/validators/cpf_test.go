package validators

import "testing"

func TestValidateCPF(t *testing.T) {
	tests := []struct {
		name     string
		cpf      string
		expected bool
	}{
		{
			name:     "valid CPF",
			cpf:      "12345678909",
			expected: true,
		},
		{
			name:     "valid CPF with formatting",
			cpf:      "123.456.789-09",
			expected: true,
		},
		{
			name:     "another valid CPF",
			cpf:      "11144477735",
			expected: true,
		},
		{
			name:     "invalid CPF - all zeros",
			cpf:      "00000000000",
			expected: false,
		},
		{
			name:     "invalid CPF - all ones",
			cpf:      "11111111111",
			expected: false,
		},
		{
			name:     "invalid CPF - all twos",
			cpf:      "22222222222",
			expected: false,
		},
		{
			name:     "invalid CPF - wrong check digit",
			cpf:      "12345678900",
			expected: false,
		},
		{
			name:     "invalid CPF - too short",
			cpf:      "123456789",
			expected: false,
		},
		{
			name:     "invalid CPF - too long",
			cpf:      "123456789012",
			expected: false,
		},
		{
			name:     "empty CPF",
			cpf:      "",
			expected: false,
		},
		{
			name:     "CPF with letters",
			cpf:      "123abc78909",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateCPF(tt.cpf)
			if result != tt.expected {
				t.Errorf("ValidateCPF(%s) = %v, expected %v", tt.cpf, result, tt.expected)
			}
		})
	}
}
