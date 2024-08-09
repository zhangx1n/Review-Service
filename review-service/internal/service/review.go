package service

import (
	"context"
	"fmt"
	pb "review-service/api/review/v1"
	"review-service/internal/biz"
	"review-service/internal/data/model"
	"time"
)

type ReviewService struct {
	pb.UnimplementedReviewServer

	uc *biz.ReviewerUsecase
}

func NewReviewService(uc *biz.ReviewerUsecase) *ReviewService {
	return &ReviewService{uc: uc}
}

func (s *ReviewService) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewReply, error) {
	fmt.Printf("[service] CreateReview, req:#%v\n", req)
	//参数转换
	//调用biz层
	var anonymous int32
	if req.Anonymous {
		anonymous = 1
	}
	reviewinfo, err := s.uc.CreateReview(ctx, &model.ReviewInfo{
		UserID:       req.UserID,
		OrderID:      req.OrderID,
		StoreID:      req.StoreID,
		Score:        req.Score,
		ExpressScore: req.ExpressScore,
		ServiceScore: req.ServiceScore,
		Content:      req.Content,
		PicInfo:      req.PicInfo,
		VideoInfo:    req.VideoInfo,
		Anonymous:    anonymous,
		Status:       0,
	})
	if err != nil {
		return nil, err
	}
	return &pb.CreateReviewReply{ReviewID: reviewinfo.ReviewID}, nil
}
func (s *ReviewService) UpdateReview(ctx context.Context, req *pb.UpdateReviewRequest) (*pb.UpdateReviewReply, error) {
	fmt.Printf("[service] UpdateReview, req:#%v\n", req)
	rowsaffected, err := s.uc.UpdateReview(ctx, &model.ReviewInfo{
		ReviewID:     req.ReviewID,
		UserID:       req.UserID,
		Score:        req.Score,
		ServiceScore: req.ServiceScore,
		ExpressScore: req.ExpressScore,
		Content:      req.Content,
		PicInfo:      req.PicInfo,
		VideoInfo:    req.VideoInfo,
	})
	return &pb.UpdateReviewReply{ReviewID: req.ReviewID, Rowaffected: rowsaffected}, err
}
func (s *ReviewService) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.DeleteReviewReply, error) {
	fmt.Printf("[service] DeleteReview, req:#%v\n", req)
	rowsaffected, err := s.uc.DeleteReview(ctx, &model.ReviewInfo{
		UserID:   req.UserID,
		ReviewID: req.ReviewID,
	})
	return &pb.DeleteReviewReply{Rowaffected: rowsaffected}, err
}
func (s *ReviewService) GetReview(ctx context.Context, req *pb.GetReviewRequest) (*pb.GetReviewReply, error) {
	fmt.Printf("[service] GetReview, req:#%v\n", req)
	reviewinfo, err := s.uc.GetReview(ctx, &model.ReviewInfo{ReviewID: req.GetReviewID()})
	return reviewinfo, err
}

func (s *ReviewService) GetAllReviews(ctx context.Context, req *pb.GetAllReviewsRequest) (*pb.GetAllReviewsReply, error) {
	fmt.Printf("[service] GetAllReviews, req:#%v\n", req)
	reviews, err := s.uc.GetAllReviews(ctx, &biz.GetAllReviewsParam{
		UserID: req.UserID,
	})
	if err != nil {
		return nil, err
	}
	allreviewsinfo := []*pb.AllReviewsInfo{}
	for i := 0; i < len(reviews); i++ {
		allreviewsinfo = append(allreviewsinfo, &pb.AllReviewsInfo{
			UserID:       reviews[i].UserID,
			OrderID:      reviews[i].OrderID,
			ReviewID:     reviews[i].ReviewID,
			Score:        reviews[i].Score,
			ServiceScore: reviews[i].ServiceScore,
			ExpressScore: reviews[i].ExpressScore,
			Content:      reviews[i].Content,
			PicInfo:      reviews[i].PicInfo,
			VideoInfo:    reviews[i].VideoInfo,
			Status:       reviews[i].Status,
		})
	}
	return &pb.GetAllReviewsReply{Data: allreviewsinfo}, nil
}

func (s *ReviewService) ReplyReview(ctx context.Context, req *pb.ReplyReviewRequest) (*pb.ReplyReviewReply, error) {
	fmt.Printf("[service] ReplyReview, req:#%v\n", req)
	reply, err := s.uc.CreateReply(ctx, &biz.ReplyParam{
		ReviewID:  req.ReviewID,
		StoreID:   req.StoreID,
		Content:   req.Content,
		PicInfo:   req.PicInfo,
		VideoInfo: req.VideoInfo,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ReplyReviewReply{ReplyID: reply.ReplyID}, nil
}

func (s *ReviewService) ReplyReviewUpdate(ctx context.Context, req *pb.ReplyReviewUpdateRequest) (*pb.ReplyReviewUpdateReply, error) {
	fmt.Printf("[service] ReplyReviewUpdate, req:#%v\n", req)
	rowsaffected, err := s.uc.UpdateReply(ctx, &biz.ReplyUpdateParam{
		ReviewID:   req.ReviewID,
		StoreID:    req.StoreID,
		ReplyID:    req.ReplyID,
		UpdateTime: time.Now(),
		Content:    req.Content,
		PicInfo:    req.PicInfo,
		VideoInfo:  req.VideoInfo,
	})
	if err != nil {
		return nil, err
	}
	return &pb.ReplyReviewUpdateReply{ReplyID: req.ReplyID, Rowsaffected: rowsaffected}, nil
}

func (s *ReviewService) AppealReview(ctx context.Context, req *pb.AppealReviewRequest) (*pb.AppealReviewReply, error) {
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

func (s *ReviewService) AuditAppeal(ctx context.Context, req *pb.AuditAppealRequest) (*pb.AuditAppealReply, error) {
	fmt.Printf("[service] AuditAppeal, req:#%v\n", req)
	err := s.uc.AuditAppeal(ctx, &biz.AuditAppealParam{
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

// AuditReviewO端审核评价
func (s *ReviewService) AuditReview(ctx context.Context, req *pb.AuditReviewRequest) (*pb.AuditReviewReply, error) {
	fmt.Printf("[service] AuditReview, req:#%v\n", req)
	err := s.uc.AuditReview(ctx, &biz.AuditReviewParam{
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

// 商家ID搜索评价列表
func (s *ReviewService) ListReviewByStoreID(ctx context.Context, req *pb.ListReviewByStoreIDRequest) (*pb.ListReviewByStoreIDReply, error) {
	fmt.Printf("[service] ListReviewByStoreID, req:#%v\n", req)
	ret, err := s.uc.ListReviewByStoreID(ctx, req.StoreID, req.Page, req.Size)
	list := make([]*pb.ReviewInfo, 0, len(ret))
	for _, r := range ret {
		list = append(list, &pb.ReviewInfo{
			UserID:       r.UserID,
			OrderID:      r.OrderID,
			Score:        r.Score,
			ServiceScore: r.ServiceScore,
			ExpressScore: r.ExpressScore,
			Content:      r.Content,
			PicInfo:      r.PicInfo,
			VideoInfo:    r.VideoInfo,
			Status:       r.Status,
		})
	}
	return &pb.ListReviewByStoreIDReply{List: list}, err
}

func (s *ReviewService) ListReviewByContent(ctx context.Context, req *pb.ListReviewByContentRequest) (*pb.ListReviewByContentReply, error) {
	fmt.Printf("[service] ListReviewByContent, req:#%v\n", req)
	ret, err := s.uc.ListReviewByContent(ctx, req.Page, req.Size)
	list := make([]*pb.ReviewInfo, 0, len(ret))
	for _, r := range ret {
		list = append(list, &pb.ReviewInfo{
			UserID:       r.UserID,
			OrderID:      r.OrderID,
			Score:        r.Score,
			ServiceScore: r.ServiceScore,
			ExpressScore: r.ExpressScore,
			Content:      r.Content,
			PicInfo:      r.PicInfo,
			VideoInfo:    r.VideoInfo,
			Status:       r.Status,
		})
	}
	return &pb.ListReviewByContentReply{List: list}, err
}
