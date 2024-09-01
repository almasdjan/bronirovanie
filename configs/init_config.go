package configs

import (
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConfig() error {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)

	viper.SetConfigFile(filepath.Join(currentDir, "config.yml"))
	viper.AddConfigPath(currentDir)
	viper.AutomaticEnv()
	viper.SetDefault("DB_PASSWORD", "defaultpassword")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config file: %v", err)
	}

	return nil
}
