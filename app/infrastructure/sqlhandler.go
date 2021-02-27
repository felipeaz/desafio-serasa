package infrastructure

import (
	"fmt"
	"os"

	"github.com/FelipeAz/desafio-serasa/app/entity"
	"github.com/FelipeAz/desafio-serasa/app/interfaces"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// SQLHandler contem a conexao do banco de dados.
type SQLHandler struct {
	db *gorm.DB
}

// NewSQLHandler inicia uma conexao com o Banco de dados SQL.
func NewSQLHandler() (interfaces.SQLHandler, error) {
	dbHandler := &SQLHandler{}
	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/Serasa?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQLUSER"), os.Getenv("MYSQLPASS"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	dbHandler.db = db
	dbHandler.migrateTables()

	return dbHandler, nil
}

// CloseConnection fecha a conexao com o banco de dados.
func (env *SQLHandler) CloseConnection() error {
	sql, err := env.db.DB()
	if err != nil {
		return err
	}

	sql.Close()

	return nil
}

func (env *SQLHandler) migrateTables() {
	env.db.Migrator().AutoMigrate(&entity.Negativacao{})
}