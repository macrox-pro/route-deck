package proxybridge

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sort"
	"strings"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// Compiler serializes a domain rule plan into macOS .pbprofile JSON.
type Compiler struct {
	Secrets port.SecretStore
	Logger  *slog.Logger
}

var _ port.ProfileCompiler = (*Compiler)(nil)

type pbProfile struct {
	DNSEnabled            bool      `json:"dnsEnabled"`
	TrafficLoggingEnabled bool      `json:"trafficLoggingEnabled"`
	CloseToMenuBar        bool      `json:"closeToMenuBar"`
	RunAtStartup          bool      `json:"runAtStartup"`
	ProxyConfigs          []pbProxy `json:"proxyConfigs"`
	ProxyRules            []pbRule  `json:"proxyRules"`
}

type pbProxy struct {
	ID         string `json:"id"`
	Name       string `json:"name,omitempty"`
	Type       string `json:"type"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	SendDomain bool   `json:"sendDomain"`
}

type pbRule struct {
	Name         string `json:"name,omitempty"`
	ProcessNames string `json:"processNames"`
	TargetHosts  string `json:"targetHosts"`
	TargetPorts  string `json:"targetPorts"`
	Protocol     string `json:"protocol"`
	Action       string `json:"action"`
	Enabled      bool   `json:"enabled"`
}

// Compile builds deterministic stdin JSON. Passwords load only here.
func (c *Compiler) Compile(ctx context.Context, profile domain.Profile, proxies []domain.ProxyEndpoint, identities []domain.ProcessIdentity) (port.CompiledProfile, error) {
	plan, err := domain.BuildRulePlan(profile, proxies, identities)
	if err != nil {
		return port.CompiledProfile{}, fmt.Errorf("compile profile: %w", err)
	}
	if len(proxies) > 100 {
		return port.CompiledProfile{}, fmt.Errorf("compile profile: too many proxies")
	}
	if len(plan.Rules) > 5000 {
		return port.CompiledProfile{}, fmt.Errorf("compile profile: too many rules")
	}

	sort.Slice(proxies, func(i, j int) bool { return proxies[i].ID < proxies[j].ID })
	pbProxies := make([]pbProxy, 0, len(proxies))
	for _, proxy := range proxies {
		password := ""
		if c.Secrets != nil && proxy.KeychainRef != "" {
			secret, err := c.Secrets.Load(ctx, proxy.KeychainRef)
			if err == nil {
				password = string(secret)
			}
		}
		if len(proxy.Username) > 255 || len(password) > 255 {
			return port.CompiledProfile{}, fmt.Errorf("compile profile: credential too long")
		}
		pbProxies = append(pbProxies, pbProxy{
			ID:         string(proxy.ID),
			Name:       proxy.Name,
			Type:       string(proxy.Kind),
			Host:       proxy.Host,
			Port:       proxy.Port,
			Username:   proxy.Username,
			Password:   password,
			SendDomain: proxy.SendDomain,
		})
	}

	pbRules := make([]pbRule, 0, len(plan.Rules))
	for i, rule := range plan.Rules {
		action := "DIRECT"
		if rule.Action.IsBlock() {
			action = "BLOCK"
		} else if id, ok := rule.Action.ProxyID(); ok {
			action = string(id)
		}
		protocol := string(rule.Protocol)
		if protocol == "" {
			protocol = string(domain.ProtocolScopeBOTH)
		}
		pbRules = append(pbRules, pbRule{
			Name:         fmt.Sprintf("rule-%d", i+1),
			ProcessNames: strings.Join(rule.ProcessNames, ";"),
			TargetHosts:  "*",
			TargetPorts:  "*",
			Protocol:     protocol,
			Action:       action,
			Enabled:      true,
		})
	}

	doc := pbProfile{
		DNSEnabled:            true,
		TrafficLoggingEnabled: false,
		CloseToMenuBar:        false,
		RunAtStartup:          false,
		ProxyConfigs:          pbProxies,
		ProxyRules:            pbRules,
	}
	raw, err := json.Marshal(doc)
	if err != nil {
		return port.CompiledProfile{}, fmt.Errorf("compile profile marshal: %w", err)
	}
	if len(raw) > 4<<20 {
		return port.CompiledProfile{}, fmt.Errorf("compile profile: exceeds 4 MiB")
	}
	if c.Logger != nil {
		c.Logger.Info("profile compiled",
			"component", "proxybridge",
			"profile_id", string(profile.ID),
			"proxy_count", len(pbProxies),
			"rule_count", len(pbRules),
		)
	}
	return port.CompiledProfile{JSON: raw, ProxyCount: len(pbProxies), RuleCount: len(pbRules)}, nil
}
