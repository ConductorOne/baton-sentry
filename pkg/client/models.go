package client

import (
	"time"
)

type Organization struct {
	Avatar                     OrganizationAvatar `json:"avatar"`
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

type OrganizationAvatar struct {
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
	Emails          []UserEmail            `json:"emails"`
	Avatar          UserAvatar             `json:"avatar"`
	CanReset2FA     bool                   `json:"canReset2fa"`
}

type UserEmail struct {
	ID         string `json:"id"`
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
}

type UserAvatar struct {
	AvatarType string  `json:"avatarType"`
	AvatarUUID *string `json:"avatarUuid"`
}

type MemberFlags struct {
	IDPProvisioned        bool `json:"idp:provisioned"`
	IDPRoleRestricted     bool `json:"idp:role-restricted"`
	SSOLinked             bool `json:"sso:linked"`
	SSOInvalid            bool `json:"sso:invalid"`
	MemberLimitRestricted bool `json:"member-limit:restricted"`
	PartnershipRestricted bool `json:"partnership:restricted"`
}
