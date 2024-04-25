package tool_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gitlab.com/zigal0/architect/internal/cli/tool"
)

func Test_ToCamelCaseWithFirstUpper(t *testing.T) {
	t.Parallel()

	testData := map[string]string{
		"":                           "",
		"m":                          "M",
		"my":                         "My",
		" my_project\n":              "MyProject",
		" my_project\t":              "MyProject",
		" my_project\r":              "MyProject",
		"my   awesome  project ":     "MyAwesomeProject",
		"_my____awesome__project_":   "MyAwesomeProject",
		"--my-awesome---project-":    "MyAwesomeProject",
		"- \n-my-awesome---project-": "MyAwesomeProject",
	}

	for input, expected := range testData {
		output := tool.ToCamelCaseWithFirstUpper(input)

		require.Equal(t, expected, output)
	}
}

func Test_ToCamelCaseWithFirstLower(t *testing.T) {
	t.Parallel()

	testData := map[string]string{
		"":                           "",
		"m":                          "m",
		"my":                         "my",
		" my_project\n":              "myProject",
		" my_project\t":              "myProject",
		" my_project\r":              "myProject",
		"my   awesome  project ":     "myAwesomeProject",
		"_my____awesome__project_":   "myAwesomeProject",
		"--my-awesome---project-":    "myAwesomeProject",
		"- \t-my-awesome---project-": "myAwesomeProject",
	}

	for input, expected := range testData {
		output := tool.ToCamelCaseWithFirstLower(input)

		require.Equal(t, expected, output)
	}
}
