package config

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type funcDownloader func() ([]byte, error)

func (fn funcDownloader) DownloadFile(ctx context.Context, filePath string) ([]byte, error) {
	return fn()
}

func TestDownloadBoilerplateConfig(t *testing.T) {
	type testcase struct {
		input  string
		expect BoilerplateConfig
	}

	tests := map[string]testcase{
		"empty file": {
			input:  ``,
			expect: BoilerplateConfig{},
		},
		"variables": {
			input: `
variables:
  - name: TemplateVersion
    type: string
    description: Internal tracking of template version - do NOT edit
    default: 2.2.0
  - name: UntypedVariable
  - name: StringVariable
    type: string
  - name: UntypedWithDefault
    default: clearly a string
  - name: DataMap
    type: map
    default:
      TrueValue:  true
      FalseValue: false
      IntValue:   42
      FloatValue: 3.14
      StringValue: "hello world"
      AnotherMap:
        IsThisComplex: true
        WillPeopleConfigureThigsThisWay: false
`,
			expect: BoilerplateConfig{
				Variables: []BoilerplateVariable{
					{
						Name:        "TemplateVersion",
						Type:        "string",
						Description: "Internal tracking of template version - do NOT edit",
						Default:     "2.2.0",
					},
					{
						Name: "UntypedVariable",
					},
					{
						Name: "StringVariable",
						Type: "string",
					},
					{
						Name:    "UntypedWithDefault",
						Default: "clearly a string",
					},
					{
						Name: "DataMap",
						Type: "map",
						Default: map[string]any{
							"TrueValue":   true,
							"FalseValue":  false,
							"IntValue":    42,
							"FloatValue":  3.14,
							"StringValue": "hello world",
							"AnotherMap": map[string]any{
								"IsThisComplex":                   true,
								"WillPeopleConfigureThigsThisWay": false,
							},
						},
					},
				},
			},
		},
		"dependencies": {
			input: `
dependencies:
  - name: versions
    template-url: ../versions
    output-folder: .
  - name: networking-data
    template-url: ../networking-data
    output-folder: ../networking-data
`,
			expect: BoilerplateConfig{
				Dependencies: []BoilerplateDependency{
					{
						Name:         "versions",
						TemplateUrl:  "../versions",
						OutputFolder: ".",
					},
					{
						Name:         "networking-data",
						TemplateUrl:  "../networking-data",
						OutputFolder: "../networking-data",
					},
				},
			},
		},
		"full config": {
			input: `
variables:
  - name: StackName
    type: string
    description: The name of the stack
  - name: Environment
    type: string
    description: The environment to deploy to
dependencies:
  - name: versions
    template-url: ../versions
    output-folder: .
  - name: networking-data
    template-url: ../networking-data
    output-folder: ../networking-data
`,

			expect: BoilerplateConfig{
				Variables: []BoilerplateVariable{
					{
						Name:        "StackName",
						Type:        "string",
						Description: "The name of the stack",
					},
					{
						Name:        "Environment",
						Type:        "string",
						Description: "The environment to deploy to",
					},
				},
				Dependencies: []BoilerplateDependency{
					{
						Name:         "versions",
						TemplateUrl:  "../versions",
						OutputFolder: ".",
					},
					{
						Name:         "networking-data",
						TemplateUrl:  "../networking-data",
						OutputFolder: "../networking-data",
					},
				},
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx := context.Background()
			downloader := funcDownloader(func() ([]byte, error) {
				return []byte(tc.input), nil
			})
			config, err := DownloadBoilerplateConfig(ctx, downloader, "boilerplate.yml")
			require.NoError(t, err, "expect boil config to parse nicely")
			require.Equal(t, tc.expect.Dependencies, config.Dependencies, "expect config dependencies to match")
			require.Equal(t, tc.expect.Variables, config.Variables, "expect config variables to match")
			//require.Equal(t, tc.expect, *config, "expect config to match")
		})
	}
}
