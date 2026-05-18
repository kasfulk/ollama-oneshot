package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	OllamaHost             string
	OllamaModel            string
	DefaultTool            string
	PromptEnhancement      bool
	PromptEnhancementModel string
	YoloMode               bool
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	viper.SetDefault("OLLAMA_HOST", "127.0.0.1:11434")
	viper.SetDefault("DEFAULT_MODEL", "kimi-k2.6:cloud")
	viper.SetDefault("DEFAULT_TOOL", "claude")
	viper.SetDefault("PROMPT_ENHANCEMENT", true)
	viper.SetDefault("PROMPT_ENHANCEMENT_MODEL", "deepseek-v4-flash")
	viper.SetDefault("YOLO_MODE", false)

	viper.AutomaticEnv()
	viper.BindEnv("OLLAMA_HOST")
	viper.BindEnv("DEFAULT_MODEL")
	viper.BindEnv("DEFAULT_TOOL")
	viper.BindEnv("PROMPT_ENHANCEMENT")
	viper.BindEnv("PROMPT_ENHANCEMENT_MODEL")
	viper.BindEnv("YOLO_MODE")

	cfg := &Config{
		OllamaHost:             viper.GetString("OLLAMA_HOST"),
		OllamaModel:            viper.GetString("DEFAULT_MODEL"),
		DefaultTool:            viper.GetString("DEFAULT_TOOL"),
		PromptEnhancement:      viper.GetBool("PROMPT_ENHANCEMENT"),
		PromptEnhancementModel: viper.GetString("PROMPT_ENHANCEMENT_MODEL"),
		YoloMode:               viper.GetBool("YOLO_MODE"),
	}

	if cfg.OllamaHost == "" {
		return nil, fmt.Errorf("OLLAMA_HOST is required")
	}

	return cfg, nil
}

func (c *Config) ApplyFlags(model, tool string, noEnhance, yoloMode bool) {
	if model != "" {
		c.OllamaModel = model
	}
	if tool != "" {
		c.DefaultTool = tool
	}
	if noEnhance {
		c.PromptEnhancement = false
	}
	if yoloMode {
		c.YoloMode = true
	}
}

func (c *Config) OllamaURL() string {
	host := c.OllamaHost
	return fmt.Sprintf("http://%s/api/generate", host)
}

func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return "/tmp"
}
