package service

import (
	"context"
	"fmt"
	"review-o/internal/biz"

	pb "review-o/api/operation/v1"
)

type OperationService struct {
	pb.UnimplementedOperationServer
	uc *biz.OperationUsecase
}

func NewOperationService(uc *biz.OperationUsecase) *OperationService {
	return &OperationService{uc: uc}
}

func (s *OperationService) AuditAppeal(ctx context.Context, req *pb.AuditAppealRequest) (*pb.AuditAppealReply, error) {
	fmt.Printf("[service] AuditAppeal, req:#%v\n", req)
	_, err := s.uc.Audit(ctx, &biz.AuditAppealParam{
		AppealID:  req.AppealID,
		OpUser:    req.OpUser,
		OpReason:  req.OpReason,
		OpRemarks: req.OpRemarks,
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AuditAppealReply{}, nil
}
func (s *OperationService) AuditReview(ctx context.Context, req *pb.AuditReviewRequest) (*pb.AuditReviewReply, error) {
	fmt.Printf("[service] AuditReview, req:#%v\n", req)
	_, err := s.uc.AuditReview(ctx, &biz.AuditReviewParam{
		ReviewID:  req.ReviewID,
		OpUser:    req.OpUser,
		OpReason:  req.OpReason,
		OpRemarks: req.OpRemarks,
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AuditReviewReply{}, nil
}
