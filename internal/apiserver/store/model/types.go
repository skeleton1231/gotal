package model

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type Extend map[string]interface{}

func (ext Extend) String() string {
	data, _ := json.Marshal(ext)
	return string(data)
}

func (ext Extend) Merge(extendShadow string) Extend {
	var extend Extend

	_ = json.Unmarshal([]byte(extendShadow), &extend)
	for k, v := range extend {
		if _, ok := ext[k]; !ok {
			ext[k] = v
		}
	}

	return ext
}

type ListMeta struct {
	TotalCount int64 `json:"totalCount,omitempty"`
}

type ObjectMeta struct {
	ID uint64 `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT;column:id"`

	Extend Extend `json:"extend,omitempty" gorm:"-" validate:"omitempty"`

	ExtendShadow string `json:"-" gorm:"column:extendShadow" validate:"omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"column:created_at"`

	UpdatedAt time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at"`

	DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deleted_at;index:idx_deleted_at"`
}

func (obj *ObjectMeta) BeforeCreate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

func (obj *ObjectMeta) BeforeUpdate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

func (obj *ObjectMeta) AfterFind(tx *gorm.DB) error {
	if err := json.Unmarshal([]byte(obj.ExtendShadow), &obj.Extend); err != nil {
		return err
	}

	return nil
}

type ListOptions struct {
	LabelSelector string `json:"labelSelector,omitempty" form:"labelSelector"`

	FieldSelector string `json:"fieldSelector,omitempty" form:"fieldSelector"`

	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`

	Offset *int64 `json:"offset,omitempty" form:"offset"`

	Limit *int64 `json:"limit,omitempty" form:"limit"`
}

type GetOptions struct {
}

type DeleteOptions struct {
	Unscoped bool `json:"unscoped"`
}

type CreateOptions struct {
	DryRun []string `json:"dryRun,omitempty"`
}

type PatchOptions struct {
	DryRun []string `json:"dryRun,omitempty"`

	Force bool `json:"force,omitempty"`
}

type UpdateOptions struct {
	DryRun []string `json:"dryRun,omitempty"`
}

type TableOptions struct {
	NoHeaders bool `json:"-"`
}

const DefaultLimit = 1000

type LimitAndOffset struct {
	Offset int
	Limit  int
}

func Unpointer(offset *int64, limit *int64) *LimitAndOffset {
	var o, l int = 0, DefaultLimit

	if offset != nil {
		o = int(*offset)
	}

	if limit != nil {
		l = int(*limit)
	}

	return &LimitAndOffset{
		Offset: o,
		Limit:  l,
	}
}
