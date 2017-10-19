package libraries

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/michaelbui/AWSSG-1710/entities"
	"github.com/michaelbui/AWSSG-1710/types"
)

type (
	SqliteConnector struct {
		configs *types.DBConfig
		conn    *gorm.DB
	}
)

func NewSqliteConnection(configs *types.DBConfig) *SqliteConnector {
	return &SqliteConnector{
		configs: configs,
	}
}

func (sql *SqliteConnector) Save(entity interface{}) error {
	if err := sql.open(); err != nil {
		return err
	}
	defer sql.close()
	var conn *gorm.DB
	if sql.conn.NewRecord(entity) {
		conn = sql.conn.Create(entity)
	} else {
		conn = sql.conn.Save(entity)
	}
	return conn.Error
}

func (sql *SqliteConnector) Find(output interface{}) error {
	if err := sql.open(); err != nil {
		return err
	}
	defer sql.close()
	conn := sql.conn.Find(output)
	return conn.Error
}

func (sql *SqliteConnector) open() error {
	conn, err := gorm.Open("sqlite3", sql.configs.Dsn)
	if err != nil {
		return err
	}
	if !sql.configs.Initiated {
		sql.init(conn)
		sql.configs.Initiated = true
	}
	sql.conn = conn
	return nil
}

func (sql *SqliteConnector) close() {
	sql.conn.Close()
}

func (sql *SqliteConnector) init(conn *gorm.DB) {
	conn.DropTableIfExists(&entities.File{})
	conn.CreateTable(&entities.File{})
}
