package config

import (
	"fmt"
	"log"
	"ogomez/mkt-export/pkg/util"
	"os"

	"github.com/go-playground/validator"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ConfigBuilder struct {
	Config Config
}

func (c ConfigBuilder) Build(cfgFile string) (Config, error) {

	// Read config file with Viper
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cobra")
	}
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	// Validate Config
	if err := viper.Unmarshal(&c.Config); err != nil {
		log.Printf("unable to unmarshall the config %v", err)
		return c.Config, err
	}
	validate := validator.New()
	if err := validate.Struct(&c.Config); err != nil {
		log.Printf("Missing required attributes %v\n", err)
		return c.Config, err
	}
	//Build Output
	if _, err := c.buildOutput(); err != nil {
		log.Printf("Can't mount Output folder %v\n", err)
		return c.Config, err
	}

	return c.Config, nil
}

// Build output folder
func (c ConfigBuilder) buildOutput() (string, error) {
	return util.BuildPath(c.Config.Output)
}

