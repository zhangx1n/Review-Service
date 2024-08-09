package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	pb "review-service/api/review/v1"
	"review-service/internal/data/model"
	"review-service/pkg/snowflake"
	"strings"
	"time"
)

// GreeterRepo is a Greater repo.
type ReviewRepo interface {
	SaveReview(context.Context, *model.ReviewInfo) (*model.ReviewInfo, error)
	GetReviewByOrderID(ctx context.Context, OrderID int64) ([]*model.ReviewInfo, error)
	GetReviewByReviewID(ctx context.Context, ReviewID int64) (*model.ReviewInfo, error)
	UpdateReview(ctx context.Context, updatereview *model.ReviewInfo) (int64, error)
	DeleteReview(ctx context.Context, deletereview *model.ReviewInfo) (int64, error)
	GetReviewByUserID(ctx context.Context, id int64) ([]*model.ReviewInfo, error)

	SaveReply(ctx context.Context, reply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error)
	UpdateReply(ctx context.Context, reply *model.ReviewReplyInfo) (int64, error)
	SaveAppeal(ctx context.Context, appeal *model.ReviewAppealInfo) (int64, error)

	AuditAppeal(ctx context.Context, audit *model.ReviewAppealInfo) error
	AuditReview(ctx context.Context, audit *model.ReviewInfo) error

	ListReviewByStoreID(ctx context.Context, id int64, offset int32, limit int32) ([]*MyReviewInfo, error)
	ListReviewByContent(ctx context.Context, offset int32, limit int32) ([]*MyReviewInfo, error)
}

type ReviewerUsecase struct {
	repo ReviewRepo
	log  *log.Helper
}

func NewReviewerUsecase(repo ReviewRepo, logger log.Logger) *ReviewerUsecase {
	return &ReviewerUsecase{repo: repo, log: log.NewHelper(logger)}
}

// c端
func (uc *ReviewerUsecase) CreateReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReview, req:%v", review)
	//1.数据校验
	//1.1参数基础校验
	//validator review.proto->review.pb.validate.go
	//1.2参数业务校验
	reviews, err := uc.repo.GetReviewByOrderID(ctx, review.OrderID)
	if err != nil {
		return nil, pb.ErrorDbFailed("查询数据库失败")
	}
	if len(reviews) > 0 {
		//已经评价过
		return nil, pb.ErrorOrderReviewed("订单:%v 已评价过", review.OrderID)
	}
	//2.生成reviewer ID
	//snowflake or company service
	//snowflake main里面init了snowflake.node用于生成
	review.ReviewID = snowflake.GenID()
	//3.查询订单与商品快照信息
	//实际业务场景下查询，rpc调用订单服务和商家服务
	//4.拼装数据入库
	return uc.repo.SaveReview(ctx, review)
}

func (uc *ReviewerUsecase) GetReview(ctx context.Context, m *model.ReviewInfo) (*pb.GetReviewReply, error) {
	uc.log.WithContext(ctx).Debugf("[biz] GetReview, reviewID:%v\n", m.ReviewID)
	review, err := uc.repo.GetReviewByReviewID(ctx, m.ReviewID)
	if err != nil {
		return nil, pb.ErrorDbFailed("查询数据库失败")
	}
	return &pb.GetReviewReply{
		UserID:       review.UserID,
		OrderID:      review.OrderID,
		Score:        review.Score,
		ServiceScore: review.ServiceScore,
		ExpressScore: review.ExpressScore,
		Content:      review.Content,
		PicInfo:      review.PicInfo,
		VideoInfo:    review.VideoInfo,
		Status:       review.Status,
	}, nil
}

func (uc *ReviewerUsecase) GetAllReviews(ctx context.Context, param *GetAllReviewsParam) ([]*model.ReviewInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] GetAllReviews, req:%v", param)
	//1.数据校验
	//1.1参数基础校验
	//validator review.proto->review.pb.validate.go
	//1.2参数业务校验
	reviews, err := uc.repo.GetReviewByUserID(ctx, param.UserID)
	if err != nil {
		return nil, pb.ErrorDbFailed("查询数据库失败")
	}
	if len(reviews) == 0 {
		//该用户没发表过评价
		return nil, pb.ErrorOrderReviewed("用户:%v 未发表过评价", param.UserID)
	}
	//3.返回数据
	return reviews, err
}

func (uc *ReviewerUsecase) UpdateReview(ctx context.Context, updatereview *model.ReviewInfo) (int64, error) {
	uc.log.WithContext(ctx).Debugf("[biz] UpdateReview, reviewID:%v\n", updatereview.ReviewID)
	//1.数据校验
	//查询是否有符合reviewid的review记录
	review, err := uc.repo.GetReviewByReviewID(ctx, updatereview.ReviewID)
	if err != nil {
		return 0, pb.ErrorDbFailed("查询数据库失败")
	}
	fmt.Printf("%v\n", review)
	//1.2水平越权校验，不是自己创建的review不能去删除
	if review.UserID != updatereview.UserID {
		return 0, pb.ErrorInvalidParams("参数有误，UserID不匹配")
	}
	//2.更新review并返回(受影响行数)
	return uc.repo.UpdateReview(ctx, updatereview)
}

func (uc *ReviewerUsecase) DeleteReview(ctx context.Context, deletereview *model.ReviewInfo) (int64, error) {
	uc.log.WithContext(ctx).Debugf("[biz] DeleteReview, reviewID:%v\n", deletereview.ReviewID)
	//1.数据校验
	//1.1查询是否有符合reviewid的review记录
	review, err := uc.repo.GetReviewByReviewID(ctx, deletereview.ReviewID)
	if err != nil {
		return 0, pb.ErrorDbFailed("查询数据库失败")
	}
	fmt.Printf("%v\n", review)
	//1.2水平越权校验，不是自己创建的review不能去删除
	if review.UserID != deletereview.UserID {
		return 0, pb.ErrorInvalidParams("参数有误，UserID不匹配")
	}
	//2.删除review并返回(受影响行数)
	return uc.repo.DeleteReview(ctx, deletereview)
}

// b端
func (uc *ReviewerUsecase) CreateReply(ctx context.Context, param *ReplyParam) (*model.ReviewReplyInfo, error) {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReply, reviewID:%v\n", param)
	//业务校验放data层了
	reply := &model.ReviewReplyInfo{
		ReplyID:   snowflake.GenID(),
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	}
	return uc.repo.SaveReply(ctx, reply)
}

func (uc *ReviewerUsecase) UpdateReply(ctx context.Context, param *ReplyUpdateParam) (int64, error) {
	uc.log.WithContext(ctx).Debugf("[biz] CreateReply, reviewID:%v\n", param)
	//业务校验放data层了
	reply := &model.ReviewReplyInfo{
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		ReplyID:   param.ReplyID,
		UpdateAt:  param.UpdateTime,
		Content:   param.Content,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	}
	return uc.repo.UpdateReply(ctx, reply)
}

func (uc *ReviewerUsecase) AppealReview(ctx context.Context, param *AppealParam) (int64, error) {
	uc.log.WithContext(ctx).Debugf("[biz] AppealReview, reviewID:%v\n", param.ReviewID)
	//1.数据校验
	//查询是否有符合reviewid的review记录
	review, err := uc.repo.GetReviewByReviewID(ctx, param.ReviewID)
	if err != nil {
		return 0, pb.ErrorDbFailed("查询数据库失败")
	}
	//1.2水平越权校验，不是自己商店下的review不能去申诉
	if review.StoreID != param.StoreID {
		return 0, pb.ErrorInvalidParams("参数有误，StoreID不匹配")
	}
	//2.创建appealID
	appeal := &model.ReviewAppealInfo{
		ReviewID:  param.ReviewID,
		StoreID:   param.StoreID,
		Reason:    param.Reason,
		Content:   param.Reason,
		PicInfo:   param.PicInfo,
		VideoInfo: param.VideoInfo,
	}
	//2.创建appeal,校验appeal表中有无对于该review的申诉
	//3.返回(申诉的appealID)
	return uc.repo.SaveAppeal(ctx, appeal)
}

func (uc *ReviewerUsecase) AuditAppeal(ctx context.Context, param *AuditAppealParam) error {
	uc.log.WithContext(ctx).Debugf("[biz] AuditAppeal, reviewID:%v\n", param)
	//业务校验放data层了
	audit := &model.ReviewAppealInfo{
		AppealID:  param.AppealID,
		OpUser:    param.OpUser,
		Reason:    param.OpReason,
		OpRemarks: param.OpRemarks,
		Status:    param.Status,
	}
	return uc.repo.AuditAppeal(ctx, audit)
}

func (uc *ReviewerUsecase) AuditReview(ctx context.Context, param *AuditReviewParam) error {
	uc.log.WithContext(ctx).Debugf("[biz] AuditReview, reviewID:%v\n", param)
	//业务校验放data层了
	audit := &model.ReviewInfo{
		ReviewID:  param.ReviewID,
		OpUser:    param.OpUser,
		OpReason:  param.OpReason,
		OpRemarks: param.OpRemarks,
		Status:    param.Status,
	}
	return uc.repo.AuditReview(ctx, audit)
}

func (uc *ReviewerUsecase) ListReviewByStoreID(ctx context.Context, id int64, page int32, size int32) ([]*MyReviewInfo, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 10
	}
	offset := (page - 1) * size
	limit := size
	uc.log.WithContext(ctx).Debugf("[biz] ListReviewByStoreID userID:%v\n", id)
	return uc.repo.ListReviewByStoreID(ctx, id, offset, limit)
}

func (uc *ReviewerUsecase) ListReviewByContent(ctx context.Context, page int32, size int32) ([]*MyReviewInfo, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 10
	}
	offset := (page - 1) * size
	limit := size
	uc.log.WithContext(ctx).Debugf("[biz] ListReviewByContent userID:%v\n")
	return uc.repo.ListReviewByContent(ctx, offset, limit)
}

type MyReviewInfo struct {
	*model.ReviewInfo
	CreateAt     MyTime `json:"create_at"`
	UpdateAt     MyTime `json:"update_at"`
	ID           int64  `json:"id,string"`
	Version      int32  `json:"version,string"`
	ReviewID     int64  `json:"review_id,string"`
	Score        int32  `json:"score,string"`
	ServiceScore int32  `json:"service_score,string"`
	ExpressScore int32  `json:"express_score,string"`
	HasMedia     int32  `json:"has_media,string"`
	OrderID      int64  `json:"order_id,string"`
	SkuID        int64  `json:"sku_id,string"`
	SpuID        int64  `json:"spu_id,string"`
	StoreID      int64  `json:"store_id,string"`
	UserID       int64  `json:"user_id,string"`
	Anonymous    int32  `json:"anonymous,string"`
	Status       int32  `json:"status,string"`
	IsDefault    int32  `json:"is_default,string"`
	HasReply     int32  `json:"has_reply,string"`
}

type MyTime time.Time

func (t *MyTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "\"")
	tmp, err := time.Parse(time.DateTime, s)
	if err != nil {
		return err
	}
	*t = MyTime(tmp)
	return nil
}
