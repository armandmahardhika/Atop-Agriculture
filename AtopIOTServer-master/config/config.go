package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// VERSION Server's version string
var VERSION = "1.0.00"

func setDefault() {
	viper.SetDefault("expire", 600)
	viper.SetDefault("refreshExpire", 1200)
	// viper.SetDefault("email", map[string]string{
	// 	"host":     "gmail.com",
	// 	"user":     "test",
	// 	"password": "test"})
	// viper.SetDefault("publickey", "./public.key")
	// viper.SetDefault("expire", 600)
	// viper.SetDefault("refreshExpire", 1200)

}

func init() {
	viper.SetConfigName("settings")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Printf("Fatal error config file: %s use default config\n", err)
		setDefault()
		if err = viper.WriteConfigAs("./settings.json"); err != nil {
			log.Println("can not written setting file.")
		}

		log.Println("writting default config file")
	}
	viper.SetDefault("version", VERSION)
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed: ", e.Name)
		if err = viper.ReadInConfig(); err != nil {
			log.Println("viper read config error!!!")
		}
	})
}
