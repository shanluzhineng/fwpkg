package gorminit

import (
	"fmt"

	"github.com/shanluzhineng/fwpkg/app"
	"github.com/shanluzhineng/fwpkg/app/web"
	"github.com/shanluzhineng/fwpkg/system/log"

	"github.com/shanluzhineng/configurationx"
	"github.com/shanluzhineng/configurationx/options/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type IGormDbManager interface {
	GetDb(key string) *gorm.DB

	GetDefaultDb() *gorm.DB
}

type gormDbManager struct {
	_dbMap map[string]*gorm.DB

	_defaultDb *gorm.DB
}

var _ IGormDbManager = (*gormDbManager)(nil)

func init() {
	web.ConfigureService(serviceConfigurator)
}

func serviceConfigurator(web web.WebApplication) {
	app.Context.RegistInstanceAs(newGormDbManager(), new(IGormDbManager))
}

func newGormDbManager() *gormDbManager {
	m := &gormDbManager{
		_dbMap: make(map[string]*gorm.DB),
	}
	m.initGorm()
	return m
}

func (m *gormDbManager) GetDb(key string) *gorm.DB {
	return m._dbMap[key]
}

func (m *gormDbManager) GetDefaultDb() *gorm.DB {
	return m._defaultDb
}

func (m *gormDbManager) initGorm() {
	for eachKey, eachDb := range configurationx.GetInstance().Db.DbList {
		if eachDb.Disable {
			continue
		}
		var currentDB *gorm.DB
		switch eachDb.DbType {
		case db.DbType_Mysql:
			currentDB = m.gormMysqlByConfig(eachDb)
		default:
			continue
		}
		if currentDB != nil {
			m._dbMap[eachKey] = currentDB
		}
	}
	if defaultDB, ok := m._dbMap[db.AliasName_Default]; ok {
		m._defaultDb = defaultDB
	}
}

// GormMysqlByConfig 初始化Mysql数据库用过传入配置
func (m *gormDbManager) gormMysqlByConfig(dbConfiguration db.SpecializedDB) *gorm.DB {
	if len(dbConfiguration.Dbname) <= 0 {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       dbConfiguration.Dsn(), // DSN data source name
		DefaultStringSize:         191,                   // string 类型字段的默认长度
		SkipInitializeWithVersion: false,                 // 根据版本自动配置
	}
	db, err := gorm.Open(mysql.New(mysqlConfig), internalGorm.Config(dbConfiguration.GeneralDB))
	if err != nil {
		err = fmt.Errorf("无法初始化数据库,数据库名:%s,异常信息:%s", dbConfiguration.Dbname, err.Error())
		log.Logger.Error(err.Error())
		return nil
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(dbConfiguration.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbConfiguration.MaxOpenConns)
	return db
}
