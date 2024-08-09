package data

import (
	"context"
	v1 "review-o/api/review/v1"
	"review-o/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type operationRepo struct {
	data *Data
	log  *log.Helper
}

// NewOperationRepo .
func NewOperationRepo(data *Data, logger log.Logger) biz.OperationRepo {
	return &operationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *operationRepo) AuditAppeal(ctx context.Context, param *biz.AuditAppealParam) (int64, error) {
	r.log.WithContext(ctx).Infof("[data]Audit,param:%v\n", param)
	reply, err := r.data.rc.AuditAppeal(ctx, &v1.AuditAppealRequest{
		AppealID:  param.AppealID,
		Status:    param.Status,
		OpUser:    param.OpUser,
		OpReason:  param.OpReason,
		OpRemarks: param.OpRemarks,
	})
	r.log.WithContext(ctx).Debugf("AduitAppeal return, ret:%v err:%v\n", reply, err)
	if err != nil {
		return 0, err
	}
	return param.AppealID, nil
}

func (r *operationRepo) AuditReview(ctx context.Context, param *biz.AuditReviewParam) (int64, error) {
	r.log.WithContext(ctx).Infof("[data]AuditReview,param:%v\n", param)
	reply, err := r.data.rc.AuditReview(ctx, &v1.AuditReviewRequest{
		ReviewID:  param.ReviewID,
		Status:    param.Status,
		OpUser:    param.OpUser,
		OpReason:  param.OpReason,
		OpRemarks: param.OpRemarks,
	})
	r.log.WithContext(ctx).Debugf("AduitReview return, ret:%v err:%v\n", reply, err)
	if err != nil {
		return 0, err
	}
	return param.ReviewID, nil
}
