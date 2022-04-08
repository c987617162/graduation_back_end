package system

import (
	"database/sql"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"server/config"
	"server/global"
	"server/model/system"
	"server/model/system/request"
)

type InitDBService struct{}

// writeMysqlConfig mysql回写配置
// Author [SliverHorn](https://github.com/SliverHorn)
// Author [songzhibin97](https://github.com/songzhibin97)
func (initDBService *InitDBService) writeMysqlConfig(mysql config.Mysql) error {
	global.GVA_CONFIG.Mysql = mysql
	global.GVA_VP.Set("mysql", mysql)
	global.GVA_VP.Set("jwt.signing-key", uuid.NewV4().String())
	return global.GVA_VP.WriteConfig()
}

// initTables 初始化表
func (initDBService *InitDBService) initTables() error {
	return global.GVA_DB.AutoMigrate(
		system.SysApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.SysAuthority{},
		system.JwtBlacklist{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
	)
}

// createDatabase 创建数据库(mysql)
func (initDBService *InitDBService) createDatabase(dsn string, driver string, createSql string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}

// InitDB 初始化数据库
func (initDBService *InitDBService) InitDB(conf request.InitDB) error {
	dsn := conf.MysqlEmptyDsn()
	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;", conf.DBName)
	if err := initDBService.createDatabase(dsn, "mysql", createSql); err != nil {
		return err
	}
	mysqlConfig := conf.ToMysqlConfig()
	if mysqlConfig.Dbname == "" {
		return nil
	} // 如果没有数据库名, 则跳出初始化数据

	if db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       mysqlConfig.Dsn(), // DSN data source name
		DefaultStringSize:         191,               // string 类型字段的默认长度
		SkipInitializeWithVersion: true,              // 根据版本自动配置
	}), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}); err != nil {
		return nil
	} else {
		global.GVA_DB = db
	}

	if err := initDBService.initTables(); err != nil {
		global.GVA_DB = nil
		return err
	}

	if err := initDBService.writeMysqlConfig(mysqlConfig); err != nil {
		return err
	}

	return nil
}
