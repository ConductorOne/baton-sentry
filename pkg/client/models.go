package client

import (
	"time"
)

type Organization struct {
	Avatar                     Avatar             `json:"avatar"`
	DateCreated                time.Time          `json:"dateCreated"`
	Features                   []string           `json:"features"`
	HasAuthProvider            bool               `json:"hasAuthProvider"`
	ID                         string             `json:"id"`
	IsEarlyAdopter             bool               `json:"isEarlyAdopter"`
	AllowMemberInvite          bool               `json:"allowMemberInvite"`
	AllowMemberProjectCreation bool               `json:"allowMemberProjectCreation"`
	AllowSuperuserAccess       bool               `json:"allowSuperuserAccess"`
	Links                      OrganizationLinks  `json:"links"`
	Name                       string             `json:"name"`
	Require2FA                 bool               `json:"require2FA"`
	Slug                       string             `json:"slug"`
	Status                     OrganizationStatus `json:"status"`
}

type Avatar struct {
	AvatarType string  `json:"avatarType"`
	AvatarUUID *string `json:"avatarUuid"`
}

type OrganizationLinks struct {
	OrganizationURL string `json:"organizationUrl"`
	RegionURL       string `json:"regionUrl"`
}

type OrganizationStatus struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type OrganizationMember struct {
	ID           string      `json:"id"`
	Email        string      `json:"email"`
	Name         string      `json:"name"`
	User         *User       `json:"user,omitempty"` // Optional - only present for active members
	OrgRole      string      `json:"orgRole"`
	Pending      bool        `json:"pending"`
	Expired      bool        `json:"expired"`
	Flags        MemberFlags `json:"flags"`
	DateCreated  time.Time   `json:"dateCreated"`
	InviteStatus string      `json:"inviteStatus"`
	InviterName  string      `json:"inviterName"`
}

type User struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Username        string         `json:"username"`
	Email           string         `json:"email"`
	AvatarURL       string         `json:"avatarUrl"`
	IsActive        bool           `json:"isActive"`
	HasPasswordAuth bool           `json:"hasPasswordAuth"`
	IsManaged       bool           `json:"isManaged"`
	DateJoined      time.Time      `json:"dateJoined"`
	LastLogin       *time.Time     `json:"lastLogin"`
	Has2FA          bool           `json:"has2fa"`
	LastActive      *time.Time     `json:"lastActive"`
	IsSuperuser     bool           `json:"isSuperuser"`
	IsStaff         bool           `json:"isStaff"`
	Experiments     map[string]any `json:"experiments"`
	Emails          []UserEmail    `json:"emails"`
	Avatar          Avatar         `json:"avatar"`
	CanReset2FA     bool           `json:"canReset2fa"`
}

type UserEmail struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
}

type MemberFlags struct {
	IDPProvisioned        bool `json:"idp:provisioned"`
	IDPRoleRestricted     bool `json:"idp:role-restricted"`
	SSOLinked             bool `json:"sso:linked"`
	SSOInvalid            bool `json:"sso:invalid"`
	MemberLimitRestricted bool `json:"member-limit:restricted"`
	PartnershipRestricted bool `json:"partnership:restricted"`
}

type Team struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	DateCreated time.Time `json:"dateCreated"`
	IsMember    bool      `json:"isMember"`
	TeamRole    *string   `json:"teamRole"`
	Flags       TeamFlags `json:"flags"`
	Access      []string  `json:"access"`
	HasAccess   bool      `json:"hasAccess"`
	IsPending   bool      `json:"isPending"`
	MemberCount int       `json:"memberCount"`
	Avatar      Avatar    `json:"avatar"`
	Projects    []Project `json:"projects,omitempty"` // Optional - some teams may not have projects
}

type Project struct {
	ID                         string     `json:"id"`
	Slug                       string     `json:"slug"`
	Name                       string     `json:"name"`
	Platform                   *string    `json:"platform"`
	DateCreated                time.Time  `json:"dateCreated"`
	IsBookmarked               bool       `json:"isBookmarked"`
	IsMember                   bool       `json:"isMember"`
	Features                   []string   `json:"features"`
	FirstEvent                 *time.Time `json:"firstEvent"`
	FirstTransactionEvent      bool       `json:"firstTransactionEvent"`
	Access                     []string   `json:"access"`
	HasAccess                  bool       `json:"hasAccess"`
	HasMinifiedStackTrace      bool       `json:"hasMinifiedStackTrace"`
	HasMonitors                bool       `json:"hasMonitors"`
	HasProfiles                bool       `json:"hasProfiles"`
	HasReplays                 bool       `json:"hasReplays"`
	HasFlags                   bool       `json:"hasFlags"`
	HasFeedbacks               bool       `json:"hasFeedbacks"`
	HasNewFeedbacks            bool       `json:"hasNewFeedbacks"`
	HasSessions                bool       `json:"hasSessions"`
	HasInsightsHttp            bool       `json:"hasInsightsHttp"`
	HasInsightsDb              bool       `json:"hasInsightsDb"`
	HasInsightsAssets          bool       `json:"hasInsightsAssets"`
	HasInsightsAppStart        bool       `json:"hasInsightsAppStart"`
	HasInsightsScreenLoad      bool       `json:"hasInsightsScreenLoad"`
	HasInsightsVitals          bool       `json:"hasInsightsVitals"`
	HasInsightsCaches          bool       `json:"hasInsightsCaches"`
	HasInsightsQueues          bool       `json:"hasInsightsQueues"`
	HasInsightsLlmMonitoring   bool       `json:"hasInsightsLlmMonitoring"`
	HasInsightsAgentMonitoring bool       `json:"hasInsightsAgentMonitoring"`
	IsInternal                 bool       `json:"isInternal"`
	IsPublic                   bool       `json:"isPublic"`
	Avatar                     Avatar     `json:"avatar"`
	Color                      string     `json:"color"`
	Status                     string     `json:"status"`
}

type TeamFlags struct {
	IDPProvisioned bool `json:"idp:provisioned"`
}
type TeamMember struct {
	ID           string          `json:"id"`
	Email        string          `json:"email"`
	Name         string          `json:"name"`
	User         *TeamMemberUser `json:"user,omitempty"` // Optional - only present for active members
	OrgRole      string          `json:"orgRole"`
	Pending      bool            `json:"pending"`
	Expired      bool            `json:"expired"`
	Flags        TeamMemberFlags `json:"flags"`
	DateCreated  time.Time       `json:"dateCreated"`
	InviteStatus string          `json:"inviteStatus"`
	InviterName  string          `json:"inviterName"`
	TeamRole     string          `json:"teamRole"`
	TeamSlug     string          `json:"teamSlug"`
}

type TeamMemberUser struct {
	ID              string                 `json:"id"`
	Name            string                 `json:"name"`
	Username        string                 `json:"username"`
	Email           string                 `json:"email"`
	AvatarURL       string                 `json:"avatarUrl"`
	IsActive        bool                   `json:"isActive"`
	HasPasswordAuth bool                   `json:"hasPasswordAuth"`
	IsManaged       bool                   `json:"isManaged"`
	DateJoined      time.Time              `json:"dateJoined"`
	LastLogin       *time.Time             `json:"lastLogin"`
	Has2FA          bool                   `json:"has2fa"`
	LastActive      *time.Time             `json:"lastActive"`
	IsSuperuser     bool                   `json:"isSuperuser"`
	IsStaff         bool                   `json:"isStaff"`
	Experiments     map[string]interface{} `json:"experiments"`
	Emails          []TeamMemberEmail      `json:"emails"`
	Avatar          Avatar                 `json:"avatar"`
	CanReset2FA     bool                   `json:"canReset2fa"`
}

type TeamMemberEmail struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
}

type TeamMemberFlags struct {
	IDPProvisioned        bool `json:"idp:provisioned"`
	IDPRoleRestricted     bool `json:"idp:role-restricted"`
	SSOLinked             bool `json:"sso:linked"`
	SSOInvalid            bool `json:"sso:invalid"`
	MemberLimitRestricted bool `json:"member-limit:restricted"`
	PartnershipRestricted bool `json:"partnership:restricted"`
}

type ProjectMember OrganizationMember

// DetailedMember represents a comprehensive member with full role information
type DetailedMember struct {
	ID           string              `json:"id"`
	Email        string              `json:"email"`
	Name         string              `json:"name"`
	User         *DetailedMemberUser `json:"user,omitempty"`
	Role         string              `json:"role"`
	OrgRole      string              `json:"orgRole"`
	RoleName     string              `json:"roleName"`
	Pending      bool                `json:"pending"`
	Expired      bool                `json:"expired"`
	Flags        DetailedMemberFlags `json:"flags"`
	DateCreated  time.Time           `json:"dateCreated"`
	InviteStatus string              `json:"inviteStatus"`
	InviterName  string              `json:"inviterName"`
	Teams        []string            `json:"teams"`
	TeamRoles    []MemberTeamRole    `json:"teamRoles"`
	InviteLink   *string             `json:"invite_link"`
	IsOnlyOwner  bool                `json:"isOnlyOwner"`
	OrgRoleList  []OrganizationRole  `json:"orgRoleList"`
	TeamRoleList []TeamRole          `json:"teamRoleList"`
}

type DetailedMemberUser struct {
	ID              string                `json:"id"`
	Name            string                `json:"name"`
	Username        string                `json:"username"`
	Email           string                `json:"email"`
	AvatarURL       string                `json:"avatarUrl"`
	IsActive        bool                  `json:"isActive"`
	HasPasswordAuth bool                  `json:"hasPasswordAuth"`
	IsManaged       bool                  `json:"isManaged"`
	DateJoined      time.Time             `json:"dateJoined"`
	LastLogin       *time.Time            `json:"lastLogin"`
	Has2FA          bool                  `json:"has2fa"`
	LastActive      *time.Time            `json:"lastActive"`
	IsSuperuser     bool                  `json:"isSuperuser"`
	IsStaff         bool                  `json:"isStaff"`
	Experiments     map[string]any        `json:"experiments"`
	Emails          []DetailedMemberEmail `json:"emails"`
	Avatar          DetailedMemberAvatar  `json:"avatar"`
	Authenticators  []any                 `json:"authenticators"`
	CanReset2FA     bool                  `json:"canReset2fa"`
}

type DetailedMemberEmail struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
}

type DetailedMemberAvatar struct {
	AvatarType string  `json:"avatarType"`
	AvatarUUID *string `json:"avatarUuid"`
}

type DetailedMemberFlags struct {
	IDPProvisioned        bool `json:"idp:provisioned"`
	IDPRoleRestricted     bool `json:"idp:role-restricted"`
	SSOLinked             bool `json:"sso:linked"`
	SSOInvalid            bool `json:"sso:invalid"`
	MemberLimitRestricted bool `json:"member-limit:restricted"`
	PartnershipRestricted bool `json:"partnership:restricted"`
}

type MemberTeamRole struct {
	TeamSlug string `json:"teamSlug"`
	Role     string `json:"role"`
}

type OrganizationRole struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Desc               string   `json:"desc"`
	Scopes             []string `json:"scopes"`
	Allowed            bool     `json:"allowed"`
	IsAllowed          bool     `json:"isAllowed"`
	IsRetired          bool     `json:"isRetired"`
	IsGlobal           bool     `json:"is_global"`
	IsGlobalAlt        bool     `json:"isGlobal"`
	IsTeamRolesAllowed bool     `json:"isTeamRolesAllowed"`
	MinimumTeamRole    string   `json:"minimumTeamRole"`
}

type TeamRole struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Desc               string   `json:"desc"`
	Scopes             []string `json:"scopes"`
	Allowed            bool     `json:"allowed"`
	IsAllowed          bool     `json:"isAllowed"`
	IsRetired          bool     `json:"isRetired"`
	IsTeamRolesAllowed bool     `json:"isTeamRolesAllowed"`
	IsMinimumRoleFor   *string  `json:"isMinimumRoleFor"`
}
type DetailedProject struct {
	ID                               string                 `json:"id"`
	Slug                             string                 `json:"slug"`
	Name                             string                 `json:"name"`
	Platform                         string                 `json:"platform"`
	DateCreated                      string                 `json:"dateCreated"`
	IsBookmarked                     bool                   `json:"isBookmarked"`
	IsMember                         bool                   `json:"isMember"`
	Features                         []string               `json:"features"`
	FirstEvent                       string                 `json:"firstEvent"`
	FirstTransactionEvent            bool                   `json:"firstTransactionEvent"`
	Access                           []string               `json:"access"`
	HasAccess                        bool                   `json:"hasAccess"`
	HasMinifiedStackTrace            bool                   `json:"hasMinifiedStackTrace"`
	HasFeedbacks                     bool                   `json:"hasFeedbacks"`
	HasMonitors                      bool                   `json:"hasMonitors"`
	HasNewFeedbacks                  bool                   `json:"hasNewFeedbacks"`
	HasProfiles                      bool                   `json:"hasProfiles"`
	HasReplays                       bool                   `json:"hasReplays"`
	HasFlags                         bool                   `json:"hasFlags"`
	HasSessions                      bool                   `json:"hasSessions"`
	HasInsightsHttp                  bool                   `json:"hasInsightsHttp"`
	HasInsightsDb                    bool                   `json:"hasInsightsDb"`
	HasInsightsAssets                bool                   `json:"hasInsightsAssets"`
	HasInsightsAppStart              bool                   `json:"hasInsightsAppStart"`
	HasInsightsScreenLoad            bool                   `json:"hasInsightsScreenLoad"`
	HasInsightsVitals                bool                   `json:"hasInsightsVitals"`
	HasInsightsCaches                bool                   `json:"hasInsightsCaches"`
	HasInsightsQueues                bool                   `json:"hasInsightsQueues"`
	HasInsightsLlmMonitoring         bool                   `json:"hasInsightsLlmMonitoring"`
	HasInsightsAgentMonitoring       bool                   `json:"hasInsightsAgentMonitoring"`
	IsInternal                       bool                   `json:"isInternal"`
	IsPublic                         bool                   `json:"isPublic"`
	Avatar                           Avatar                 `json:"avatar"`
	Color                            string                 `json:"color"`
	Status                           string                 `json:"status"`
	Team                             Team                   `json:"team"`
	Teams                            []ProjectTeam          `json:"teams"`
	LatestRelease                    Release                `json:"latestRelease"`
	Options                          map[string]interface{} `json:"options"`
	DigestsMinDelay                  int                    `json:"digestsMinDelay"`
	DigestsMaxDelay                  int                    `json:"digestsMaxDelay"`
	SubjectPrefix                    string                 `json:"subjectPrefix"`
	AllowedDomains                   []string               `json:"allowedDomains"`
	ResolveAge                       int                    `json:"resolveAge"`
	DataScrubber                     bool                   `json:"dataScrubber"`
	DataScrubberDefaults             bool                   `json:"dataScrubberDefaults"`
	SafeFields                       []string               `json:"safeFields"`
	StoreCrashReports                int                    `json:"storeCrashReports"`
	SensitiveFields                  []string               `json:"sensitiveFields"`
	SubjectTemplate                  string                 `json:"subjectTemplate"`
	SecurityToken                    string                 `json:"securityToken"`
	SecurityTokenHeader              *string                `json:"securityTokenHeader"`
	VerifySSL                        bool                   `json:"verifySSL"`
	ScrubIPAddresses                 bool                   `json:"scrubIPAddresses"`
	ScrapeJavaScript                 bool                   `json:"scrapeJavaScript"`
	GroupingConfig                   string                 `json:"groupingConfig"`
	GroupingEnhancements             string                 `json:"groupingEnhancements"`
	GroupingEnhancementsBase         *string                `json:"groupingEnhancementsBase"`
	DerivedGroupingEnhancements      string                 `json:"derivedGroupingEnhancements"`
	SecondaryGroupingExpiry          int                    `json:"secondaryGroupingExpiry"`
	SecondaryGroupingConfig          string                 `json:"secondaryGroupingConfig"`
	FingerprintingRules              string                 `json:"fingerprintingRules"`
	Organization                     Organization           `json:"organization"`
	Plugins                          []Plugin               `json:"plugins"`
	Platforms                        []string               `json:"platforms"`
	ProcessingIssues                 int                    `json:"processingIssues"`
	DefaultEnvironment               string                 `json:"defaultEnvironment"`
	RelayPiiConfig                   *string                `json:"relayPiiConfig"`
	BuiltinSymbolSources             []string               `json:"builtinSymbolSources"`
	DynamicSamplingBiases            []DynamicSamplingBias  `json:"dynamicSamplingBiases"`
	DynamicSamplingMinimumSampleRate bool                   `json:"dynamicSamplingMinimumSampleRate"`
	EventProcessing                  EventProcessing        `json:"eventProcessing"`
	SymbolSources                    string                 `json:"symbolSources"`
	TempestFetchScreenshots          bool                   `json:"tempestFetchScreenshots"`
	TempestFetchDumps                bool                   `json:"tempestFetchDumps"`
	IsDynamicallySampled             bool                   `json:"isDynamicallySampled"`
	AutofixAutomationTuning          string                 `json:"autofixAutomationTuning"`
	SeerScannerAutomation            bool                   `json:"seerScannerAutomation"`
	HighlightTags                    []string               `json:"highlightTags"`
	HighlightContext                 map[string]interface{} `json:"highlightContext"`
	HighlightPreset                  HighlightPreset        `json:"highlightPreset"`
}

type ProjectTeam struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Release struct {
	Version string `json:"version"`
}

type Status struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Links struct {
	OrganizationUrl string `json:"organizationUrl"`
	RegionUrl       string `json:"regionUrl"`
}

type Plugin struct {
	ID                    string                 `json:"id"`
	Name                  string                 `json:"name"`
	Slug                  string                 `json:"slug"`
	ShortName             string                 `json:"shortName"`
	Type                  string                 `json:"type"`
	CanDisable            bool                   `json:"canDisable"`
	IsTestable            bool                   `json:"isTestable"`
	HasConfiguration      bool                   `json:"hasConfiguration"`
	Metadata              map[string]interface{} `json:"metadata"`
	Contexts              []string               `json:"contexts"`
	Status                string                 `json:"status"`
	Assets                []string               `json:"assets"`
	Doc                   string                 `json:"doc"`
	FirstPartyAlternative *string                `json:"firstPartyAlternative"`
	DeprecationDate       *string                `json:"deprecationDate"`
	AltIsSentryApp        *string                `json:"altIsSentryApp"`
	Enabled               bool                   `json:"enabled"`
	Version               string                 `json:"version"`
	Author                Author                 `json:"author"`
	IsDeprecated          bool                   `json:"isDeprecated"`
	IsHidden              bool                   `json:"isHidden"`
	Description           string                 `json:"description"`
	Features              []string               `json:"features"`
	FeatureDescriptions   []FeatureDescription   `json:"featureDescriptions"`
	ResourceLinks         []ResourceLink         `json:"resourceLinks"`
}

type Author struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type FeatureDescription struct {
	Description string `json:"description"`
	FeatureGate string `json:"featureGate"`
}

type ResourceLink struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type DynamicSamplingBias struct {
	ID     string `json:"id"`
	Active bool   `json:"active"`
}

type EventProcessing struct {
	SymbolicationDegraded bool `json:"symbolicationDegraded"`
}

type HighlightPreset struct {
	Tags    []string               `json:"tags"`
	Context map[string]interface{} `json:"context"`
}

type AddOrganizationMemberBody struct {
	Email string `json:"email"`
	// Optional.
	//  Possible values are:
	//, "owner", "manager", "member", "billing"
	OrgRole string `json:"orgRole,omitempty"`
}
