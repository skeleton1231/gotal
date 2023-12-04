package model

import (
	"fmt"
	"time"

	"github.com/skeleton1231/gotal/pkg/util/common"
)

type User struct {
	ObjectMeta
	Name  string `gorm:"size:255;not null" json:"name"`
	Email string `gorm:"size:255" json:"email"`
	// 	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"column:created_at"`
	EmailVerifiedAt time.Time `gorm:"column:email_verified_at" json:"-"`
	Password        string    `gorm:"size:255;not null" json:"-"`
	RememberToken   string    `gorm:"size:100" json:"-"`
	StripeID        string    `gorm:"size:255" json:"stripeId"`
	DiscordID       uint64    `gorm:"default:0" json:"discordId"`
	PMType          string    `gorm:"size:255" json:"-"`
	PMLastFour      string    `gorm:"size:4" json:"-"`
	TrialEndsAt     time.Time `gorm:"column:trial_ends_at" json:"-"`
	TotalCredits    int       `gorm:"default:0" json:"totalCredits"`
}

// TableName overrides the table name used by User to `users`.
func (User) TableName() string {
	return "users"
}

// UserList is the whole list of all users which have been stored in stroage.
type UserList struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard list metadata.
	// +optional
	ListMeta `json:",inline"`

	Items []*User `json:"items"`
}

func (u *User) Compare(pwd string) error {
	if err := common.Compare(u.Password, pwd); err != nil {
		return fmt.Errorf("failed to compile password: %w", err)
	}

	return nil
}
