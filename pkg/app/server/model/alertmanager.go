package model

import (
	"time"
)

type AlertMessage struct {
	Email       []string    `json:"Email"`
	Labels      Labels      `json:"Labels"`
	Annotations Annotations `json:"Annotations"`
	StartsAt    time.Time   `json:"StartsAt"`
	EndsAt      time.Time   `json:"EndsAt"`
}
type Labels struct {
	Alertname string `json:"alertname"`
	IP        string `json:"IP"`
}
type Annotations struct {
	Summary string `json:"summary"`
}

type AlertSentMessage struct {
	Labels      Labels      `json:"Labels"`
	Annotations Annotations `json:"Annotations"`
	StartsAt    time.Time   `json:"StartsAt"`
	EndsAt      time.Time   `json:"EndsAt"`
}

// Alertmanager配置文件结构
type AlertmanagerYaml struct {
	Global    Global      `yaml:"global"`
	Route     Route       `yaml:"route"`
	Receivers []Receivers `yaml:"receivers"`
}

// Global
type Global struct {
	SmtpSmarthost    string `yaml:"smtp_smarthost"`
	SmtpFrom         string `yaml:"smtp_from"`
	SmtpAuthUsername string `yaml:"smtp_auth_username"`
	SmtpAuthPassword string `yaml:"smtp_auth_password"`
	ResolveTimeout   string `yaml:"resolve_timeout"`
	SmtpRequireTls   bool   `yaml:"smtp_require_tls"`
}

// Route
type Route struct {
	GroupBy        []string `yaml:"group_by"`
	GroupWait      string   `yaml:"group_wait"`
	GroupInterval  string   `yaml:"group_interval"`
	RepeatInterval string   `yaml:"repeat_interval"`
	Receiver       string   `yaml:"receiver"`
}

// Receivers
type Receivers struct {
	Name         string         `yaml:"name"`
	EmailConfigs []EmailConfigs `yaml:"email_configs"`
}

// EmailConfigs
type EmailConfigs struct {
	To string `yaml:"to"`
}
