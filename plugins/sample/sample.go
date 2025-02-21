package sample

import (
	"flag"

	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/plugin"
)

type SamplePlugin struct {
	logger  handler.Logger
	enabled bool
}

func NewPlugin() plugin.Plugin {
	return &SamplePlugin{
		enabled: *flag.Bool("sample-plugin-enabled", true, "Enable/disable the sample plugin"),
	}
}

func (p *SamplePlugin) Name() string {
	return "SamplePlugin"
}

func (p *SamplePlugin) Init(logger handler.Logger) error {
	p.logger = logger
	p.logger.Info("Initializing %s", "SamplePlugin")
	return nil
}

func (p *SamplePlugin) Start() error {
	p.logger.Info("%s is now running", "SamplePlugin")
	return nil
}

func (p *SamplePlugin) Stop() error {
	p.logger.Info("%s has been stopped", "SamplePlugin")
	return nil
}

func (p *SamplePlugin) IsEnabled() bool {
	return p.enabled
}

func (p *SamplePlugin) Register(manager *plugin.PluginManager) {
	manager.Register(p)
}
