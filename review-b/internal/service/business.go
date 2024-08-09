package service

import (
	"context"
	"fmt"
	"review-b/internal/biz"

	pb "review-b/api/business/v1"
)

type BusinessService struct {
	pb.UnimplementedBusinessServer
	uc *biz.BusinessUsecase
}

func NewBusinessService(uc *biz.BusinessUsecase) *BusinessService {
	return &BusinessService{uc: uc}
}

func (s *BusinessService) ReplyReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	replyid, err := s.uc.CreateReply(ctx, &biz.ReplyParam{
		ReviewID:  req.ReviewID,
		StoreID:   req.StoreID,
		Content:   req.Content,
		PicInfo:   req.PicInfo,
		VideoInfo: req.VideoInfo,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ReplyReviewReply{ReplyID: replyid}, nil
}

func (s *BusinessService) AppealReview(ctx context.Context, req *pb.AppealReviewRequest) (*pb.AppealReviewReply, error) {
	fmt.Printf("[service] AppealReview, req:#%v\n", req)
	appealID, err := s.uc.AppealReview(ctx, &biz.AppealParam{
		ReviewID:  req.ReviewID,
		StoreID:   req.StoreID,
		Reason:    req.Reason,
		Content:   req.Content,
		PicInfo:   req.PicInfo,
		VideoInfo: req.VideoInfo,
	})
	if err != nil {
		return nil, err
	}
	return &pb.AppealReviewReply{AppealID: appealID}, nil
}

func (s *BusinessService) ReplyReviewUpdate(ctx context.Context, req *pb.ReplyReviewUpdateRequest) (*pb.ReplyReviewUpdateReply, error) {
	fmt.Printf("[service] UpdateReview, req:#%v\n", req)
	rowaffectedd, err := s.uc.UpdateReply(ctx, &biz.UpdateParam{
		ReviewID:  req.ReviewID,
		StoreID:   req.StoreID,
		ReplyID:   req.ReplyID,
		Content:   req.Content,
		PicInfo:   req.PicInfo,
		VideoInfo: req.VideoInfo,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ReplyReviewUpdateReply{ReplyID: req.ReplyID, Rowsaffected: rowaffectedd}, nil
}
