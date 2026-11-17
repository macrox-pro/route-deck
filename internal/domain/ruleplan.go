package domain

import "fmt"

// PlannedRule is a wire-neutral ProxyBridge rule before JSON serialization.
type PlannedRule struct {
	Name         string
	ProcessNames []string
	Protocol     ProtocolScope
	Action       RoutingAction
}

// RulePlan is the ordered compile result plus coverage facts for DeriveProtection.
type RulePlan struct {
	Rules []PlannedRule
	Apps  []AppCoverage
}

// BuildRulePlan compiles profile policies into rules and AppCoverage facts (ADR-0008 / ADR-0014).
func BuildRulePlan(profile Profile, proxies []ProxyEndpoint, identities []ProcessIdentity) (RulePlan, error) {
	known := make([]ProxyID, 0, len(proxies))
	proxyByID := make(map[ProxyID]ProxyEndpoint, len(proxies))
	for _, proxy := range proxies {
		known = append(known, proxy.ID)
		proxyByID[proxy.ID] = proxy
	}
	if err := ValidateProxyRefs(profile, known); err != nil {
		return RulePlan{}, err
	}
	if err := ValidateMatchConflicts(profile, identities); err != nil {
		return RulePlan{}, fmt.Errorf("build rule plan: %w", err)
	}
	identityByID := make(map[ProcessIdentityID]ProcessIdentity, len(identities))
	for _, identity := range identities {
		identityByID[identity.ID] = identity
	}

	// Count stored identities per bundle id (ADR-0014 collapse of pre-afb231c duplicates).
	bundleCounts := make(map[string]int)
	for _, identity := range identities {
		if identity.BundleID == "" {
			continue
		}
		bundleCounts[identity.BundleID]++
	}

	emittedBundle := make(map[string]struct{})
	var rules []PlannedRule
	var apps []AppCoverage
	for _, policy := range profile.Policies {
		if !policy.Action.IsSet() {
			return RulePlan{}, fmt.Errorf("build rule plan: %w", ErrUnsetAction)
		}
		identity, ok := identityByID[policy.ProcessID]
		if !ok {
			return RulePlan{}, fmt.Errorf("build rule plan: unknown process %q", policy.ProcessID)
		}

		// Collapse: several identities share a bundle id → one rule with that bundle id only.
		if identity.BundleID != "" && bundleCounts[identity.BundleID] > 1 {
			if _, done := emittedBundle[identity.BundleID]; done {
				continue
			}
			emittedBundle[identity.BundleID] = struct{}{}
		}

		names := identity.MatchableNames()
		if identity.BundleID != "" && bundleCounts[identity.BundleID] > 1 {
			names = []string{identity.BundleID}
		}
		if len(names) == 0 {
			return RulePlan{}, fmt.Errorf("build rule plan: no matchable names for %q", policy.ProcessID)
		}

		scope := policy.ProtocolScope
		if scope == "" {
			scope = ProtocolScopeBOTH
		}
		limitation := policy.UDPLimitation
		if limitation == "" {
			limitation = UDPLimitationAbsent
		}

		coverage := AppCoverage{
			ID:     identity.ApplicationID,
			Action: policy.Action,
		}

		switch {
		case policy.Action.IsDirect(), policy.Action.IsBlock():
			rules = append(rules, PlannedRule{
				ProcessNames: names,
				Protocol:     scope,
				Action:       policy.Action,
			})
			if scope == ProtocolScopeTCP || scope == ProtocolScopeBOTH {
				coverage.TCPCovered = policy.Action.IsBlock()
			}
			if scope == ProtocolScopeUDP || scope == ProtocolScopeBOTH {
				coverage.UDPRequired = true
				coverage.UDPCovered = policy.Action.IsBlock()
			}
			if policy.Action.IsDirect() {
				coverage.TCPCovered = false
				coverage.UDPCovered = false
			}
		default:
			proxyID, _ := policy.Action.ProxyID()
			proxy := proxyByID[proxyID]
			inScopeRules, outScopeRules, cov := compileProxyPolicy(names, scope, limitation, proxy, policy.Action)
			rules = append(rules, inScopeRules...)
			rules = append(rules, outScopeRules...)
			coverage = cov
			coverage.ID = identity.ApplicationID
		}
		apps = append(apps, coverage)
	}
	return RulePlan{Rules: rules, Apps: apps}, nil
}

func compileProxyPolicy(names []string, scope ProtocolScope, limitation UDPLimitation, proxy ProxyEndpoint, action RoutingAction) (inScope, outScope []PlannedRule, coverage AppCoverage) {
	coverage.Action = action
	switch scope {
	case ProtocolScopeTCP:
		inScope = append(inScope, PlannedRule{ProcessNames: names, Protocol: ProtocolScopeTCP, Action: action})
		outScope = append(outScope, PlannedRule{ProcessNames: names, Protocol: ProtocolScopeUDP, Action: BlockAction()})
		coverage.TCPCovered = true
		coverage.UDPRequired = true
		coverage.UDPCovered = true // Block covers UDP
	case ProtocolScopeUDP:
		inScope = append(inScope, compileUDPProxyRule(names, limitation, proxy, action)...)
		outScope = append(outScope, PlannedRule{ProcessNames: names, Protocol: ProtocolScopeTCP, Action: BlockAction()})
		coverage.TCPCovered = true // Block covers TCP
		coverage.UDPRequired = true
		coverage.UDPCovered = udpCovered(proxy, limitation)
	case ProtocolScopeBOTH:
		if proxy.Kind == ProxyKindHTTP {
			inScope = append(inScope, PlannedRule{ProcessNames: names, Protocol: ProtocolScopeTCP, Action: action})
			inScope = append(inScope, compileUDPProxyRule(names, limitation, proxy, action)...)
			coverage.TCPCovered = true
			coverage.UDPRequired = true
			coverage.UDPCovered = udpCovered(proxy, limitation)
		} else {
			inScope = append(inScope, PlannedRule{ProcessNames: names, Protocol: ProtocolScopeBOTH, Action: action})
			coverage.TCPCovered = true
			coverage.UDPRequired = true
			coverage.UDPCovered = true
		}
	}
	return inScope, outScope, coverage
}

func compileUDPProxyRule(names []string, limitation UDPLimitation, proxy ProxyEndpoint, action RoutingAction) []PlannedRule {
	if proxy.Kind == ProxyKindHTTP {
		if limitation == UDPLimitationDirectUDP {
			return []PlannedRule{{ProcessNames: names, Protocol: ProtocolScopeUDP, Action: DirectAction()}}
		}
		return []PlannedRule{{ProcessNames: names, Protocol: ProtocolScopeUDP, Action: BlockAction()}}
	}
	return []PlannedRule{{ProcessNames: names, Protocol: ProtocolScopeUDP, Action: action}}
}

func udpCovered(proxy ProxyEndpoint, limitation UDPLimitation) bool {
	if proxy.Kind == ProxyKindHTTP {
		return false // BLOCK or DIRECT-UDP: UDP not covered per CONCEPTION
	}
	_ = limitation
	return true
}
