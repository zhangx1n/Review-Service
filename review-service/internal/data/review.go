package data

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	pb "review-service/api/review/v1"
	"review-service/internal/biz"
	"review-service/internal/data/model"
	"review-service/internal/data/query"
	"review-service/pkg/snowflake"
)

type reviewRepo struct {
	data *Data
	log  *log.Helper
}

// NewReviewRepo .
func NewReviewRepo(data *Data, logger log.Logger) biz.ReviewRepo {
	return &reviewRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *reviewRepo) SaveReview(ctx context.Context, review *model.ReviewInfo) (*model.ReviewInfo, error) {
	err := r.data.query.ReviewInfo.WithContext(ctx).Save(review)
	return review, err
}

func (r *reviewRepo) GetReviewByOrderID(ctx context.Context, OrderID int64) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Debug().
		Where(r.data.query.ReviewInfo.OrderID.Eq(OrderID)).Find()
}

func (r *reviewRepo) GetReviewByReviewID(ctx context.Context, ReviewID int64) (*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Debug().Where(r.data.query.ReviewInfo.ReviewID.Eq(ReviewID)).First()
}

func (r *reviewRepo) GetReviewByUserID(ctx context.Context, id int64) ([]*model.ReviewInfo, error) {
	return r.data.query.ReviewInfo.WithContext(ctx).Debug().
		Where(r.data.query.ReviewInfo.UserID.Eq(id)).Find()
}

func (r *reviewRepo) UpdateReview(ctx context.Context, updatereview *model.ReviewInfo) (int64, error) {
	info, err := r.data.query.ReviewInfo.WithContext(ctx).Debug().
		Where(r.data.query.ReviewInfo.ReviewID.Eq(updatereview.ReviewID)).Updates(updatereview)
	return info.RowsAffected, err
}

func (r *reviewRepo) DeleteReview(ctx context.Context, deletereview *model.ReviewInfo) (int64, error) {
	info, err := r.data.query.ReviewInfo.WithContext(ctx).Debug().
		Where(r.data.query.ReviewInfo.ReviewID.Eq(deletereview.ReviewID)).Delete()
	return info.RowsAffected, err
}

func (r *reviewRepo) SaveReply(ctx context.Context, reply *model.ReviewReplyInfo) (*model.ReviewReplyInfo, error) {
	//1.数据校验
	//1.1数据合法性(已回复的不能再回复)
	//查评价id，看是否已回复
	review, err := r.data.query.ReviewInfo.WithContext(ctx).
		Where(r.data.query.ReviewInfo.ReviewID.Eq(reply.ReviewID)).First()
	if err != nil {
		return nil, err
	}
	if review.HasReply == 1 {
		return nil, pb.ErrorReviewReplied("该评价已回复")
	}
	//1.2水平越权校验(只能回复自己商家下的review)
	if review.StoreID != reply.StoreID {
		return nil, pb.ErrorInvalidParams("参数错误，商家ID不匹配")
	}
	//2.更新数据(review,reviewreply同时更新)
	//事务
	err = r.data.query.Transaction(func(tx *query.Query) error {
		//回复表插入
		if err = tx.ReviewReplyInfo.WithContext(ctx).Debug().
			Save(reply); err != nil {
			r.log.WithContext(ctx).Errorf("savereply create reply failed ,err:%v\n", err)
			return err
		}
		//评价表hasreply字段更新
		if _, err = tx.ReviewInfo.WithContext(ctx).Debug().
			Where(tx.ReviewInfo.ReviewID.Eq(reply.ReviewID)).Update(tx.ReviewInfo.HasReply, 1); err != nil {
			r.log.WithContext(ctx).Errorf("savereply update review failed ,err:%v\n", err)
			return err
		}
		return nil
	})
	//3.返回
	return reply, err
}

func (r *reviewRepo) UpdateReply(ctx context.Context, reply *model.ReviewReplyInfo) (int64, error) {
	//1.数据校验
	//1.1数据合法性(必须是已经回复了的review)
	//查评价id，看是否已回复
	reviewreply, err := r.data.query.ReviewReplyInfo.WithContext(ctx).
		Where(r.data.query.ReviewReplyInfo.ReplyID.Eq(reply.ReplyID)).First()
	if err != nil {
		//没这条ID的回复...
		return 0, err
	}
	if reviewreply.Content == "" {
		return 0, pb.ErrorReviewReplied("该评价还未回复，不用修改")
	}
	//1.2水平越权校验(只能修改自己商家下的reply)
	if reviewreply.StoreID != reply.StoreID {
		return 0, pb.ErrorInvalidParams("参数错误，商家ID不匹配")
	}
	//2.更新数据(reviewreply更新)
	info, err := r.data.query.ReviewReplyInfo.WithContext(ctx).Debug().
		Where(r.data.query.ReviewReplyInfo.ReplyID.Eq(reviewreply.ReplyID)).Updates(reply)
	//3.返回
	return info.RowsAffected, err
}

func (r *reviewRepo) SaveAppeal(ctx context.Context, appeal *model.ReviewAppealInfo) (int64, error) {
	appeals, err := r.data.query.ReviewAppealInfo.WithContext(ctx).
		Where(r.data.query.ReviewAppealInfo.ReviewID.Eq(appeal.ReviewID)).First()
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	if err == nil && appeals.Status > 10 {
		return 0, pb.ErrorAppealAudited("该评价已有审核过的申诉记录")
	}
	//if appeals != nil {
	//	return 0, pb.ErrorReviewAppealed("该review正在申诉，不能重复发起申诉")
	//}
	if appeals == nil {
		//1.没有申诉记录则创建
		appeal.AppealID = snowflake.GenID()
	} else {
		//2.有申诉记录但待审核，需要更新
		//只要更新申诉内容,appeal_id复用查询到的appeal记录的appeal_id
		//storeID,reviewID已经保证一致(若已存在),只可能更新原因、内容等字段...
		appeal.AppealID = appeals.AppealID
	}
	//.save方法根据主键更新或新建
	err = r.data.query.ReviewAppealInfo.WithContext(ctx).Debug().Save(appeal)
	return appeal.AppealID, err
}

func (r *reviewRepo) AuditReview(ctx context.Context, audit *model.ReviewInfo) error {
	review_info, err := r.data.query.ReviewInfo.WithContext(ctx).Debug().
		Where(r.data.query.ReviewInfo.ReviewID.Eq(audit.ReviewID)).First()
	if review_info == nil {
		//没有这个申诉
		return pb.ErrorDbFailed("没有这个评价ID！")
	}
	_, err = r.data.query.ReviewInfo.WithContext(ctx).Debug().
		Where(r.data.query.ReviewInfo.ReviewID.Eq(audit.ReviewID)).Updates(audit)
	return err
}

func (r *reviewRepo) AuditAppeal(ctx context.Context, audit *model.ReviewAppealInfo) error {
	appeal_info, err := r.data.query.ReviewAppealInfo.WithContext(ctx).Debug().
		Where(r.data.query.ReviewAppealInfo.AppealID.Eq(audit.AppealID)).First()
	if appeal_info == nil {
		//没有这个申诉
		return pb.ErrorDbFailed("没有这个申诉ID！")
	}
	audit.ReviewID = appeal_info.ReviewID
	err = r.data.query.Transaction(func(tx *query.Query) error {
		if _, err := tx.ReviewAppealInfo.WithContext(ctx).Debug().
			Where(r.data.query.ReviewAppealInfo.AppealID.Eq(audit.AppealID)).
			Updates(map[string]interface{}{
				"status":  audit.Status,
				"op_user": audit.OpUser,
				"reason":  audit.Reason,
			}); err != nil {
			return err
		}
		if audit.Status == 20 { //申诉通过则需要隐藏评价
			if _, err := tx.ReviewInfo.WithContext(ctx).Where(tx.ReviewInfo.ReviewID.Eq(audit.ReviewID)).
				Update(tx.ReviewInfo.Status, 40); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

// 根据storeID分页查询评价
func (r *reviewRepo) ListReviewByStoreID(ctx context.Context, id int64, offset int32, limit int32) ([]*biz.MyReviewInfo, error) {
	//ES里查询评价
	resp, err := r.data.es.Search().Index("review").From(int(offset)).Size(int(limit)).
		Query(&types.Query{
			Bool: &types.BoolQuery{
				Filter: []types.Query{
					{
						Term: map[string]types.TermQuery{
							"store_id": {Value: id},
						},
					},
				},
			},
		},
		).Do(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("es result: total:%v\n", resp.Hits.Total.Value)
	//反序列化数据
	//resp.Hits.Hits[0].Source_ = json.RawMessage -> model.ReviewInfo
	list := make([]*biz.MyReviewInfo, 0, resp.Hits.Total.Value)
	for _, hit := range resp.Hits.Hits {
		tmp := &biz.MyReviewInfo{}
		if err = json.Unmarshal(hit.Source_, tmp); err != nil {
			r.log.Errorf("json.unmarshal(hit.source_) failed, err:%v\n", err)
			continue
		}
		list = append(list, tmp)
	}
	return list, nil
}

func (r *reviewRepo) ListReviewByContent(ctx context.Context, offset int32, limit int32) ([]*biz.MyReviewInfo, error) {
	//ES里查询评价-评价不为空
	resp, err := r.data.es.Search().Index("review").From(int(offset)).Size(int(limit)).
		Query(&types.Query{
			Bool: &types.BoolQuery{
				Must: []types.Query{
					{
						Exists: &types.ExistsQuery{
							Field: "content",
						},
					},
				},
			},
		},
		).Do(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("es result: total:%v\n", resp.Hits.Total.Value)
	//反序列化数据
	//resp.Hits.Hits[0].Source_ = json.RawMessage -> model.ReviewInfo
	list := make([]*biz.MyReviewInfo, 0, resp.Hits.Total.Value)
	for _, hit := range resp.Hits.Hits {
		tmp := &biz.MyReviewInfo{}
		if err = json.Unmarshal(hit.Source_, tmp); err != nil {
			r.log.Errorf("json.unmarshal(hit.source_) failed, err:%v\n", err)
			continue
		}
		list = append(list, tmp)
	}
	return list, nil
}
