package dao

import (
	"database/sql"
	"fmt"
	"log"

	"go.uber.org/zap"
	_ "github.com/lib/pq"
)

//go:generate mockery --name=DBService
type DBService interface {
	AddRow(message string) error
}

func GetDBService (config *YsqlConfig) DBService {
	db,err := config.GetDB()
	if err != nil {
		log.Fatal("DB service initialization failed")
	}
	return &YsqlDB{db}
}



type YsqlConfig struct {
	Host, User, Password, DbName string
	Port                         int
}

type YsqlDB struct {
	*sql.DB
}

func (c *YsqlConfig) getConnString() string {
	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DbName)
}

func (c *YsqlConfig) GetDB() (*sql.DB, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("logger initialization failed")
		return nil, err
	}
	defer logger.Sync()
	db, err := sql.Open("postgres", c.getConnString())
	if err != nil {
		logger.Error("connection establishment failed", zap.String("err", err.Error()))
		return nil, err
	}
	return db, nil
}

func (db *YsqlDB) AddRow(message string) error {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("logger initialization failed")
		return err
	}
	defer logger.Sync()
	insertStmt := "INSERT INTO kafkaMessages (id, message) VALUES (nextval('poc_sequence'),'" + message + "')"
	_, err = db.Exec(insertStmt)
	if err != nil {
		logger.Error("db insert failed", zap.String("err", err.Error()))
		return err
	}
	logger.Info("successfully inserted")
	return nil
}
