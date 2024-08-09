package data

import (
	"errors"
	"github.com/elastic/go-elasticsearch/v8"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"review-service/internal/conf"
	"review-service/internal/data/query"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewReviewRepo, NewDB, NewEsclient)

// Data .
type Data struct {
	// TODO wrapped database client
	//db *gorm.DB
	query *query.Query
	log   *log.Helper
	es    *elasticsearch.TypedClient
}

// NewData .
func NewData(db *gorm.DB, esClient *elasticsearch.TypedClient, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	//GEN生成query代码设置数据库对象
	query.SetDefault(db)
	return &Data{query: query.Q, log: log.NewHelper(logger), es: esClient}, cleanup, nil
}

func NewEsclient(config *conf.ElasticSearch) (*elasticsearch.TypedClient, error) {
	// ES 配置
	cfg := elasticsearch.Config{Addresses: config.Addresses}

	return elasticsearch.NewTypedClient(cfg)
}
func NewDB(cfg *conf.Data) (*gorm.DB, error) {
	switch strings.ToLower(cfg.Database.GetDriver()) {
	case "mysql":
		return gorm.Open(mysql.Open(cfg.Database.GetSource()))
	case "sqlite":
		return gorm.Open(sqlite.Open(cfg.Database.GetSource()))
	}
	return nil, errors.New("connect db fail: unsupported db driver")
}
