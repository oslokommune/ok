package config

import (
	"context"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
)

const (
	defaultProfileNameTemplate = "{{.SessionName}}-{{.AccountName}}-{{.RoleName}}"
	awsConfigTemplate          = `; run 'ok aws generate' to update the configuration
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
	ssoClient := sso.New(sso.Options{
		Region: opts.SsoRegion,
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

	// Print summary header
	fmt.Printf("\nGenerating AWS config for %d profiles\n", len(profiles))
	fmt.Println("\nAWS CLI config (e.g., `$HOME/.aws/config`):")
	fmt.Println("---")

	config, err := g.generateConfig(profiles)
	if err != nil {
		return fmt.Errorf("generating config: %w", err)
	}

	fmt.Print(config)
	fmt.Println("---")

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
func sanitizeName(s string) string {
	if s == "" {
		return s
	}

	var result strings.Builder
	runes := []rune(s)

	for i := 0; i < len(runes); i++ {
		curr := runes[i]

		// Skip invalid characters
		if !unicode.IsLetter(curr) && !unicode.IsNumber(curr) && curr != '-' {
			if result.Len() > 0 && result.String()[result.Len()-1] != '-' {
				result.WriteRune('-')
			}
			continue
		}

		// Add hyphen in specific cases
		if i > 0 {
			prev := runes[i-1]
			isCurrentUpper := unicode.IsUpper(curr)
			isPrevLower := unicode.IsLower(prev)
			isPrevNumber := unicode.IsNumber(prev)

			needsHyphen := false

			if isCurrentUpper && isPrevLower {
				needsHyphen = true
			} else if isCurrentUpper && i < len(runes)-1 && unicode.IsLower(runes[i+1]) && !isPrevLower {
				needsHyphen = true
			} else if isPrevNumber && unicode.IsLetter(curr) {
				needsHyphen = true
			}

			if needsHyphen && result.Len() > 0 && result.String()[result.Len()-1] != '-' {
				result.WriteRune('-')
			}
		}

		result.WriteRune(unicode.ToLower(curr))
	}

	return strings.Trim(result.String(), "-")
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
