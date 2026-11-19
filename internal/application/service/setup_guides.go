package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"sync"

	"github.com/macrox-pro/route-deck/internal/application/port"
	"github.com/macrox-pro/route-deck/internal/domain"
)

// Setup guide typed errors (ADR-0013).
var (
	ErrUnknownGuideID       = errors.New("unknown setup guide id")
	ErrGuideHasNoProxy      = errors.New("setup guide has no local proxy")
	ErrGuideNeedsPort       = errors.New("setup guide needs a port")
	ErrEndpointAlreadyAdded = errors.New("proxy endpoint already added")
)

// Guide detection and add modes.
const (
	GuideKindBridge  = "bridge"
	GuideKindClient  = "client"
	GuideKindGeneric = "generic"

	DetectionNone         = "none"
	DetectionUnavailable  = "unavailable"
	DetectionNotInstalled = "not_installed"
	DetectionInstalled    = "installed"

	AddModeNone         = "none"
	AddModeOneClick     = "one_click"
	AddModeFormFallback = "form_fallback"
)

// GuideStep is one numbered instruction.
type GuideStep struct {
	Title   string
	Body    string
	Warning bool
}

// GuideEndpoint is a documented local listener.
type GuideEndpoint struct {
	Host string
	Port int
	Kind domain.ProxyKind
	Note string
}

// GuideDefinition is static catalog content.
type GuideDefinition struct {
	ID          string
	Title       string
	Kind        string
	BundleIDs   []string
	Endpoints   []GuideEndpoint
	DocsURL     string
	Steps       []GuideStep
	AddMode     string
	DefaultName string
}

// GuideView is a detected guide for the UI.
type GuideView struct {
	ID            string
	Title         string
	Kind          string
	Detection     string
	InstalledPath string
	PortOpen      bool
	PortChecked   int
	AlreadyAdded  bool
	AddMode       string
	DocsURL       string
	Steps         []GuideStep
	Endpoints     []GuideEndpoint
	ConnectHint   bool
}

// SuggestionView is a banner candidate.
type SuggestionView struct {
	GuideID string
	Title   string
	Host    string
	Port    int
	Kind    string
}

// ChecklistProxyStep is the closed wire token for the proxy checklist row (UI copy maps by ID on the frontend).
type ChecklistProxyStep string

// ChecklistProxyStep values for ChecklistView.ProxyStep.
const (
	ChecklistProxyMissing    ChecklistProxyStep = "missing"
	ChecklistProxyAvailable  ChecklistProxyStep = "available"
	ChecklistProxyUDPLimited ChecklistProxyStep = "udp_limited"
	ChecklistProxyChecking   ChecklistProxyStep = "checking"
)

// ChecklistView is the first-run checklist.
type ChecklistView struct {
	BridgeFound      bool
	ProxyAvailable   bool
	ProxyPartial     bool
	ProxyStep        ChecklistProxyStep
	ApplicationAdded bool
	ProfileApplied   bool
	Complete         bool
	Dismissed        bool
}

// checklistFingerprint records step booleans at dismiss time for regression detection.
type checklistFingerprint struct {
	BridgeFound      bool
	ProxyAvailable   bool
	ApplicationAdded bool
	ProfileApplied   bool
}

// SetupGuides lists guides, adds proxies, suggestions, and checklist (ADR-0013 / ADR-0016).
type SetupGuides struct {
	Lister   port.InstalledApplicationLister
	Prober   port.LocalPortProber
	Proxies  port.ProxyStore
	Apps     port.ApplicationStore
	Catalog  *Catalog
	Health   *HealthMonitor
	Dismiss  port.DismissStore
	Snapshot func(ctx context.Context) (SnapshotResult, error)
	Logger   *slog.Logger

	mu                   sync.Mutex
	installed            []port.InstalledApplication
	installedOK          bool
	dismissed            map[string]struct{}
	checklistDismissed   bool
	checklistFingerprint checklistFingerprint
}

func (s *SetupGuides) logger() *slog.Logger {
	if s != nil && s.Logger != nil {
		return s.Logger
	}
	return slog.Default()
}

// DefaultGuideCatalog returns the ADR-0013 guide list.
func DefaultGuideCatalog() []GuideDefinition {
	warnTUN := GuideStep{
		Title:   "Use proxy-only mode",
		Body:    "Turn off TUN, system VPN, system proxy, or Enhanced Mode in this app. Those modes can capture other apps too, so a RouteDeck Direct policy may not be truly direct.",
		Warning: true,
	}
	warnLoop := GuideStep{
		Title:   "Do not route this app through itself",
		Body:    "Do not put this app in a RouteDeck Proxy policy that points at its own local port. That creates a routing loop.",
		Warning: true,
	}
	return []GuideDefinition{
		{
			ID:        "proxybridge",
			Title:     "ProxyBridge",
			Kind:      GuideKindBridge,
			BundleIDs: []string{"com.interceptsuite.ProxyBridge"},
			DocsURL:   "https://proxybridge.interceptsuite.com/",
			AddMode:   AddModeNone,
			Steps: []GuideStep{
				{Title: "Install ProxyBridge", Body: "Install the official ProxyBridge app from InterceptSuite. RouteDeck was tested with CLI 4.1.3."},
				{Title: "Approve the network extension", Body: "Open ProxyBridge once and approve the macOS network extension when prompted."},
				{Title: "Prefer the CLI path in Settings", Body: "RouteDeck discovers documented install locations. Set a custom CLI path in Settings only if discovery fails."},
				{Title: "Do not run the ProxyBridge GUI with RouteDeck", Body: "If the ProxyBridge GUI holds the extension, RouteDeck shows Conflict and will not take over a foreign process.", Warning: true},
			},
		},
		{
			ID:          "shadowrocket",
			Title:       "Shadowrocket",
			Kind:        GuideKindClient,
			BundleIDs:   []string{"com.liguangming.Shadowrocket"},
			Endpoints:   []GuideEndpoint{{Host: "127.0.0.1", Port: 1082, Kind: domain.ProxyKindSOCKS5, Note: "Proxy Share mixed port (SOCKS5 and HTTP)"}},
			DocsURL:     "https://apps.apple.com/app/shadowrocket/id932747118",
			AddMode:     AddModeOneClick,
			DefaultName: "Shadowrocket local",
			Steps: []GuideStep{
				{Title: "Open Shadowrocket", Body: "Launch Shadowrocket on this Mac and connect to a working node (top switch on)."},
				{Title: "Enable Proxy Share", Body: "Settings → Proxy → Proxy Share → Enable Share. Note the port (default 1082)."},
				{Title: "Limitation on Mac", Body: "Shadowrocket on Mac runs as a system VPN. Other apps may still tunnel through it even when RouteDeck sets Direct.", Warning: true},
				warnTUN,
				warnLoop,
			},
		},
		{
			ID:          "v2rayn",
			Title:       "v2rayN",
			Kind:        GuideKindClient,
			BundleIDs:   []string{"2dust.v2rayN"},
			Endpoints:   []GuideEndpoint{{Host: "127.0.0.1", Port: 10808, Kind: domain.ProxyKindSOCKS5, Note: "Mixed inbound in current builds; legacy HTTP was 10809"}},
			DocsURL:     "https://github.com/2dust/v2rayN",
			AddMode:     AddModeOneClick,
			DefaultName: "v2rayN local",
			Steps: []GuideStep{
				{Title: "Start v2rayN and a node", Body: "Import a subscription, select a node, and start the core."},
				{Title: "Confirm the local port", Body: "Settings → Parameter Settings. Use the SOCKS or mixed port (default 10808)."},
				{Title: "Turn off system proxy and TUN", Body: "Leave system proxy and TUN off so RouteDeck owns per-app routing.", Warning: true},
				warnTUN,
				warnLoop,
			},
		},
		{
			ID:          "hiddify",
			Title:       "Hiddify",
			Kind:        GuideKindClient,
			BundleIDs:   nil, // Detection unavailable until Mac bundle ID is sourced
			Endpoints:   []GuideEndpoint{{Host: "127.0.0.1", Port: 12334, Kind: domain.ProxyKindSOCKS5, Note: "Mixed port; closed while VPN mode is on"}},
			DocsURL:     "https://hiddify.com/app/Share-your-connection/",
			AddMode:     AddModeOneClick,
			DefaultName: "Hiddify local",
			Steps: []GuideStep{
				{Title: "Set Proxy service only", Body: "Settings → Inbound → Service mode → Proxy service only (not VPN)."},
				{Title: "Note the local proxy port", Body: "Default mixed port is 12334 (SOCKS5 and HTTP). Optional password is not used by one-click add."},
				warnTUN,
				warnLoop,
			},
		},
		{
			ID:          "happ",
			Title:       "Happ",
			Kind:        GuideKindClient,
			BundleIDs:   []string{"com.happ"},
			Endpoints:   []GuideEndpoint{{Host: "127.0.0.1", Port: 10808, Kind: domain.ProxyKindSOCKS5, Note: "SOCKS5; HTTP is often 10809"}},
			DocsURL:     "https://www.happ.su/main/faq/local-network-connections",
			AddMode:     AddModeOneClick,
			DefaultName: "Happ local",
			Steps: []GuideStep{
				{Title: "Prefer Proxy mode over TUN", Body: "Use a mode that exposes the local SOCKS5 port. Turn TUN off for RouteDeck per-app routing."},
				{Title: "Confirm SOCKS5 port 10808", Body: "If Allow LAN Connections is required for the local listener on your build, enable it only on a trusted network. LAN share often has no password.", Warning: true},
				warnTUN,
				warnLoop,
			},
		},
		{
			ID:          "singbox",
			Title:       "sing-box SFM",
			Kind:        GuideKindClient,
			BundleIDs:   []string{"io.nekohasekai.sfamt.standalone"},
			Endpoints:   nil,
			DocsURL:     "https://sing-box.sagernet.org/clients/apple/",
			AddMode:     AddModeFormFallback,
			DefaultName: "sing-box local",
			Steps: []GuideStep{
				{Title: "Add a mixed inbound", Body: "In the sing-box config, add a mixed inbound listening on 127.0.0.1 and choose a free port."},
				{Title: "Keep TUN off", Body: "Do not enable a TUN inbound while RouteDeck applies per-app rules.", Warning: true},
				{Title: "Enter the port in RouteDeck", Body: "This guide has no default port. Use Add this proxy to open the Proxies form, then type your listen port."},
				warnLoop,
			},
		},
		{
			ID:          "surge",
			Title:       "Surge",
			Kind:        GuideKindClient,
			BundleIDs:   []string{"com.nssurge.surge-mac"},
			Endpoints:   []GuideEndpoint{{Host: "127.0.0.1", Port: 6153, Kind: domain.ProxyKindSOCKS5, Note: "Default SOCKS5; HTTP is 6152"}},
			DocsURL:     "https://nssurge.com/",
			AddMode:     AddModeOneClick,
			DefaultName: "Surge local",
			Steps: []GuideStep{
				{Title: "Confirm listen ports", Body: "Default HTTP 6152 and SOCKS5 6153 (http-listen / socks5-listen). RouteDeck adds SOCKS5 6153."},
				{Title: "Turn Enhanced Mode and system proxy off", Body: "Enhanced Mode and system proxy capture traffic outside RouteDeck policies.", Warning: true},
				warnTUN,
				warnLoop,
			},
		},
		{
			ID:          "generic",
			Title:       "Other local proxy",
			Kind:        GuideKindGeneric,
			BundleIDs:   nil,
			Endpoints:   nil,
			DocsURL:     "",
			AddMode:     AddModeFormFallback,
			DefaultName: "Local proxy",
			Steps: []GuideStep{
				{Title: "Find the local SOCKS5 or HTTP port", Body: "In your VPN or proxy app, find the listen address (usually 127.0.0.1) and port."},
				{Title: "Use proxy-only mode", Body: "Disable TUN, system VPN, and system proxy so traffic is not captured outside RouteDeck.", Warning: true},
				{Title: "Add the proxy in RouteDeck", Body: "Press Add this proxy to open the Proxies form with 127.0.0.1 and SOCKS5, then enter the port and save."},
				warnLoop,
			},
		},
	}
}

func guideByID(id string) (GuideDefinition, bool) {
	for _, g := range DefaultGuideCatalog() {
		if g.ID == id {
			return g, true
		}
	}
	return GuideDefinition{}, false
}

// RefreshInstalled reloads the installed-application cache.
func (s *SetupGuides) RefreshInstalled(ctx context.Context) error {
	if s == nil || s.Lister == nil {
		s.mu.Lock()
		s.installed = nil
		s.installedOK = false
		s.mu.Unlock()
		return nil
	}
	rows, err := s.Lister.List(ctx)
	if err != nil {
		s.logger().Warn("setup guides list installed failed", "component", "setup_guides", "err", err)
		s.mu.Lock()
		s.installed = nil
		s.installedOK = false
		s.mu.Unlock()
		return fmt.Errorf("setup guides refresh installed: %w", err)
	}
	s.mu.Lock()
	s.installed = rows
	s.installedOK = true
	s.mu.Unlock()
	return nil
}

func (s *SetupGuides) installedCached(ctx context.Context) []port.InstalledApplication {
	s.mu.Lock()
	ok := s.installedOK
	rows := append([]port.InstalledApplication(nil), s.installed...)
	s.mu.Unlock()
	if ok {
		return rows
	}
	_ = s.RefreshInstalled(ctx)
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]port.InstalledApplication(nil), s.installed...)
}

func matchInstalled(def GuideDefinition, apps []port.InstalledApplication) (path string, detection string) {
	if def.Kind == GuideKindGeneric {
		return "", DetectionNone
	}
	if len(def.BundleIDs) == 0 {
		return "", DetectionUnavailable
	}
	want := map[string]struct{}{}
	for _, id := range def.BundleIDs {
		want[id] = struct{}{}
	}
	for _, app := range apps {
		if _, ok := want[app.BundleID]; ok {
			return app.Path, DetectionInstalled
		}
	}
	return "", DetectionNotInstalled
}

func (s *SetupGuides) probeEndpoint(ctx context.Context, ep GuideEndpoint) bool {
	if ep.Port < 1 || s == nil || s.Prober == nil {
		return false
	}
	open, err := s.Prober.Open(ctx, ep.Port)
	if err != nil {
		s.logger().Debug("setup guides port probe error", "component", "setup_guides", "port", ep.Port, "err", err)
		return false
	}
	return open
}

func (s *SetupGuides) alreadyAdded(ctx context.Context, host string, port int) (bool, error) {
	if s == nil || s.Proxies == nil || host == "" || port < 1 {
		return false, nil
	}
	proxies, err := s.Proxies.List(ctx)
	if err != nil {
		return false, err
	}
	for _, p := range proxies {
		if p.Host == host && p.Port == port {
			return true, nil
		}
	}
	return false, nil
}

func sortRank(v GuideView) int {
	switch {
	case v.Kind == GuideKindBridge:
		return 0
	case v.Kind == GuideKindGeneric:
		return 50
	case v.Detection == DetectionInstalled && v.PortOpen:
		return 10
	case v.Detection == DetectionInstalled:
		return 20
	case v.Detection == DetectionNotInstalled:
		return 30
	case v.Detection == DetectionUnavailable:
		return 40
	default:
		return 45
	}
}

// List returns detected guides. refresh=true reloads the installed-app cache.
func (s *SetupGuides) List(ctx context.Context, refresh bool) ([]GuideView, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if refresh {
		_ = s.RefreshInstalled(ctx)
	}
	apps := s.installedCached(ctx)
	out := make([]GuideView, 0, len(DefaultGuideCatalog()))
	for _, def := range DefaultGuideCatalog() {
		path, detection := matchInstalled(def, apps)
		view := GuideView{
			ID:            def.ID,
			Title:         def.Title,
			Kind:          def.Kind,
			Detection:     detection,
			InstalledPath: path,
			AddMode:       def.AddMode,
			DocsURL:       def.DocsURL,
			Steps:         append([]GuideStep(nil), def.Steps...),
			Endpoints:     append([]GuideEndpoint(nil), def.Endpoints...),
		}
		if len(def.Endpoints) > 0 {
			ep := def.Endpoints[0]
			view.PortChecked = ep.Port
			view.PortOpen = s.probeEndpoint(ctx, ep)
			view.ConnectHint = detection == DetectionInstalled && !view.PortOpen
			added, err := s.alreadyAdded(ctx, ep.Host, ep.Port)
			if err != nil {
				return nil, fmt.Errorf("setup guides list: %w", err)
			}
			view.AlreadyAdded = added
		}
		out = append(out, view)
	}
	sort.SliceStable(out, func(i, j int) bool {
		ri, rj := sortRank(out[i]), sortRank(out[j])
		if ri != rj {
			return ri < rj
		}
		return strings.ToLower(out[i].Title) < strings.ToLower(out[j].Title)
	})
	return out, nil
}

// DocsURL returns the catalog docs URL for a guide.
func (s *SetupGuides) DocsURL(guideID string) (string, error) {
	def, ok := guideByID(guideID)
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnknownGuideID, guideID)
	}
	if def.DocsURL == "" {
		return "", fmt.Errorf("%w: %s has no docs url", ErrUnknownGuideID, guideID)
	}
	return def.DocsURL, nil
}

// AddProxy saves the guide's default endpoint and rechecks health.
func (s *SetupGuides) AddProxy(ctx context.Context, guideID string) (domain.ProxyEndpoint, error) {
	def, ok := guideByID(guideID)
	if !ok {
		return domain.ProxyEndpoint{}, fmt.Errorf("%w: %s", ErrUnknownGuideID, guideID)
	}
	switch def.AddMode {
	case AddModeNone:
		return domain.ProxyEndpoint{}, fmt.Errorf("%w: %s", ErrGuideHasNoProxy, guideID)
	case AddModeFormFallback:
		return domain.ProxyEndpoint{}, fmt.Errorf("%w: %s", ErrGuideNeedsPort, guideID)
	}
	if len(def.Endpoints) == 0 {
		return domain.ProxyEndpoint{}, fmt.Errorf("%w: %s", ErrGuideNeedsPort, guideID)
	}
	ep := def.Endpoints[0]
	added, err := s.alreadyAdded(ctx, ep.Host, ep.Port)
	if err != nil {
		return domain.ProxyEndpoint{}, fmt.Errorf("setup guides add proxy: %w", err)
	}
	if added {
		return domain.ProxyEndpoint{}, fmt.Errorf("%w: %s:%d", ErrEndpointAlreadyAdded, ep.Host, ep.Port)
	}
	if s.Catalog == nil {
		return domain.ProxyEndpoint{}, fmt.Errorf("setup guides add proxy: catalog not configured")
	}
	id, err := NewID()
	if err != nil {
		return domain.ProxyEndpoint{}, err
	}
	name := def.DefaultName
	if name == "" {
		name = def.Title + " local"
	}
	name, err = s.uniqueProxyName(ctx, name)
	if err != nil {
		return domain.ProxyEndpoint{}, err
	}
	endpoint, err := domain.NewProxyEndpoint(domain.ProxyID(id), name, ep.Kind, ep.Host, ep.Port, domain.KeychainRef(id))
	if err != nil {
		return domain.ProxyEndpoint{}, fmt.Errorf("setup guides add proxy: %w", err)
	}
	if err := s.Catalog.SaveProxy(ctx, endpoint); err != nil {
		return domain.ProxyEndpoint{}, err
	}
	if s.Health != nil {
		if err := s.Health.Recheck(ctx); err != nil {
			s.logger().Warn("setup guides recheck after add", "component", "setup_guides", "err", err)
		}
	}
	return endpoint, nil
}

func (s *SetupGuides) uniqueProxyName(ctx context.Context, base string) (string, error) {
	if s.Proxies == nil {
		return base, nil
	}
	proxies, err := s.Proxies.List(ctx)
	if err != nil {
		return "", err
	}
	used := map[string]struct{}{}
	for _, p := range proxies {
		used[p.Name] = struct{}{}
	}
	if _, ok := used[base]; !ok {
		return base, nil
	}
	for i := 2; i < 100; i++ {
		candidate := fmt.Sprintf("%s (%d)", base, i)
		if _, ok := used[candidate]; !ok {
			return candidate, nil
		}
	}
	return base + " " + mustShortID(), nil
}

func mustShortID() string {
	id, err := NewID()
	if err != nil {
		return "x"
	}
	return id[:8]
}

// Suggestions returns install+port-open guides not yet in the catalog and not dismissed.
func (s *SetupGuides) Suggestions(ctx context.Context) ([]SuggestionView, error) {
	guides, err := s.List(ctx, false)
	if err != nil {
		return nil, err
	}
	s.mu.Lock()
	dismissed := s.dismissed
	s.mu.Unlock()
	out := make([]SuggestionView, 0)
	for _, g := range guides {
		if g.Kind != GuideKindClient {
			continue
		}
		if g.Detection != DetectionInstalled || !g.PortOpen || g.AlreadyAdded {
			continue
		}
		if len(g.Endpoints) == 0 {
			continue
		}
		ep := g.Endpoints[0]
		key := dismissKey(g.ID, ep.Port)
		if dismissed != nil {
			if _, ok := dismissed[key]; ok {
				continue
			}
		}
		out = append(out, SuggestionView{
			GuideID: g.ID,
			Title:   g.Title,
			Host:    ep.Host,
			Port:    ep.Port,
			Kind:    string(ep.Kind),
		})
	}
	return out, nil
}

func dismissKey(guideID string, port int) string {
	return fmt.Sprintf("%s:%d", guideID, port)
}

// LoadDismiss loads checklist and suggestion dismiss from the store into memory (ADR-0016).
func (s *SetupGuides) LoadDismiss(ctx context.Context) error {
	if s == nil || s.Dismiss == nil {
		return nil
	}
	checklist, err := s.Dismiss.GetChecklist(ctx)
	if err != nil {
		return fmt.Errorf("setup guides load checklist dismiss: %w", err)
	}
	suggestions, err := s.Dismiss.ListSuggestions(ctx)
	if err != nil {
		return fmt.Errorf("setup guides load suggestion dismiss: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checklistDismissed = checklist.Dismissed
	s.checklistFingerprint = checklistFingerprint{
		BridgeFound:      checklist.BridgeFound,
		ProxyAvailable:   checklist.ProxyAvailable,
		ApplicationAdded: checklist.ApplicationAdded,
		ProfileApplied:   checklist.ProfileApplied,
	}
	s.dismissed = map[string]struct{}{}
	for _, row := range suggestions {
		s.dismissed[dismissKey(row.GuideID, row.Port)] = struct{}{}
	}
	return nil
}

// DismissSuggestion hides a suggestion across process restart (ADR-0016).
func (s *SetupGuides) DismissSuggestion(ctx context.Context, guideID string, port int) error {
	if _, ok := guideByID(guideID); !ok {
		return fmt.Errorf("%w: %s", ErrUnknownGuideID, guideID)
	}
	if s.Dismiss != nil {
		if err := s.Dismiss.AddSuggestion(ctx, guideID, port); err != nil {
			return fmt.Errorf("setup guides dismiss suggestion: %w", err)
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.dismissed == nil {
		s.dismissed = map[string]struct{}{}
	}
	s.dismissed[dismissKey(guideID, port)] = struct{}{}
	return nil
}

// Checklist derives first-run steps from snapshot and catalog.
func (s *SetupGuides) Checklist(ctx context.Context) (ChecklistView, error) {
	view, err := s.deriveChecklist(ctx)
	if err != nil {
		return ChecklistView{}, err
	}
	s.mu.Lock()
	cleared := false
	if s.checklistDismissed && checklistRegressed(s.checklistFingerprint, view) {
		s.checklistDismissed = false
		s.checklistFingerprint = checklistFingerprint{}
		cleared = true
	}
	view.Dismissed = s.checklistDismissed
	s.mu.Unlock()
	if cleared && s.Dismiss != nil {
		if err := s.Dismiss.SaveChecklist(ctx, port.ChecklistDismiss{}); err != nil {
			return view, fmt.Errorf("setup guides clear checklist dismiss: %w", err)
		}
	}
	return view, nil
}

// DismissChecklist hides the first-run checklist until step regression (ADR-0016).
func (s *SetupGuides) DismissChecklist(ctx context.Context) error {
	view, err := s.deriveChecklist(ctx)
	if err != nil {
		return err
	}
	fp := checklistFingerprint{
		BridgeFound:      view.BridgeFound,
		ProxyAvailable:   view.ProxyAvailable,
		ApplicationAdded: view.ApplicationAdded,
		ProfileApplied:   view.ProfileApplied,
	}
	if s.Dismiss != nil {
		if err := s.Dismiss.SaveChecklist(ctx, port.ChecklistDismiss{
			Dismissed:        true,
			BridgeFound:      fp.BridgeFound,
			ProxyAvailable:   fp.ProxyAvailable,
			ApplicationAdded: fp.ApplicationAdded,
			ProfileApplied:   fp.ProfileApplied,
		}); err != nil {
			return fmt.Errorf("setup guides dismiss checklist: %w", err)
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.checklistDismissed = true
	s.checklistFingerprint = fp
	return nil
}

func (s *SetupGuides) deriveChecklist(ctx context.Context) (ChecklistView, error) {
	view := ChecklistView{
		ProxyStep: ChecklistProxyMissing,
	}
	if s.Snapshot == nil {
		return view, fmt.Errorf("setup guides checklist: snapshot not configured")
	}
	snap, err := s.Snapshot(ctx)
	if err != nil {
		return view, fmt.Errorf("setup guides checklist: %w", err)
	}
	view.BridgeFound = snap.Bridge != domain.BridgeStateNotInstalled && snap.Bridge != domain.BridgeStateUnknown
	applyProxyChecklistStep(&view, snap.Proxies)
	if s.Apps != nil {
		apps, err := s.Apps.List(ctx)
		if err != nil {
			return view, fmt.Errorf("setup guides checklist apps: %w", err)
		}
		view.ApplicationAdded = len(apps) > 0
	}
	view.ProfileApplied = snap.AppliedConfirmed
	view.Complete = view.BridgeFound && view.ProxyAvailable && view.ApplicationAdded && view.ProfileApplied
	return view, nil
}

func applyProxyChecklistStep(view *ChecklistView, proxies []ProxyHealthRow) {
	best := ChecklistProxyMissing
	for _, p := range proxies {
		candidate := checklistProxyStepForHealth(p.Health)
		if checklistProxyStepRank(candidate) > checklistProxyStepRank(best) {
			best = candidate
		}
		if best == ChecklistProxyAvailable {
			break
		}
	}
	view.ProxyStep = best
	view.ProxyAvailable = best != ChecklistProxyMissing
	view.ProxyPartial = best == ChecklistProxyUDPLimited || best == ChecklistProxyChecking
}

func checklistProxyStepForHealth(health domain.ProxyHealth) ChecklistProxyStep {
	switch health {
	case domain.ProxyHealthAvailable:
		return ChecklistProxyAvailable
	case domain.ProxyHealthTCPAvailableUDPUnavailable:
		return ChecklistProxyUDPLimited
	case domain.ProxyHealthChecking:
		return ChecklistProxyChecking
	case domain.ProxyHealthUnknown,
		domain.ProxyHealthUnavailable,
		domain.ProxyHealthAuthenticationFailed:
		return ChecklistProxyMissing
	}
	return ChecklistProxyMissing
}

func checklistProxyStepRank(step ChecklistProxyStep) int {
	switch step {
	case ChecklistProxyAvailable:
		return 3
	case ChecklistProxyUDPLimited:
		return 2
	case ChecklistProxyChecking:
		return 1
	case ChecklistProxyMissing:
		return 0
	}
	return 0
}

func checklistRegressed(fp checklistFingerprint, view ChecklistView) bool {
	return (fp.BridgeFound && !view.BridgeFound) ||
		(fp.ProxyAvailable && !view.ProxyAvailable) ||
		(fp.ApplicationAdded && !view.ApplicationAdded) ||
		(fp.ProfileApplied && !view.ProfileApplied)
}
