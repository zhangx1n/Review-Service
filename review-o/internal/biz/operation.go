package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

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

// GreeterRepo is a Greater repo.
type OperationRepo interface {
	AuditAppeal(ctx context.Context, param *AuditAppealParam) (int64, error)
	AuditReview(ctx context.Context, param *AuditReviewParam) (int64, error)
}

// OprationUsecase is a Greeter usecase.
type OperationUsecase struct {
	repo OperationRepo
	log  *log.Helper
}

// NewGreeterUsecase new a Greeter usecase.
func NewOperationUsecase(repo OperationRepo, logger log.Logger) *OperationUsecase {
	return &OperationUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *OperationUsecase) Audit(ctx context.Context, param *AuditAppealParam) (int64, error) {
	uc.log.WithContext(ctx).Debugf("[biz] Audit, replyID:%v\n", param.AppealID)
	return uc.repo.AuditAppeal(ctx, param)
}

func (uc *OperationUsecase) AuditReview(ctx context.Context, param *AuditReviewParam) (int64, error) {
	uc.log.WithContext(ctx).Debugf("[biz] AuditReview, replyID:%v\n", param.ReviewID)
	return uc.repo.AuditReview(ctx, param)
}
