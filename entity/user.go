package entity

type User struct {
	ID              string
	Authentications []*Authentication
}

type Authentication struct {
	Provider Provider
	SocialID string
}

type Provider string

const (
	ProviderGoogle  Provider = "Google"
	ProviderTwitter Provider = "Twitter"
	ProviderApple   Provider = "Apple"
)
