package sqlite

// applicationModel maps domain.Application.
type applicationModel struct {
	ID          string `gorm:"column:id;primaryKey"`
	DisplayName string `gorm:"column:display_name;not null"`
	BundleID    string `gorm:"column:bundle_id;not null;uniqueIndex"`
	IconRef     string `gorm:"column:icon_ref;not null"`
}

func (applicationModel) TableName() string { return "applications" }

// processIdentityModel maps domain.ProcessIdentity.
type processIdentityModel struct {
	ID             string           `gorm:"column:id;primaryKey"`
	ApplicationID  string           `gorm:"column:application_id;not null;index"`
	BundleID       string           `gorm:"column:bundle_id;not null"`
	ExecutableName string           `gorm:"column:executable_name;not null"`
	Path           string           `gorm:"column:path;not null"`
	SigningHint    string           `gorm:"column:signing_hint;not null"`
	Application    applicationModel `gorm:"foreignKey:ApplicationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (processIdentityModel) TableName() string { return "process_identities" }

// proxyModel maps domain.ProxyEndpoint. Passwords are never stored.
type proxyModel struct {
	ID          string `gorm:"column:id;primaryKey"`
	Name        string `gorm:"column:name;not null"`
	Kind        string `gorm:"column:kind;not null"`
	Host        string `gorm:"column:host;not null"`
	Port        int    `gorm:"column:port;not null"`
	Username    string `gorm:"column:username;not null"`
	KeychainRef string `gorm:"column:keychain_ref;not null"`
	SendDomain  bool   `gorm:"column:send_domain;not null"`
}

func (proxyModel) TableName() string { return "proxies" }

// profileModel maps domain.Profile.
type profileModel struct {
	ID   string `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name;not null;uniqueIndex"`
}

func (profileModel) TableName() string { return "profiles" }

// applicationPolicyModel maps one ApplicationPolicy inside a profile.
type applicationPolicyModel struct {
	ProfileID     string               `gorm:"column:profile_id;primaryKey"`
	ProcessID     string               `gorm:"column:process_id;primaryKey"`
	ActionKind    string               `gorm:"column:action_kind;not null"`
	ProxyID       *string              `gorm:"column:proxy_id"`
	ProtocolScope string               `gorm:"column:protocol_scope;not null"`
	UDPLimitation string               `gorm:"column:udp_limitation;not null"`
	Profile       profileModel         `gorm:"foreignKey:ProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Process       processIdentityModel `gorm:"foreignKey:ProcessID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Proxy         *proxyModel          `gorm:"foreignKey:ProxyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (applicationPolicyModel) TableName() string { return "application_policies" }

// settingsModel is a single-row settings table (id always 1).
type settingsModel struct {
	ID               int    `gorm:"column:id;primaryKey"`
	DesiredProfileID string `gorm:"column:desired_profile_id;not null"`
	DesiredSet       bool   `gorm:"column:desired_set;not null"`
	ProxyBridgePath  string `gorm:"column:proxybridge_path;not null"`
	DebugLogging     bool   `gorm:"column:debug_logging;not null"`
	// default:system is required so SQLite AutoMigrate can ADD COLUMN on existing DBs
	// (NOT NULL without DEFAULT fails with "default value NULL").
	Appearance string `gorm:"column:appearance;not null;default:system"`
}

func (settingsModel) TableName() string { return "settings" }
