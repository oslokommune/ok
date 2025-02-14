// Package config provides AWS configuration management functionality, primarily focused on
// AWS SSO (Single Sign-On) profile generation and authentication.
//
// The package handles:
//   - AWS SSO authentication flow
//   - Account and role discovery
//   - AWS config file generation
//   - Profile name customization through templates
//
// Example usage:
//
//	err := config.Generate(config.Options{
//	    SsoStartUrl: "https://my-sso.awsapps.com/start",
//	    SsoRegion:   "us-east-1",
//	    Region:      "eu-west-1",
//	})
package config
