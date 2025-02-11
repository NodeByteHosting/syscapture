package plugins

import (
	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/plugin"
)

type SamplePlugin struct {
	logger handler.Logger
}

func NewPlugin() plugin.Plugin {
	return &SamplePlugin{}
}

func (p *SamplePlugin) Name() string {
	return "SamplePlugin"
}

func (p *SamplePlugin) Init(logger handler.Logger) error {
	p.logger = logger
	p.logger.Info("Initializing SamplePlugin")
	return nil
}

func (p *SamplePlugin) Start() error {
	p.logger.Info("Starting SamplePlugin")
	return nil
}

func (p *SamplePlugin) Stop() error {
	p.logger.Info("Stopping SamplePlugin")
	return nil
}
