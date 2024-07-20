package gorminit

import (
	"log"
	"os"
	"time"

	"github.com/shanluzhineng/configurationx/options/db"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dbBase interface {
	GetLogMode() string
}

var internalGorm = new(_gorm)

type _gorm struct{}

func (g *_gorm) Config(dbConfig db.GeneralDB) *gorm.Config {
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	_default := logger.New(newWriter(dbConfig.LogZap, log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
		SlowThreshold: 200 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})
	var logMode dbBase = &dbConfig

	switch logMode.GetLogMode() {
	case "silent", "Silent":
		config.Logger = _default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = _default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = _default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = _default.LogMode(logger.Info)
	default:
		config.Logger = _default.LogMode(logger.Info)
	}
	return config
}
