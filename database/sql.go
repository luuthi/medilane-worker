package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"medilane-worker/config"
	"medilane-worker/models"
	"sync"
)

type MySqlDB interface {
	Init(cfg *config.Config)
	InsertManyNotification([]models.Notification) error
	GetFcmToken([]uint, *[]models.FcmToken) error
}

var once sync.Once

type SqlConnector struct {
	DB *gorm.DB
}

// singleton for api
var singletonSqlConnector MySqlDB

func GetInstance() MySqlDB {
	once.Do(func() { // <-- atomic, does not allow repeating
		singletonSqlConnector = &SqlConnector{}
	})
	return singletonSqlConnector
}

func SetInstance(obj MySqlDB) {
	singletonSqlConnector = obj
}

//NewClient new client for worker
func NewClient(cfg *config.Config) MySqlDB {
	sqlConnector := &SqlConnector{}
	sqlConnector.Init(cfg)
	return sqlConnector
}
func (obj *SqlConnector) Init(cfg *config.Config) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.DB.User, cfg.DB.Password, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)

	db, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   false,
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err.Error())
	}
	obj.DB = db
}

func (obj *SqlConnector) InsertManyNotification(data []models.Notification) error {
	return obj.DB.Table("notification").CreateInBatches(data, 10).Error
}

func (obj *SqlConnector) GetFcmToken(userIds []uint, tokens *[]models.FcmToken) error {
	return obj.DB.Table("fcm_token").Where("user IN ?", userIds).Find(&tokens).Error
}
