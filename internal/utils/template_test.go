package utils_test

import (
	"testing"

	"github.com/argoproj-labs/gitops-promoter/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestRenderStringTemplate(t *testing.T) {
	tests := map[string]struct {
		template string
		data     any
		expected string
		wantErr  bool
	}{
		"can render template successfully": {
			template: "Name: {{ .Name }}",
			data: map[string]string{
				"Name": "John",
			},
			expected: "Name: John",
		},
		"can render using sprig functions": {
			template: "Name: {{ trunc 1 .Name }}",
			data: map[string]string{
				"Name": "John",
			},
			expected: "Name: J",
		},
		"cannot render using sensitive sprig functions": {
			template: "{{ env HOME }}",
			wantErr:  true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := utils.RenderStringTemplate(test.template, test.data)
			if test.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, test.expected, result)
		})
	}
}
