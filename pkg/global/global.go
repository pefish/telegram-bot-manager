package global

import "github.com/pefish/go-commander"

type Config struct {
	commander.BasicConfig
	Token          string `json:"token" default:"" usage:"Bot token."`
	CommandsJsFile string `json:"commands_js_file" default:"" usage:"Commands js file."`
}

var GlobalConfig Config
