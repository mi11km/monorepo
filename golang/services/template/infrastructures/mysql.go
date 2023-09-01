package infrastructures

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

func (cfg *MySQLConfig) FormatDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name,
	)
}

type MySQL struct {
	dsn string
	db  *sql.DB
}

func NewMySQL(dsn string) (*MySQL, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return &MySQL{
		dsn: dsn,
		db:  db,
	}, nil
}

func (m *MySQL) Close() error {
	return m.db.Close()
}

func (m *MySQL) Ping() error {
	return m.db.Ping()
}

/*
usage:
	mysql, err := infrastructures.NewMySQL(cfg.MySQL.FormatDSN())
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	if err := mysql.Ping(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
	defer func() {
		if err := mysql.Close(); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}()
*/
