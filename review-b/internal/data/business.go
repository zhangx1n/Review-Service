package data

import (
	"context"
	v1 "review-b/api/review/v1"

	"review-b/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type businessRepo struct {
	data *Data
	log  *log.Helper
}

// NewBusinessRepo .
func NewBusinessRepo(data *Data, logger log.Logger) biz.BusinessRepo {
	return &businessRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *businessRepo) Reply(ctx context.Context, param *biz.ReplyParam) (int64, error) {
	r.log.WithContext(ctx).Infof("[data] Reply,param:%v\n", param)
	reply, err := r.data.rc.ReplyReview(ctx, &v1.ReplyReviewRequest{
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	})
	r.log.WithContext(ctx).Debugf("ReplyReview return, ret:%v err:%v\n", reply, err)
	if err != nil {
		return 0, err
	}
	return reply.ReplyID, nil
}

func (r *businessRepo) Appeal(ctx context.Context, param *biz.AppealParam) (int64, error) {
	r.log.WithContext(ctx).Infof("[data] Appeal,param:%v\n", param)
	reply, err := r.data.rc.AppealReview(ctx, &v1.AppealReviewRequest{
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Reason:    param.Reason,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	})
	r.log.WithContext(ctx).Debugf("AppealReview return, ret:%v err:%v\n", reply, err)
	if err != nil {
		return 0, err
	}
	return reply.AppealID, nil
}

func (r *businessRepo) Update(ctx context.Context, param *biz.UpdateParam) (int64, error) {
	r.log.WithContext(ctx).Infof("[data]Update,param:%v\n", param)
	reply, err := r.data.rc.ReplyReviewUpdate(ctx, &v1.ReplyReviewUpdateRequest{
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		ReplyID:   param.ReplyID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	})
	r.log.WithContext(ctx).Debugf("UpdateReply return, ret:%v err:%v\n", reply, err)
	if err != nil {
		return 0, err
	}
	return reply.Rowsaffected, nil
}
