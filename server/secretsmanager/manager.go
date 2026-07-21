package secretsmanager

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// ProviderConfig maps directly to each provider block in the JSON
type ProviderConfig struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scope        []string `json:"scope"`
	RedirectURI  string   `json:"redirect_uri"`
}

// SecretsManager handles loading and retrieving secrets safely
type SecretsManager struct {
	mu        sync.RWMutex
	providers map[string]ProviderConfig
}

// NewSecretsManager initializes and loads the manager from a JSON path
func NewSecretsManager(filePath string) (*SecretsManager, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets file: %w", err)
	}

	var data map[string]ProviderConfig
	if err := json.Unmarshal(fileBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to parse secrets JSON: %w", err)
	}

	return &SecretsManager{
		providers: data,
	}, nil
}

// Get returns the config for a specific provider and a boolean indicating if it exists
func (sm *SecretsManager) Get(provider string) (ProviderConfig, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	config, exists := sm.providers[provider]
	return config, exists
}
