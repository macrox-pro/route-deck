package domain

// ProxyFact is the id and health of one proxy supplied by the caller.
type ProxyFact struct {
	ID     ProxyID
	Health ProxyHealth
}

// AppCoverage is one application's action and protocol coverage facts.
type AppCoverage struct {
	ID          ApplicationID
	Action      RoutingAction
	TCPCovered  bool
	UDPCovered  bool
	UDPRequired bool
}

// ProtectionFacts is the input DeriveProtection reads. The caller supplies coverage.
type ProtectionFacts struct {
	Bridge           BridgeState
	VersionKnown     bool
	Desired          ProfileID
	DesiredSet       bool
	Applied          ProfileID
	AppliedConfirmed bool
	ConfirmFailed    bool
	RollbackFailed   bool
	RecheckPending   bool
	DatabaseDamaged  bool
	Proxies          []ProxyFact
	Apps             []AppCoverage
	DomainRules      bool
	DNSIncomplete    bool
}

// DeriveProtection returns the first matching protection step.
func DeriveProtection(facts ProtectionFacts) ProtectionState {
	lookup := func(id ProxyID) (ProxyHealth, bool) {
		for _, proxy := range facts.Proxies {
			if proxy.ID == id {
				return proxy.Health, true
			}
		}
		return "", false
	}

	if facts.RecheckPending {
		return ProtectionStateUnknown
	}

	for _, app := range facts.Apps {
		id, ok := app.Action.ProxyID()
		if !ok {
			continue
		}
		health, found := lookup(id)
		if !found {
			continue
		}
		if health == ProxyHealthUnknown || health == ProxyHealthChecking {
			return ProtectionStateUnknown
		}
	}
	if facts.Bridge == BridgeStateUnknown || !facts.VersionKnown || facts.DatabaseDamaged {
		return ProtectionStateUnknown
	}

	if facts.Bridge == BridgeStateStarting || facts.Bridge == BridgeStateStopping {
		return ProtectionStateInitializing
	}

	if (facts.ConfirmFailed && facts.RollbackFailed) || (facts.ConfirmFailed && facts.Applied == "") {
		return ProtectionStateUnprotected
	}

	if facts.Bridge == BridgeStateFailed {
		return ProtectionStateFailed
	}

	if facts.Bridge == BridgeStateConflict {
		return ProtectionStateUnprotected
	}

	if facts.Bridge == BridgeStateNotInstalled || facts.Bridge == BridgeStateDetected {
		return ProtectionStateInactive
	}

	if facts.Bridge == BridgeStateStopped && facts.DesiredSet {
		return ProtectionStateUnprotected
	}

	if facts.Bridge == BridgeStateStopped && !facts.DesiredSet {
		return ProtectionStateInactive
	}

	if facts.Bridge == BridgeStateReady && !facts.AppliedConfirmed {
		return ProtectionStateUnprotected
	}

	if facts.Bridge != BridgeStateReady {
		return ""
	}

	protected := 0
	anyCovered := false
	allFullyCovered := true
	gap := facts.DomainRules && facts.DNSIncomplete
	allProxiesAvailable := true

	for _, app := range facts.Apps {
		isProtected := app.Action.IsBlock()
		proxyID, isProxy := app.Action.ProxyID()
		if isProxy {
			isProtected = true
		}
		if !isProtected {
			continue
		}
		protected++
		fullyCovered := app.TCPCovered && (!app.UDPRequired || app.UDPCovered)
		if !fullyCovered {
			allFullyCovered = false
			gap = true
		}
		if app.TCPCovered || app.UDPCovered {
			anyCovered = true
		}
		if isProxy {
			health, found := lookup(proxyID)
			if found {
				if health != ProxyHealthAvailable {
					allProxiesAvailable = false
				}
				if health == ProxyHealthUnavailable ||
					health == ProxyHealthAuthenticationFailed ||
					health == ProxyHealthTCPAvailableUDPUnavailable {
					gap = true
				}
			}
		}
	}

	if protected == 0 {
		return ProtectionStateInactive
	}
	if allFullyCovered && allProxiesAvailable && (!facts.DomainRules || !facts.DNSIncomplete) {
		return ProtectionStateActive
	}
	if anyCovered && gap {
		return ProtectionStateDegraded
	}
	if !anyCovered {
		return ProtectionStateUnprotected
	}
	return ""
}
