package sql

import (
	"log"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"fmt"
)

type Config struct {
	MySQL struct {
		Host        string `yaml:"host"`
		Port        int    `yaml:"port"`
		User        string `yaml:"user"`
		Password    string `yaml:"password"`
		DBName      string `yaml:"dbname"`
		Charset     string `yaml:"charset"`
		ParseTime   bool   `yaml:"parseTime"` // 将时间戳转换为时间
		MaxIdleConns int    `yaml:"maxIdleConns"`
		MaxOpenConns int    `yaml:"maxOpenConns"`
	} `yaml:"mysql"` // mysql 配置
}

func LoadConfig() *Config {
	viper.SetConfigName("mysql") // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath("/home/www/flowers/config") // 配置文件路径
	
	if err := viper.ReadInConfig(); err != nil { // 读取配置文件
		log.Fatalf("Error reading config file: %v", err)
	}
	
	var cfg Config // 定义配置结构体 用于保存配置文件中的数据 
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode config: %v", err)
	}
	return &cfg // 返回配置结构体指针
}

func InitDB(cfg *Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.DBName,
		cfg.MySQL.Charset,
		cfg.MySQL.ParseTime,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)

	return db
}