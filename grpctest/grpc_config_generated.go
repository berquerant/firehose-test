// Code generated by "goconfig -type Executable string,Port int,HealthWait time.Duration -option -configOption Option -output grpc_config_generated.go"; DO NOT EDIT.

package grpctest

import "time"

type ConfigItem[T any] struct {
	modified     bool
	value        T
	defaultValue T
}

func (s *ConfigItem[T]) Set(value T) {
	s.modified = true
	s.value = value
}
func (s *ConfigItem[T]) Get() T {
	if s.modified {
		return s.value
	}
	return s.defaultValue
}
func (s *ConfigItem[T]) Default() T {
	return s.defaultValue
}
func (s *ConfigItem[T]) IsModified() bool {
	return s.modified
}
func NewConfigItem[T any](defaultValue T) *ConfigItem[T] {
	return &ConfigItem[T]{
		defaultValue: defaultValue,
	}
}

type Config struct {
	Executable *ConfigItem[string]
	Port       *ConfigItem[int]
	HealthWait *ConfigItem[time.Duration]
}
type ConfigBuilder struct {
	executable string
	port       int
	healthWait time.Duration
}

func (s *ConfigBuilder) Executable(v string) *ConfigBuilder {
	s.executable = v
	return s
}
func (s *ConfigBuilder) Port(v int) *ConfigBuilder {
	s.port = v
	return s
}
func (s *ConfigBuilder) HealthWait(v time.Duration) *ConfigBuilder {
	s.healthWait = v
	return s
}
func (s *ConfigBuilder) Build() *Config {
	return &Config{
		Executable: NewConfigItem(s.executable),
		Port:       NewConfigItem(s.port),
		HealthWait: NewConfigItem(s.healthWait),
	}
}

func NewConfigBuilder() *ConfigBuilder { return &ConfigBuilder{} }
func (s *Config) Apply(opt ...Option) {
	for _, x := range opt {
		x(s)
	}
}

type Option func(*Config)

func WithExecutable(v string) Option {
	return func(c *Config) {
		c.Executable.Set(v)
	}
}
func WithPort(v int) Option {
	return func(c *Config) {
		c.Port.Set(v)
	}
}
func WithHealthWait(v time.Duration) Option {
	return func(c *Config) {
		c.HealthWait.Set(v)
	}
}
