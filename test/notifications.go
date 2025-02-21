package notify

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/nodebytehosting/syscapture/plugins/notify"
	"github.com/nodebytehosting/syscapture/internal/handler"
	"github.com/nodebytehosting/syscapture/internal/config"
	"gopkg.in/yaml.v3"
)

var appConfig *config.Config

func setup(t *testing.T) {
	yamlData := `
notifications:
  enabled: true
  discord_webhook: "https://example.com/webhook"
`
	var configData config.Config
	if err := yaml.Unmarshal([]byte(yamlData), &configData); err != nil {
		t.Fatal(err)
	}
	appConfig = &configData
}

func TestNotifyPlugin_Init(t *testing.T) {
	setup(t)
	plugin := notify.NewNotifyPlugin()
	logger := handler.NewSysCaptureLogger()
	plugin.SetDiscordWebhook(appConfig.Notifications.DiscordWebhook)
	plugin.Init(logger)
	
	// Test initialization
	err := plugin.Init(logger)
	assert.NoError(t, err)
	plugin.GetLogger().Info("NotifyPlugin initialized successfully")
}

func TestNotifyPlugin_Start(t *testing.T) {
	setup(t)
	plugin := notify.NewNotifyPlugin()
	logger := handler.NewSysCaptureLogger()
	plugin.SetDiscordWebhook(appConfig.Notifications.DiscordWebhook)
	plugin.Init(logger)
	
	// Test starting the plugin
	err := plugin.Start()
	assert.NoError(t, err)
	plugin.GetLogger().Info("NotifyPlugin started successfully")
}

func TestNotifyPlugin_Stop(t *testing.T) {
	setup(t)
	plugin := notify.NewNotifyPlugin()
	plugin.SetDiscordWebhook(appConfig.Notifications.DiscordWebhook)
	logger := handler.NewSysCaptureLogger()
	plugin.Init(logger)
	
	// Test stopping the plugin
	err := plugin.Stop()
	assert.NoError(t, err)
	plugin.GetLogger().Info("NotifyPlugin stopped successfully")
}
