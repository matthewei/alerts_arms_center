package configfileparser

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

type ConfigFileParser struct {
	vyaml *viper.Viper
}

func New() (*ConfigFileParser, error) {
	//获取项目的执行路径
	path, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(path, "config")
	vyaml := viper.New()
	vyaml.AddConfigPath(configPath)
	vyaml.SetConfigName("config") // 配置文件名
	vyaml.SetConfigType("yaml")
	if err := vyaml.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
		}
	}
	return &ConfigFileParser{vyaml}, nil
}

func (c ConfigFileParser) ReadSection(k string, v interface{}) error {
	var sections = make(map[string]interface{})
	err := c.vyaml.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	if _, ok := sections[k]; !ok {
		sections[k] = v
	}
	return nil
}

func (c ConfigFileParser) ReadAllSection(v interface{}) error {
	err := c.vyaml.Unmarshal(v)
	if err != nil {
		return err
	}
	return nil
}
