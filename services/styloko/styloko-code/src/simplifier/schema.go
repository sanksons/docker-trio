package simplifier

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type UpdateJob struct {
	Id          bson.ObjectId `bson:"_id" json:"id"`
	Controller  string        `bson:"controller" json:"controller"`
	Api         string        `bson:"api" json:"api"`
	Type        string        `bson:"type" json:"type"`
	JobName     string        `bson:"job_name" json:"job_name"`
	UserId      int           `bson:"user_id" json:"user_id"`
	Data        []interface{} `bson:"data" json:"data"`
	CreatedAt   time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at" json:"updated_at"`
	IsPicked    bool          `bson:"is_picked" json:"is_picked"`
	IsCompleted bool          `bson:"is_completed" json:"is_completed"`
	TotalChunks int           `bson:"total_chunks" json:"total_chunks"`
}

type ErrorStruct struct {
	JobName string `json:"job_name"`
	ErrMsg  string `json:"message"`
}

type FlorestResp struct {
	Status   interface{}   `json:"status"`
	Data     []interface{} `json:"data"`
	MetaData interface{}   `json:"_metaData"`
}
