package sqlite

// applicationModel maps domain.Application.
type applicationModel struct {
	ID          string `gorm:"column:id;primaryKey"`
	DisplayName string `gorm:"column:display_name;not null"`
	BundleID    string `gorm:"column:bundle_id;not null"`
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
	KeychainRef string `gorm:"column:keychain_ref;not null"`
	SendDomain  bool   `gorm:"column:send_domain;not null"`
}

func (proxyModel) TableName() string { return "proxies" }

// profileModel maps domain.Profile.
type profileModel struct {
	ID   string `gorm:"column:id;primaryKey"`
	Name string `gorm:"column:name;not null"`
}

func (profileModel) TableName() string { return "profiles" }

// applicationPolicyModel maps one ApplicationPolicy inside a profile.
type applicationPolicyModel struct {
	ProfileID  string               `gorm:"column:profile_id;primaryKey"`
	ProcessID  string               `gorm:"column:process_id;primaryKey"`
	ActionKind string               `gorm:"column:action_kind;not null"`
	ProxyID    *string              `gorm:"column:proxy_id"`
	Profile    profileModel         `gorm:"foreignKey:ProfileID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Process    processIdentityModel `gorm:"foreignKey:ProcessID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Proxy      *proxyModel          `gorm:"foreignKey:ProxyID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (applicationPolicyModel) TableName() string { return "application_policies" }
