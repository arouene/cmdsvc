package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func init() {
	// default configuration
	viper.SetDefault("General.debug", false)
	viper.SetDefault("General.default_workdir", "")
	viper.SetDefault("General.address", ":3000")

	// additional path to search for a config file
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.cmdsvc/")
	viper.AddConfigPath("$HOME/.config/cmdsvc/")
	viper.AddConfigPath("/etc/cmdsvc/")

	readConfigFile(configFileName, configFileType)
	readConfigFile(authFileName, authFileType)

	viper.MergeInConfig()

	// get configuration from env vars
	viper.SetEnvPrefix("cmdsvc")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", " ", "_"))
	viper.AutomaticEnv()

	// general flags
	viper.UnmarshalKey("General", &General)

	// auth config
	viper.UnmarshalKey("Auth", &Auths)

	// services config
	viper.UnmarshalKey("Services", &Services)

	// Update service workdirs with default workdir if present
	for _, svc := range Services {
		if svc.WorkDir == "" && General.Workdir != "" {
			svc.WorkDir = General.Workdir
		}
	}
}

// readConfigFile open files and merge configurations
func readConfigFile(fileName, fileType string) {
	viper.SetConfigName(fileName)
	viper.SetConfigType(fileType)

	if err := viper.MergeInConfig(); err != nil {
		log.Printf("Config file %s cannot be read\n", fileName+"."+fileType)
		log.Printf("%s\n", err)
		os.Exit(1)
	} else {
		log.Printf("Config file used: %s\n", viper.ConfigFileUsed())
	}
}
