package model

import (
	"errors"
	"fmt"
	"time"

	pb "github.com/skeleton1231/gotal/internal/proto/user"
	"github.com/skeleton1231/gotal/pkg/log"
	"github.com/skeleton1231/gotal/pkg/util/common"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type User struct {
	ObjectMeta
	Name string `json:"name,omitempty" gorm:"column:name;type:varchar(255);not null" validate:"required"`
	// Required: true
	Email           string    `json:"email" gorm:"column:email" validate:"required,email,min=1,max=100"`
	EmailVerifiedAt time.Time `gorm:"column:email_verified_at" json:"-"`
	Password        string    `json:"-" gorm:"column:password" validate:"required"`
	RememberToken   string    `gorm:"size:100" json:"-"`
	StripeID        string    `gorm:"size:255" json:"stripeId"`
	DiscordID       uint64    `gorm:"default:0" json:"discordId"`
	PMType          string    `gorm:"size:255" json:"-"`
	PMLastFour      string    `gorm:"size:4" json:"-"`
	TrialEndsAt     time.Time `gorm:"column:trial_ends_at" json:"-"`
	TotalCredits    int       `gorm:"default:0" json:"totalCredits"`
	Token           string    `json:"token,omitempty" gorm:"-"`
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

// ToProto converts User model to protobuf message
func UserToProto(u *User) *pb.User {

	// 使用u.Extend构造*structpb.Struct
	extendProto, err := structpb.NewStruct(u.Extend)
	if err != nil {
		// 处理错误，例如打印日志或者返回一个错误
		// 注意：在生产代码中，您应该处理这个错误而不是忽略它
		log.Errorf("Error converting extend to proto struct: %v", err)
	}

	// 创建并返回pb.User
	return &pb.User{
		Meta: &pb.ObjectMeta{
			Id:           u.ID,
			Extend:       extendProto, // 使用转换后的extend
			ExtendShadow: u.ExtendShadow,
			CreatedAt:    timestamppb.New(u.CreatedAt),
			UpdatedAt:    timestamppb.New(u.UpdatedAt),
			// 注意处理DeletedAt字段
			IsDeleted: u.Status == 1, // 示例，根据您的业务逻辑调整
			Status:    int32(u.Status),
		},
		Name:            u.Name,
		Email:           u.Email,
		EmailVerifiedAt: timestamppb.New(u.EmailVerifiedAt),
		RememberToken:   u.RememberToken,
		StripeId:        u.StripeID,
		DiscordId:       u.DiscordID,
		PmType:          u.PMType,
		PmLastFour:      u.PMLastFour,
		TrialEndsAt:     timestamppb.New(u.TrialEndsAt),
		TotalCredits:    int32(u.TotalCredits),
		// Token字段通常用于认证响应，而不是作为用户模型的一部分发送，这里也不包括它
	}
}

func ProtoToUser(pbUser *pb.User) (*User, error) {

	if pbUser == nil {
		return nil, errors.New("userProto is nil")
	}
	// 转换extend字段
	var extend Extend
	if pbUser.Meta.GetExtend() != nil {
		extendMap, err := structpbToMap(pbUser.Meta.GetExtend())
		if err != nil {
			return nil, err // 如果转换失败，返回错误
		}
		extend = extendMap
	}

	user := &User{
		ObjectMeta: ObjectMeta{
			ID:           pbUser.Meta.GetId(),
			Extend:       extend,
			ExtendShadow: pbUser.Meta.GetExtendShadow(),
			CreatedAt:    pbUser.Meta.GetCreatedAt().AsTime(),
			UpdatedAt:    pbUser.Meta.GetUpdatedAt().AsTime(),
			Status:       int(pbUser.Meta.GetStatus()),
		},
		Name:            pbUser.GetName(),
		Email:           pbUser.GetEmail(),
		EmailVerifiedAt: pbUser.GetEmailVerifiedAt().AsTime(),
		RememberToken:   pbUser.GetRememberToken(),
		StripeID:        pbUser.GetStripeId(),
		DiscordID:       pbUser.GetDiscordId(),
		PMType:          pbUser.GetPmType(),
		PMLastFour:      pbUser.GetPmLastFour(),
		TrialEndsAt:     pbUser.GetTrialEndsAt().AsTime(),
		TotalCredits:    int(pbUser.GetTotalCredits()),
	}

	return user, nil
}

// structpbToMap 是一个辅助函数，用于将*structpb.Struct转换为map[string]interface{}
func structpbToMap(structProto *structpb.Struct) (Extend, error) {
	extend := make(Extend)
	if structProto == nil {
		return extend, nil
	}

	for k, v := range structProto.Fields {
		extend[k] = v.AsInterface()
	}

	return extend, nil
}
