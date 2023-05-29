package config

const (
	configFileName = "config"
	configFileType = "toml"
)

type GeneralOptions struct {
 	Address string `json:"address"`
	CrtFile string `json:"crt_file"`
	KeyFile string `json:"key_file"`
	Workdir string `json:"default_workdir"`
	Debug   bool   `json:"debug"`
}

var General GeneralOptions
