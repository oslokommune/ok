package add

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateStackName(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		outputFolder string
		expected     string
	}{
		{ ///////////////////////////////////////////////////////////
			name: "just StackName",
			input: `StackName: "whatever"

Serverless:
  Enable: true
`,
			outputFolder: "my-db",
			expected: `StackName: "my-db"

Serverless:
  Enable: true
`,
		}, ///////////////////////////////////////////////////////////
		{ ///////////////////////////////////////////////////////////
			name: "StackName with data",
			input: `StackName: "some-app"
app-data.StackName: "some-app-data"

AppName: "hello-world"
`,
			outputFolder: "app-foo",
			expected: `StackName: "app-foo"
app-data.StackName: "app-foo-data"

AppName: "hello-world"
`,
		}, ///////////////////////////////////////////////////////////
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := updateStackName(tt.input, tt.outputFolder)
			require.Equal(t, tt.expected, result)
		})
	}
}
