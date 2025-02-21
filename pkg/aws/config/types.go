package config

import (
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/aws/aws-sdk-go-v2/service/ssooidc"
)

// Options configures the AWS SSO config generation
type Options struct {
	SsoStartUrl string
	SsoRegion   string
	Region      string
	SessionName string
	Template    string
}

// Account represents an AWS SSO account
type Account struct {
	ID    string
	Name  string
	Roles []string
}

// Profile represents an AWS SSO profile configuration
type Profile struct {
	Name        string
	AccountID   string
	RoleName    string
	AccountName string
}

// profileNameData contains data for profile name generation
type profileNameData struct {
	SessionName string
	AccountName string
	AccountID   string
	RoleName    string
}

// profileConfigData contains data for profile config generation
type profileConfigData struct {
	SessionName   string
	StartURL      string
	SSORegion     string
	ProfileRegion string
	Profiles      []Profile
}

// ConfigGenerator handles AWS SSO configuration generation
type ConfigGenerator struct {
	ssooidcClient *ssooidc.Client
	ssoClient     *sso.Client
	options       Options
}
