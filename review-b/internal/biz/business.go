package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type ReplyParam struct {
	ReviewID  int64
	StoreID   int64
	Content   string
	PicInfo   string
	VideoInfo string
}

type AppealParam struct {
	ReviewID  int64
	StoreID   int64
	Reason    string
	Content   string
	PicInfo   string
	VideoInfo string
}

type UpdateParam struct {
	ReviewID  int64
	StoreID   int64
	ReplyID   int64
	Content   string
	PicInfo   string
	VideoInfo string
}

// GreeterRepo is a Greater repo.
type BusinessRepo interface {
	Reply(context.Context, *ReplyParam) (int64, error)
	Appeal(ctx context.Context, param *AppealParam) (int64, error)
	Update(ctx context.Context, param *UpdateParam) (int64, error)
}

// GreeterUsecase is a Greeter usecase.
type BusinessUsecase struct {
	repo BusinessRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewBusinessUsecase(repo BusinessRepo, logger log.Logger) *BusinessUsecase {
	return &BusinessUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateReply创建回复，service层调用这个方法
func (uc *BusinessUsecase) CreateReply(ctx context.Context, param *ReplyParam) (int64, error) {
	uc.log.WithContext(ctx).Infof("[biz] CreateReply param:%v\n", param)
	return uc.repo.Reply(ctx, param)
}

// AppealReview申诉评价，service层调用这个方法
func (uc *BusinessUsecase) AppealReview(ctx context.Context, param *AppealParam) (int64, error) {
	uc.log.WithContext(ctx).Debugf("[biz] AppealReview, reviewID:%v\n", param.ReviewID)
	return uc.repo.Appeal(ctx, param)
}

func (uc *BusinessUsecase) UpdateReply(ctx context.Context, param *UpdateParam) (int64, error) {
	uc.log.WithContext(ctx).Debugf("[biz] UpdateReply, replyID:%v\n", param.ReviewID)
	return uc.repo.Update(ctx, param)
}
