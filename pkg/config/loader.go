package config

import (
	"context"
	"log/slog"
	"os"
	"slices"
	"sync/atomic"
	"time"

	"gopkg.in/yaml.v3"
)

const RefreshInterval = 5 * time.Second

type DynamicConfig[T any] interface {
	Config() *T
	RegisterConfigUpdateHook(hook func(c *T))
	Start(ctx context.Context) error
}

// dynamicConfigImpl represents the structure of our configuration
type dynamicConfigImpl[T any] struct {
	config            atomic.Pointer[T]
	rawConfig         []byte
	filepath          string
	configUpdateHooks []func(c *T)
}

func (c *dynamicConfigImpl[T]) Config() *T {
	return c.config.Load()
}

func (c *dynamicConfigImpl[T]) RegisterConfigUpdateHook(hook func(c *T)) {
	c.configUpdateHooks = append(c.configUpdateHooks, hook)
}

func (c *dynamicConfigImpl[T]) Start(ctx context.Context) error {
	err := c.load()
	if err != nil {
		return err
	}

	return c.watch(ctx)
}

func NewDynamicConfig[T any](filepath string) DynamicConfig[T] {
	c := &dynamicConfigImpl[T]{
		filepath: filepath,
	}
	return c
}

// loadConfig loads the configuration from the given filename
func (c *dynamicConfigImpl[T]) load() error {
	bytes, err := os.ReadFile(c.filepath)
	if err != nil {
		return err
	}

	// If the configuration hasn't changed, don't reload it
	if slices.Equal(c.rawConfig, bytes) {
		return nil
	}

	var config T
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		return err
	}

	slog.Info("Config loaded", "config", string(bytes))
	c.rawConfig = bytes
	c.config.Store(&config)
	return nil
}

// watchConfig listens for changes to the configuration file and reloads it
func (c *dynamicConfigImpl[T]) watch(ctx context.Context) error {
	ticker := time.NewTicker(RefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := c.load()
			if err != nil {
				slog.Error("Failed to reload config", "error", err)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

type StaticConfig[T any] struct {
	cfg *T
}

func NewStaticConfig[T any](cfg *T) DynamicConfig[T] {
	return &StaticConfig[T]{cfg: cfg}
}

func NewStaticConfigByJsonOrDie[T any](cfg []byte) DynamicConfig[T] {
	var c T
	err := yaml.Unmarshal(cfg, &c)
	if err != nil {
		panic(err)
	}
	return &StaticConfig[T]{cfg: &c}
}

func (s *StaticConfig[T]) Config() *T {
	return s.cfg
}

func (s *StaticConfig[T]) RegisterConfigUpdateHook(_ func(c *T)) {
}

func (s *StaticConfig[T]) Start(ctx context.Context) error {
	return nil
}
