package biz

import "time"

type GetAllReviewsParam struct {
	UserID int64
}

type ReplyParam struct {
	ReviewID  int64
	StoreID   int64
	Content   string
	PicInfo   string
	VideoInfo string
}

type ReplyUpdateParam struct {
	ReviewID   int64
	StoreID    int64
	ReplyID    int64
	UpdateTime time.Time
	Content    string
	PicInfo    string
	VideoInfo  string
}

type AuditAppealParam struct {
	AppealID  int64
	OpUser    string
	OpReason  string
	OpRemarks string
	Status    int32
}

type AuditReviewParam struct {
	ReviewID  int64
	OpUser    string
	OpReason  string
	OpRemarks string
	Status    int32
}

type AppealParam struct {
	ReviewID  int64
	StoreID   int64
	AppealID  int64
	Reason    string
	Content   string
	PicInfo   string
	VideoInfo string
}
