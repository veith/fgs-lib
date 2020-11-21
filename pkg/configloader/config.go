package configloader

import (
	"fmt"
	"github.com/spf13/viper"
)

/**
 * This will look for a .shortnamerc in your app folder
 * or a /etc/shortname/.shortnamerc
 * or a .shortname in the home directory
 */
func Load(shortname string) {
	viper.SetConfigName("." + shortname + "rc")
	viper.SetConfigType("yaml")                    // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/" + shortname + "/") // path to look for the config file in
	viper.AddConfigPath("$HOME/." + shortname)     // look in home too
	viper.AddConfigPath(".")                       // optionally look for config in the working directory
	err := viper.ReadInConfig()                    // Find and read the config file
	if err != nil {                                // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

}

/**
 * Set the defaults, you also can simply use `viper.SetDefault(key, value)`
 */
func SetDefault(key string, value string) {
	viper.SetDefault(key, value)
}
