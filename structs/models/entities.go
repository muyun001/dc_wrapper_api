package models

import "time"

type Request struct {
	Url       string `json:"url" bson:"url"`
	UserAgent string `json:"user_agent" bson:"user_agent"`
	Cookie    string `json:"cookie" bson:"cookie"`
	Body      string `json:"body" bson:"body"`
}

type Config struct {
	District       string   `json:"district" bson:"district"`
	ResponseTypes  []string `json:"response_types" bson:"response_types"`
	FollowRedirect bool     `json:"follow_redirect" bson:"follow_redirect"`
	Priority       string   `json:"priority" bson:"priority"`
}

type DcRequest struct {
	UniqueKey string    `json:"unique_key" bson:"unique_key"`
	UniqueMD5 string    `json:"unique_md5" bson:"unique_md5"`
	Request   Request   `json:"request" bson:"request"`
	Config    Config    `json:"config" bson:"config"`
	Status    int       `json:"status" bson:"status"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type DcResponse struct {
	Header  string `json:"header"`
	Body    string `json:"body"`
	Capture string `json:"capture"`
}
