package config

import (
	"testing"
)

func TestSanitizeName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple case",
			input:    "HelloWorld",
			expected: "hello-world",
		},
		{
			name:     "multiple words with acronym",
			input:    "MyAWSRole",
			expected: "my-aws-role",
		},
		{
			name:     "special characters",
			input:    "*Hello@_World_foo_",
			expected: "hello-world-foo",
		},
		{
			name:     "multiple acronyms",
			input:    "AWSIAMRole",
			expected: "awsiam-role",
		},
		{
			name:     "consecutive capitals",
			input:    "MyABCRole",
			expected: "my-abc-role",
		},
		{
			name:     "number followed by word",
			input:    "hello2World",
			expected: "hello2-world",
		},
		{
			name:     "preserve existing hyphens",
			input:    "example-dev",
			expected: "example-dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitizeName(tt.input)
			if got != tt.expected {
				t.Errorf("sanitizeName(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGenerateConfig(t *testing.T) {
	tests := []struct {
		name     string
		profiles []Profile
		options  Options
		want     string
	}{
		{
			name: "single profile",
			profiles: []Profile{
				{
					Name:        "dev-admin",
					AccountID:   "123456789012",
					RoleName:    "AdminRole",
					AccountName: "dev",
				},
			},
			options: Options{
				SessionName: "mysession",
				SsoStartUrl: "https://my-start-url.awsapps.com/start",
				SsoRegion:   "us-east-1",
				Region:      "eu-west-1",
			},
			want: `; run 'ok aws generate' to update the configuration
[sso-session mysession]
sso_start_url = https://my-start-url.awsapps.com/start
sso_region = us-east-1
sso_registration_scopes = sso:account:access

[profile dev-admin]
sso_session = mysession
sso_account_id = 123456789012
sso_role_name = AdminRole
region = eu-west-1
`,
		},
		{
			name: "multiple profiles",
			profiles: []Profile{
				{
					Name:        "dev-admin",
					AccountID:   "123456789012",
					RoleName:    "AdminRole",
					AccountName: "dev",
				},
				{
					Name:        "prod-read",
					AccountID:   "987654321098",
					RoleName:    "ReadOnlyRole",
					AccountName: "prod",
				},
			},
			options: Options{
				SessionName: "mysession",
				SsoStartUrl: "https://my-start-url.awsapps.com/start",
				SsoRegion:   "us-east-1",
			},
			want: `; run 'ok aws generate' to update the configuration
[sso-session mysession]
sso_start_url = https://my-start-url.awsapps.com/start
sso_region = us-east-1
sso_registration_scopes = sso:account:access

[profile dev-admin]
sso_session = mysession
sso_account_id = 123456789012
sso_role_name = AdminRole
region = us-east-1

[profile prod-read]
sso_session = mysession
sso_account_id = 987654321098
sso_role_name = ReadOnlyRole
region = us-east-1
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a generator with test options
			generator, err := NewConfigGenerator(tt.options)
			if err != nil {
				t.Fatalf("NewConfigGenerator() error = %v", err)
			}

			got, err := generator.generateConfig(tt.profiles)
			if err != nil {
				t.Fatalf("generateConfig() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("generateConfig() output mismatch:\ngot:\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestGenerateProfileName(t *testing.T) {
	tests := []struct {
		name         string
		data         profileNameData
		template     string
		expected     string
		expectError  bool
		errorMessage string
	}{
		{
			name: "default template",
			data: profileNameData{
				SessionName: "MySession",
				AccountName: "my-account",
				AccountID:   "123456789012",
				RoleName:    "AdminRole",
			},
			template: "",
			expected: "my-session-my-account-admin-role",
		},
		{
			name: "custom template with role mapping",
			data: profileNameData{
				SessionName: "mysso",
				AccountName: "my-account",
				AccountID:   "123456789012",
				RoleName:    "AdminRole",
			},
			template: "{{.AccountName}}-{{if eq .RoleName \"AdminRole\"}}MyAdminRole{{else}}{{.RoleName}}{{end}}",
			expected: "my-account-my-admin-role",
		},
		{
			name: "template with account ID",
			data: profileNameData{
				SessionName: "mysso",
				AccountName: "MyAccount",
				AccountID:   "123456789012",
				RoleName:    "ReadRole",
			},
			template: "{{.AccountID}}-{{.RoleName}}",
			expected: "123456789012-read-role",
		},
		{
			name: "invalid template",
			data: profileNameData{
				SessionName: "mysso",
				AccountName: "MyAccount",
				AccountID:   "123456789012",
				RoleName:    "ReadRole",
			},
			template:    "{{.InvalidField}}",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateProfileName(tt.data, tt.template)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if got != tt.expected {
				t.Errorf("generateProfileName() = %q, want %q", got, tt.expected)
			}
		})
	}
}
