package config

import (
	"time"

	"github.com/prometheus/alertmanager/template"
)

type Config struct {
	Interval time.Duration `yaml:"interval"`
	Notify   *Notify       `yaml:"notify"`
	Evaluate *Evaluate     `yaml:"evaluate"`
	Env      string        `yaml:"env"`
}

type Notify struct {
	DingTalk *DingTalk `yaml:"dingtalk"`
}

type DingTalk struct {
	Url     string `yaml:"url"`
	Secret  string `yaml:"secret"`
	Mention string `yaml:"mention"`
	Message string `yaml:"message"`
}
type Evaluate struct {
	Data template.Data `yaml:"data"`
	Type EvaluateType  `yaml:"type"`
}
type EvaluateType string

const (
	EvaluateEqual   EvaluateType = "equal"
	EvaluateInclude EvaluateType = "include"
)
