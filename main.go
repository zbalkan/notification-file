package main

import (
	"context"
	"fmt"
	"os"

	"github.com/crowdsecurity/crowdsec/pkg/protobufs"
	"github.com/hashicorp/go-hclog"
	plugin "github.com/hashicorp/go-plugin"
	"gopkg.in/yaml.v2"
)

var logger hclog.Logger = hclog.New(&hclog.LoggerOptions{
	Name:       "file-plugin",
	Level:      hclog.LevelFromString("DEBUG"),
	Output:     os.Stderr,
	JSONFormat: true,
})

type PluginConfig struct {
	Name     string  `yaml:"name"`
	LogLevel *string `yaml:"log_level"`

	LogPath string `yaml:"log_path"`
}

type FilePlugin struct {
	ConfigByName map[string]PluginConfig
}

func (n *FilePlugin) Configure(ctx context.Context, config *protobufs.Config) (*protobufs.Empty, error) {
	d := PluginConfig{}
	if err := yaml.Unmarshal(config.Config, &d); err != nil {
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

	f, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return nil, err
	}
	// defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("[%s] %s\n", notification.Name, notification.Text)); err != nil {
		return nil, err
	}

	logger.Info(fmt.Sprintf("Appended new alert: [%s] %s\n", notification.Name, notification.Text))

	return nil, nil
}

func main() {
	var handshake = plugin.HandshakeConfig{
		ProtocolVersion:  1,
		MagicCookieKey:   "CROWDSEC_PLUGIN_KEY",
		MagicCookieValue: os.Getenv("CROWDSEC_PLUGIN_KEY"),
	}

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
