package configs

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config is the main configuration structure
type Config struct {
	Version  string          `json:"version"`
	Gateway  GatewayConfig   `json:"gateway"`
	LLM      LLMConfig       `json:"llm"`
	Channels ChannelsConfig  `json:"channels"`
	Agents   []AgentConfig   `json:"agents,omitempty"`
}

// GatewayConfig for the WebSocket server
type GatewayConfig struct {
	Port     int    `json:"port"`
	Host     string `json:"host,omitempty"`
	Token    string `json:"token,omitempty"`
	TLS      bool   `json:"tls,omitempty"`
	CertFile string `json:"certFile,omitempty"`
	KeyFile  string `json:"keyFile,omitempty"`
}

// LLMConfig for language model settings
type LLMConfig struct {
	Provider     string  `json:"provider"`
	Model        string  `json:"model"`
	APIKey       string  `json:"apiKey,omitempty"`
	MaxTokens    int     `json:"maxTokens"`
	Temperature  float64 `json:"temperature"`
	SystemPrompt string  `json:"systemPrompt"`
}

// ChannelsConfig for messaging channels
type ChannelsConfig struct {
	Telegram TelegramConfig `json:"telegram,omitempty"`
	Discord  DiscordConfig  `json:"discord,omitempty"`
}

// TelegramConfig for Telegram bot
type TelegramConfig struct {
	Enabled    bool    `json:"enabled"`
	Token      string  `json:"token,omitempty"`
	AllowedIDs []int64 `json:"allowedIds,omitempty"`
	AdminIDs   []int64 `json:"adminIds,omitempty"`
}

// DiscordConfig for Discord bot
type DiscordConfig struct {
	Enabled      bool     `json:"enabled"`
	Token        string   `json:"token,omitempty"`
	GuildID      string   `json:"guildId,omitempty"`
	AllowedRoles []string `json:"allowedRoles,omitempty"`
	Prefix       string   `json:"prefix,omitempty"`
}

// AgentConfig for multi-agent support
type AgentConfig struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Workspace    string `json:"workspace,omitempty"`
	Model        string `json:"model,omitempty"`
	SystemPrompt string `json:"systemPrompt,omitempty"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Version: "1",
		Gateway: GatewayConfig{
			Port: 8080,
			Host: "0.0.0.0",
		},
		LLM: LLMConfig{
			Provider:     "anthropic",
			Model:        "claude-sonnet-4-20250514",
			MaxTokens:    4096,
			Temperature:  0.7,
			SystemPrompt: "You are a helpful AI assistant powered by HiveClaw. Be concise and helpful.",
		},
		Channels: ChannelsConfig{
			Telegram: TelegramConfig{Enabled: false},
			Discord:  DiscordConfig{Enabled: false, Prefix: "!"},
		},
		Agents: []AgentConfig{
			{
				ID:   "main",
				Name: "Main Agent",
			},
		},
	}
}

// Load loads configuration from a file
func Load(path string) (*Config, error) {
	if path == "" {
		home, _ := os.UserHomeDir()
		candidates := []string{
			filepath.Join(home, ".hiveclaw", "config.json"),
			"./config.json",
			"./hiveclaw.json",
		}

		for _, c := range candidates {
			if _, err := os.Stat(c); err == nil {
				path = c
				break
			}
		}
	}

	if path == "" {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// Save saves configuration to a file
func (c *Config) Save(path string) error {
	if path == "" {
		home, _ := os.UserHomeDir()
		configDir := filepath.Join(home, ".hiveclaw")
		os.MkdirAll(configDir, 0755)
		path = filepath.Join(configDir, "config.json")
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// GetAPIKey returns the API key, checking env vars as fallback
func (c *Config) GetAPIKey() string {
	if c.LLM.APIKey != "" {
		return c.LLM.APIKey
	}

	switch c.LLM.Provider {
	case "anthropic":
		return os.Getenv("ANTHROPIC_API_KEY")
	case "openrouter":
		return os.Getenv("OPENROUTER_API_KEY")
	default:
		return ""
	}
}

// GetConfigPath returns the default config path
func GetConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".hiveclaw", "config.json")
}
