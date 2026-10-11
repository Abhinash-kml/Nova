package users

import "time"

type UserId struct {
	Id string `uri:"id" json:"id" binding:"required,uuid"`
}

type DeleteOptions struct {
	Type string `form:"type" json:"type" binding:"required,oneof=soft hard"`
}

type GetDTO struct {
	UserId
}

type GetAllDTO struct {
	Cursor        string `form:"cursor,default=AAAAAAAAAAAAAAAAAAAAAA" binding:"required"`
	Limit         int    `form:"limit,default=10" binding:"required,gte=10,lte=20"`
	SortDirection string `form:"sort,default=asc" binding:"required,oneof=asc desc"`
	SortBy        string `form:"sortby,default=id" binding:"required"`
	// Filter string `form:"filter,default=id" binding:"required"`
}

type CreateDTO struct {
	Username    string `json:"username" binding:"required,gte=5,lte=20"`
	DisplayName string `json:"display_name" binding:"required,gte=5,lte=20"`
	Email       string `json:"email" binding:"required,email"`
	Country     string `json:"country" binding:"required"`
	State       string `json:"state" binding:"required"`
	LangTag     string `json:"lang_tag" binding:"required"`
	Timezone    string `json:"time_zone" binding:"required"`
	AvatarURL   string `json:"avatar_url" binding:"required"`
}

type FieldUpdate struct {
	Field    string `json:"field" binding:"required"`
	DataType string `json:"datatype" binding:"required"`
	Value    string `json:"value" binding:"required"`
}

type FieldUpdates struct {
	Updates []FieldUpdate `json:"updates" binding:"required"`
}

type UpdateDTO struct {
	UserId
	Username    *string `json:"username"`
	DisplayName *string `json:"display_name"`
	Email       *string `json:"email"`
	Country     *string `json:"country"`
	State       *string `json:"state"`
	LangTag     *string `json:"lang_tag"`
	Timezone    *string `json:"timezone"`
}

type ReplacementData struct {
	Username    string `json:"username" binding:"required,gte=5,lte=20"`
	DisplayName string `json:"displayname" binding:"required,gte=5,lte=10"`
	Email       string `json:"email" binding:"required,email"`
	Country     string `json:"country" binding:"required"`
	State       string `json:"state" binding:"required"`
	LangTag     string `json:"lang_tag" binding:"required"`
	Timezone    string `json:"timezone" binding:"required"`
}

type ReplaceDTO struct {
	UserId
	ReplacementData
}

type DeleteDTO struct {
	UserId
	DeleteOptions
}

type BanDetailsDTO struct {
	By     string    `json:"by"`
	Reason string    `json:"reason"`
	On     time.Time `json:"on"`
	Till   time.Time `json:"till"`
}

type BulkCreateDTO struct {
	Users []CreateDTO `json:"users" binding:"required"`
}

type BulkModifyDTO struct {
	Updates []UpdateDTO `json:"updates" binding:"required"`
}

type BulkDeleteDTO struct {
	Users []DeleteDTO `json:"users" binding:"required"`
}
