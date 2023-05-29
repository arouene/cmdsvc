package config

const (
	authFileName   = "auth"
	authFileType   = "toml"
)

type Auth struct {
	User   string
	Token  string
	Groups []string
}

type AuthsList []Auth

var Auths AuthsList
