package config

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
)

const (
	defaultProfileNameTemplate = "{{.SessionName}}-{{.AccountName}}-{{.RoleName}}"
	awsConfigTemplate          = `; Config starts here. Run 'ok aws generate' to generate new configuration.
[sso-session {{.SessionName}}]
sso_start_url = {{.StartURL}}
sso_region = {{.SSORegion}}
sso_registration_scopes = sso:account:access
{{- range .Profiles}}

[profile {{.Name}}]
sso_session = {{$.SessionName}}
sso_account_id = {{.AccountID}}
sso_role_name = {{.RoleName}}
region = {{$.ProfileRegion}}
{{- end}}
`
)

// NewConfigGenerator creates a new ConfigGenerator with the given options.
// It validates required fields and sets default values where appropriate.
func NewConfigGenerator(opts Options) (*ConfigGenerator, error) {
	if opts.SsoStartUrl == "" {
		return nil, fmt.Errorf("SSO start URL is required")
	}
	if opts.SsoRegion == "" {
		return nil, fmt.Errorf("SSO region is required")
	}

	if opts.Region == "" {
		opts.Region = opts.SsoRegion
	}

	if opts.SessionName == "" {
		opts.SessionName = strings.Split(strings.TrimPrefix(opts.SsoStartUrl, "https://"), ".")[0]
	}

	ssooidcClient := ssooidc.New(ssooidc.Options{
		Region: opts.SsoRegion,
	})
	//
	ssoClient := sso.New(sso.Options{
		Region: opts.SsoRegion,
		// We configure retries to mitigate TooManyRequestsException when fetching accounts and roles
		RetryMaxAttempts: 100,
		RetryMode:        aws.RetryModeStandard,
	})

	return &ConfigGenerator{
		ssooidcClient: ssooidcClient,
		ssoClient:     ssoClient,
		options:       opts,
	}, nil
}

func (g *ConfigGenerator) Generate(ctx context.Context) error {
	accounts, err := GetAccounts(ctx, g.ssoClient, g.ssooidcClient, g.options.SsoStartUrl, g.options.SsoRegion)
	if err != nil {
		return fmt.Errorf("getting accounts: %w", err)
	}
	if len(accounts) == 0 {
		return fmt.Errorf("no accounts found")
	}
	profiles, err := g.generateProfiles(accounts)
	if err != nil {
		return fmt.Errorf("generating profiles: %w", err)
	}

	fmt.Fprintf(os.Stderr, "\nGenerating AWS CLI configuration (i.e., $HOME/.aws/config) for %d profiles\n\n", len(profiles))

	config, err := g.generateConfig(profiles)
	if err != nil {
		return fmt.Errorf("generating config: %w", err)
	}

	fmt.Print(config)

	return nil
}

func (g *ConfigGenerator) generateProfiles(accounts []Account) ([]Profile, error) {
	var profiles []Profile
	for _, account := range accounts {
		for _, roleName := range account.Roles {
			profileName, err := generateProfileName(profileNameData{
				SessionName: g.options.SessionName,
				AccountName: account.Name,
				AccountID:   account.ID,
				RoleName:    roleName,
			}, g.options.Template)
			if err != nil {
				return nil, fmt.Errorf("generating profile name for account %s role %s: %w",
					account.Name, roleName, err)
			}

			profiles = append(profiles, Profile{
				Name:        profileName,
				AccountID:   account.ID,
				RoleName:    roleName,
				AccountName: account.Name,
			})
		}
	}
	return profiles, nil
}

func (g *ConfigGenerator) generateConfig(profiles []Profile) (string, error) {
	data := profileConfigData{
		SessionName:   g.options.SessionName,
		StartURL:      g.options.SsoStartUrl,
		SSORegion:     g.options.SsoRegion,
		ProfileRegion: g.options.Region,
		Profiles:      profiles,
	}

	tmpl, err := template.New("config").Parse(awsConfigTemplate)
	if err != nil {
		return "", fmt.Errorf("parsing config template: %w", err)
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("executing config template: %w", err)
	}

	return buf.String(), nil
}

// sanitizeName converts a string into an opinionated profile name by applying
// consistent formatting rules and removing invalid characters.
func sanitizeName(input string) string {
	sep := "-"
	// Split on acronym boundaries
	reg := regexp.MustCompile(`([A-Z]+)([A-Z][a-z])`)
	result := reg.ReplaceAllString(input, "${1}"+sep+"${2}")

	// Split on word boundaries
	reg = regexp.MustCompile(`([a-z0-9])([A-Z])`)
	result = reg.ReplaceAllString(result, "${1}"+sep+"${2}")

	// Convert to lowercase
	result = strings.ToLower(result)

	// Replace non-alphanumeric chars with separator
	reg = regexp.MustCompile(`[^a-z0-9]+`)
	result = reg.ReplaceAllString(result, sep)

	// Trim hyphens from start/end
	return strings.Trim(result, sep)
}

func generateProfileName(data profileNameData, templateStr string) (string, error) {
	if templateStr == "" {
		templateStr = defaultProfileNameTemplate
	}

	tmpl, err := template.New("profile-name").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("parsing profile name template: %w", err)
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("executing profile name template: %w", err)
	}

	return sanitizeName(result.String()), nil
}

// Generate creates AWS SSO profile configurations based on the provided options.
// It handles the SSO authentication flow and generates AWS config file content.
func Generate(opts Options) error {
	generator, err := NewConfigGenerator(opts)
	if err != nil {
		return fmt.Errorf("creating config generator: %w", err)
	}
	return generator.Generate(context.Background())
}
