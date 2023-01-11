package profile

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"i69social/src/api/models/fake_user"
	"log"
	"strings"
	"time"
)

const (
	MinPasswordLen          = 6
	MaxPasswordLen          = 255
	MinLoginLen             = 3
	MaxLoginLen             = 50
	MinEmailLen             = 3
	MaxEmailLen             = 255
	RegularUserPrivilege    = "regular"
	ChatterUserPrivilege    = "chatter"
	AdminUserPrivilege      = "admin"
	SuperAdminUserPrivilege = "super_admin"
	InvitationKeySize       = 10
)

type CustomJSONMessage struct {
	json.RawMessage
}

func (t CustomJSONMessage) String() string { return string(t.RawMessage) }

func (t CustomJSONMessage) Value() (driver.Value, error) {
	return []byte(t.RawMessage), nil
}

func (t *CustomJSONMessage) Scan(src interface{}) error {
	var source []byte
	switch orig := src.(type) {
	case []byte:
		source = orig
	default:
		return errors.New("Incompatible type for CustomJSONMessage")
	}

	*t = CustomJSONMessage{RawMessage: json.RawMessage(source)}
	return nil
}

type Profile struct {
	// user identifier (required)
	Id        uint              `json:"id" gorm:"primary_key"`
	Email     string            `json:"email"`
	Password  string            `json:"password,omitempty"`
	FirstName string            `json:"first_name"`
	LastName  string            `json:"last_name"`
	Roles     CustomJSONMessage `json:"roles"`
	IsActive  bool              `json:"is_active"`
	FakeUsers []fake_user.User  `json:"fake_users" gorm:"foreignkey:OwnerId"`
	CreatedAt time.Time         `json:"-"`
	UpdatedAt time.Time         `json:"-"`
	DeletedAt *time.Time        `json:"-"`
}

func (p *Profile) HasRole(role string) bool {
	var roles []string
	err := json.Unmarshal(p.Roles.RawMessage, &roles)
	if err != nil {
		log.Printf("roles %s unmarshal err: %s", string(p.Roles.RawMessage), err)
		return false
	}
	for _, r := range roles {
		if r == role {
			return true
		}
	}

	return false
}

func (p *Profile) GetPrivileges() map[string]bool {
	var roles []string
	err := json.Unmarshal(p.Roles.RawMessage, &roles)
	if err != nil {
		return map[string]bool{}
	}

	var privileges = make(map[string]bool)
	for _, role := range roles {
		privileges[role] = true
	}
	privileges[RegularUserPrivilege] = true

	return privileges
}

func (Profile) TableName() string {
	return "users"
}

type GenerateInvitationRequest struct {
	// (required)
	Email                 string `json:"email" binding:"required" validate:"max=255,min=5"`
	IsAdminPermission     bool   `json:"is_admin_permission"`
	IsChatAdminPermission bool   `json:"is_chat_admin_permission"`
}

func (r *GenerateInvitationRequest) TrimSpaces() {
	r.Email = strings.TrimSpace(r.Email)
}

func (r GenerateInvitationRequest) ToRoles() []string {
	if !r.IsAdminPermission {
		return []string{ChatterUserPrivilege}
	}

	roles := []string{AdminUserPrivilege}
	if r.IsChatAdminPermission {
		roles = append(roles, ChatterUserPrivilege)
	}

	return roles
}

type SignUpRequest struct {
	InvitationKey string `json:"invitation_key" binding:"required" validate:"max=255,min=5"`
	// (required)
	Password  string `json:"password" binding:"required" validate:"max=255,min=6"`
	FirstName string `json:"first_name" binding:"required" validate:"max=255,min=6"`
	LastName  string `json:"last_name" binding:"required" validate:"max=255,min=6"`
}

func (u *SignUpRequest) TrimSpaces() {
	u.FirstName = strings.TrimSpace(u.FirstName)
	u.LastName = strings.TrimSpace(u.LastName)
}

func (u *SignUpRequest) Validate() (err error) {
	if len(u.Password) == 0 {
		err = fmt.Errorf("field password cannot be empty")
		return
	}

	if len(u.Password) < MinPasswordLen {
		err = fmt.Errorf("field password is too short. "+
			"given length: `%d`, min length: `%d`", len(u.Password), MinPasswordLen)
		return
	}

	if len(u.Password) > MaxPasswordLen {
		err = fmt.Errorf("field password is too long. "+
			"given length: `%d`, max length: `%d`", len(u.Password), MaxPasswordLen)
		return
	}

	return
}

type SignInRequest struct {
	// (required)
	Email string `json:"email" binding:"required" validate:"max=255,min=5"`
	// (required)
	Password string `json:"password" binding:"required" validate:"max=255,min=6"`
}

func (u *SignInRequest) TrimSpaces() {
	u.Email = strings.TrimSpace(u.Email)
}
