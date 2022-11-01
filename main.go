package main

import (
	"context"
	"fmt"
	"os"

	"github.com/crowdsecurity/crowdsec/pkg/protobufs"
	"github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

var logger hclog.Logger = hclog.New(&hclog.LoggerOptions{
	Name:       "file-plugin",
	Level:      hclog.LevelFromString("DEBUG"),
	Output:     os.Stderr,
	JSONFormat: true,
})

type PluginConfig struct {
	Name         string     `yaml:"name"`
	LogLevel     *string    `yaml:"log_level"`
	LogPath      string     `yaml:"log_path"`
	LogFormatter string     `yaml:"log_formatter"`
	Rotate       FileRotate `yaml:"rotate"`
	LogFormat    LogFormat  `yaml:"log_format"`
}

type FileRotate struct {
	Enabled  bool  `yaml:"enabled"`
	MaxSize  int   `yaml:"max_size"`
	MaxFiles int   `yaml:"max_files"`
	MaxAge   int   `yaml:"max_age"`
	Compress *bool `yaml:"compress"`
}

type LogFormat struct {
	CustomFormat     string `yaml:"custom_format"`
	FormatterName    string `yaml:"formatter_name"`
	CustomTimeFormat string `yaml:"custom_time_format"`
}

type FilePlugin struct {
	ConfigByName map[string]PluginConfig
}

func (n *FilePlugin) Configure(ctx context.Context, config *protobufs.Config) (*protobufs.Empty, error) {
	logger.Info("Configured called")
	d := PluginConfig{}
	if err := yaml.Unmarshal(config.Config, &d); err != nil {
		return nil, err
	}
	if err := d.SetDefaultLoggerConfig(); err != nil {
		logger.Info(fmt.Sprintf("Error happened %s", err.Error()))
		return nil, err
	}
	n.ConfigByName[d.Name] = d
	return &protobufs.Empty{}, nil
}

func (n *FilePlugin) Notify(ctx context.Context, notification *protobufs.Notification) (*protobufs.Empty, error) {
	if _, ok := n.ConfigByName[notification.Name]; !ok {
		return nil, fmt.Errorf("invalid plugin config name %s", notification.Name)
	}
	cfg := n.ConfigByName[notification.Name]
	if cfg.LogLevel != nil && *cfg.LogLevel != "" {
		logger.SetLevel(hclog.LevelFromString(*cfg.LogLevel))
	} else {
		logger.SetLevel(hclog.Info)
	}
	log.Info(notification.Text)

	logger.Info(fmt.Sprintf("Appended new alert: [%s] %s\n", notification.Name, notification.Text))

	return nil, nil
}

func main() {
	var handshake = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "CROWDSEC_PLUGIN_KEY",
		MagicCookieValue: os.Getenv("CROWDSEC_PLUGIN_KEY"),
	}
	logger.Info("Plugin called")
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshake,
		Plugins: map[string]plugin.Plugin{
			"file": &protobufs.NotifierPlugin{
				Impl: &FilePlugin{ConfigByName: make(map[string]PluginConfig)},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
		Logger:     logger,
	})
}

func (n *PluginConfig) SetDefaultLoggerConfig() error {

	/*Configure logs*/
	_maxsize := 500
	if n.Rotate.MaxSize != 0 {
		_maxsize = n.Rotate.MaxSize
	}
	_maxfiles := 3
	if n.Rotate.MaxFiles != 0 {
		_maxfiles = n.Rotate.MaxFiles
	}
	_maxage := 28
	if n.Rotate.MaxAge != 0 {
		_maxage = n.Rotate.MaxAge
	}
	_compress := true
	if n.Rotate.Compress != nil {
		_compress = *n.Rotate.Compress
	}
	// check if file exists
	_, err := os.Stat(n.LogPath)
	// create file if not exists, purposefully ignore errors
	if os.IsNotExist(err) {
		file, _ := os.OpenFile(n.LogPath, os.O_RDWR|os.O_CREATE, 0644)
		file.Close()
	}

	LogOutput := &lumberjack.Logger{
		Filename: n.LogPath,
	}
	if n.Rotate.Enabled {
		LogOutput.MaxSize = _maxsize
		LogOutput.MaxSize = _maxfiles
		LogOutput.MaxSize = _maxage
		LogOutput.Compress = _compress
	}
	log.SetOutput(LogOutput)
	log.SetLevel(log.InfoLevel)
	if n.LogFormat.CustomFormat != "" {
		log.SetFormatter(&easy.Formatter{
			TimestampFormat: "2006-01-02 15:04:05",
			LogFormat:       n.LogFormat.CustomFormat,
		})
	} else {
		logFormatter := &log.TextFormatter{TimestampFormat: n.LogFormat.CustomTimeFormat, FullTimestamp: true}
		log.SetFormatter(logFormatter)
	}
	return nil
}
