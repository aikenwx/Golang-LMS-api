package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Change accordingly

const DB_USERNAME = "root"
const DB_PASSWORD = ""
const DB_NAME = "production_lms_db"
const DB_HOST = "localhost"
const DB_PORT = "3306"

type Credentials struct {
	Username string
	Password string
	Name     string
	Host     string
	Port     string
	Debug    bool
}

type Connection struct {
	db *gorm.DB
}

func InitDefaultConnection() *Connection {
	return NewConnection(&Credentials{
		Username: DB_USERNAME,
		Password: DB_PASSWORD,
		Name:     DB_NAME,
		Host:     DB_HOST,
		Port:     DB_PORT,
	})
}

func NewConnection(credentials *Credentials) *Connection {
	return &Connection{db: connectDB(credentials)}
}

func (connection *Connection) GetDb() *gorm.DB {
	return connection.db
}

func connectDB(credentials *Credentials) *gorm.DB {
	var err error
	dsn := credentials.Username + ":" + credentials.Password + "@tcp" + "(" + credentials.Host + ":" + credentials.Port + ")/" +
		credentials.Name + "?" + "parseTime=true&loc=Local"

	var db *gorm.DB
	if credentials.Debug {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		fmt.Println("Error connecting to database : error=%v", err)
		return nil
	}

	return db
}
