package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
	"review-service/internal/conf"
	"strings"
)

//gorm gen生成代码配置

var flagconf string

func connectDB(cfg *conf.Data_Database) *gorm.DB {
	if cfg == nil {
		panic(errors.New("GEN: connectDB fail, need cfg"))
	}
	switch strings.ToLower(cfg.GetDriver()) {
	case "mysql":
		db, err := gorm.Open(mysql.Open(cfg.GetSource()))
		if err != nil {
			panic(fmt.Errorf("connect db fail: %w", err))
		}
		return db
	case "sqlite":
		db, err := gorm.Open(sqlite.Open(cfg.GetSource()))
		if err != nil {
			panic(fmt.Errorf("connect db fail: %w", err))
		}
		return db
	default:
		panic(errors.New("GEN:connectDB failed: unsupported db driver"))
	}
}

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:       "../../internal/data/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	g.UseDB(connectDB(bc.Data.Database))
	g.ApplyBasic(g.GenerateAllTable()...)
	g.Execute()
}
