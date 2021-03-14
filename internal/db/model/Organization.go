package model

import (
	"database/sql"
	"strconv"
	"time"

	v1 "github.com/ekspand/trusty/api/v1"
)

// Organization represents an organization account.
type Organization struct {
	ID         int64         `db:"id"`
	ExternalID sql.NullInt64 `db:"extern_id"`
	Provider   string        `db:"provider"`
	Login      string        `db:"login"`
	AvatarURL  string        `db:"avatar_url"`
	Name       string        `db:"name"`
	Email      string        `db:"email"`
	Company    string        `db:"company"`
	CreatedAt  time.Time     `db:"created_at"`
	UpdatedAt  time.Time     `db:"updated_at"`
	Type       string        `db:"type"`
}

// ToDto converts model to v1.Organization DTO
func (u *Organization) ToDto() *v1.Organization {
	user := &v1.Organization{
		ID:        strconv.FormatUint(uint64(u.ID), 10),
		Provider:  u.Provider,
		Login:     u.Login,
		Name:      u.Name,
		Email:     u.Email,
		Company:   u.Company,
		AvatarURL: u.AvatarURL,
		Type:      u.Type,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}

	if u.ExternalID.Valid {
		user.ExternalID = strconv.FormatUint(uint64(u.ExternalID.Int64), 10)
	}

	return user
}
