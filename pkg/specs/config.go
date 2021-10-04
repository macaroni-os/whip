/*

Copyright (C) 2021  Daniele Rondina <geaaru@sabayonlinux.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/
package specs

import (
	v "github.com/spf13/viper"

	"gopkg.in/yaml.v2"
)

const (
	WHIP_CONFIGNAME = "whip.yml"
	WHIP_ENV_PREFIX = "WHIP"
)

type Config struct {
	Viper *v.Viper `yaml:"-" json:"-"`

	General CGeneral `mapstructure:"general" json:"general,omitempty" yaml:"general,omitempty"`
	Logging CLogging `mapstructure:"logging" json:"logging,omitempty" yaml:"logging,omitempty"`

	Loader   string   `mapstructure:"loader,omitempty" json:"loader,omitempty" yaml:"loader,omitempty"`
	SpecDirs []string `mapstructure:"specs_dirs,omitempty" json:"specs_dirs,omitempty" yaml:"specs_dirs,omitempty"`
}

type CGeneral struct {
	Debug bool `mapstructure:"debug,omitempty" json:"debug,omitempty" yaml:"debug,omitempty"`
}

type CLogging struct {
	// Path of the logfile
	Path string `mapstructure:"path,omitempty" json:"path,omitempty" yaml:"path,omitempty"`
	// Enable/Disable logging to file
	EnableLogFile bool `mapstructure:"enable_logfile,omitempty" json:"enable_logfile,omitempty" yaml:"enable_logfile,omitempty"`
	// Enable JSON format logging in file
	JsonFormat bool `mapstructure:"json_format,omitempty" json:"json_format,omitempty" yaml:"json_format,omitempty"`

	// Log level
	Level string `mapstructure:"level,omitempty" json:"level,omitempty" yaml:"level,omitempty"`

	// Enable emoji
	EnableEmoji bool `mapstructure:"enable_emoji,omitempty" json:"enable_emoji,omitempty" yaml:"enable_emoji,omitempty"`
	// Enable/Disable color in logging
	Color bool `mapstructure:"color,omitempty" json:"color,omitempty" yaml:"color,omitempty"`
}

func NewConfig(viper *v.Viper) *Config {
	if viper == nil {
		viper = v.New()
	}

	GenDefault(viper)
	return &Config{Viper: viper}
}

func (c *Config) GetGeneral() *CGeneral {
	return &c.General
}

func (c *Config) GetLogging() *CLogging {
	return &c.Logging
}

func (c *Config) Unmarshal() error {
	c.Viper.ReadInConfig()

	err := c.Viper.Unmarshal(&c)

	return err
}

func (c *Config) Yaml() ([]byte, error) {
	return yaml.Marshal(c)
}

func (c *Config) AddSpecsDirs(dirs []string) {
	if len(dirs) > 0 {
		c.SpecDirs = append(c.SpecDirs, dirs...)
	}
}

func GenDefault(viper *v.Viper) {
	viper.SetDefault("general.debug", false)
	viper.SetDefault("general.p2pmode", false)

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.enable_logfile", false)
	viper.SetDefault("logging.path", "./logs/whip.log")
	viper.SetDefault("logging.json_format", false)
	viper.SetDefault("logging.enable_emoji", true)
	viper.SetDefault("logging.color", true)

	viper.SetDefault("loader", "dir")
	viper.SetDefault("specs_dirs", []string{})
}

func (g *CGeneral) HasDebug() bool {
	return g.Debug
}
