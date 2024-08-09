package job

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/segmentio/kafka-go"
	"review-job/internal/conf"
	"time"
)

// 评价数据流处理任务
// 自定义执行job的结构体，实现transport.server
type JobWorker struct {
	kafkaReader *kafka.Reader //Kafka reader
	esClient    *ESClient     //ES client
	log         *log.Helper
}

type ESClient struct {
	*elasticsearch.TypedClient
	index string
}

type Msg struct {
	Type     string `json:"type"`
	Database string `json:"database"`
	Table    string `json:"table"`
	IsDdl    bool   `json:"isDdl"`
	Data     []map[string]interface{}
}

func NewJobWorker(kafkaReader *kafka.Reader, esclient *ESClient, logger log.Logger) *JobWorker {
	return &JobWorker{
		kafkaReader: kafkaReader,
		esClient:    esclient,
		log:         log.NewHelper(logger),
	}
}

func NewKafkaReader(cfg *conf.Kafka) *kafka.Reader {
	return kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:  cfg.Brokers,
			GroupID:  cfg.GroupId,
			Topic:    cfg.Topic,
			MaxBytes: 1e7,
		},
	)
}

func NewEsclient(config *conf.ElasticSearch) *ESClient {
	// ES 配置
	cfg := elasticsearch.Config{Addresses: config.Addresses}

	// 创建客户端连接
	client, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		fmt.Printf("elasticsearch.NewTypedClient failed, err:%v\n", err)
		return nil
	}
	return &ESClient{
		TypedClient: client,
		index:       config.Index,
	}
}

func (jw JobWorker) Start(ctx context.Context) error {
	jw.log.Debug("job worker start")
	// 1.kafka中获取数据变更消息
	for {
		m, err := jw.kafkaReader.ReadMessage(ctx)
		if errors.Is(err, context.Canceled) {
			return nil
		}
		if err != nil {
			jw.log.Errorf("readmessage failed,err:%v\n", err)
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		// 2.将评价数据完整写入ES
		msg := new(Msg)
		err = json.Unmarshal(m.Value, msg)
		if err != nil {
			log.Errorf("unmarshal from kafka failed, err:%v\n", err)
			continue
		}
		//data process...
		if msg.Type == "INSERT" {
			//add
			for i := range msg.Data {
				jw.indexDocument(msg.Data[i])
			}
		} else {
			//update
			for i := range msg.Data {
				jw.updateDocument(msg.Data[i])
			}
		}
	}
	return nil
}
func (jw JobWorker) Stop(ctx context.Context) (err error) {
	jw.log.Debug("job worker stop.")
	if err = jw.kafkaReader.Close(); err != nil {
		jw.log.Errorf("reader close failed,err:%v\n", err)
	}
	return err
}

// Review 评价数据
type Review struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"userID"`
	Score       uint8     `json:"score"`
	Content     string    `json:"content"`
	Tags        []Tag     `json:"tags"`
	Status      int       `json:"status"`
	PublishTime time.Time `json:"publishDate"`
}

// Tag 评价标签
type Tag struct {
	Code  int    `json:"code"`
	Title string `json:"title"`
}

// indexDocument 索引文档
func (jw JobWorker) indexDocument(d map[string]interface{}) {
	// 添加文档
	resp, err := jw.esClient.Index(jw.esClient.index).
		Id(d["review_id"].(string)).
		Document(d).
		Do(context.Background())
	if err != nil {
		fmt.Printf("indexing document failed, err:%v\n", err)
		return
	}
	fmt.Printf("result:%#v\n", resp.Result)
}

// updateDocument 更新文档
func (jw JobWorker) updateDocument(d map[string]interface{}) {
	// 修改后的结构体变量
	resp, err := jw.esClient.Update(jw.esClient.index, d["review_id"].(string)).
		Doc(d). // 使用结构体变量更新
		Do(context.Background())
	if err != nil {
		jw.log.Errorf("update document failed, err:%v\n", err)
		return
	}
	jw.log.Errorf("result:%v\n", resp.Result)
}
