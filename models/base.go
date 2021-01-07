package models

import (
	"github.com/luoxiaojun1992/go-skeleton/services/db/sql/mysql"
)

type BaseModel struct {
}

type ModelInterface interface {
	Connection() string
}

type QueryBuilder struct {
	Model    ModelInterface
	DBClient *mysql.ConnWrapper
}

// NOTE When query with struct, GORM will only query with those fields has non-zero value, that means if your field’s value is 0, '', false or other zero values, it won’t be used to build query conditions
func (baseModel *BaseModel) Query(model ModelInterface) *QueryBuilder {
	return &QueryBuilder{
		Model:    model,
		DBClient: mysql.Connection(model.Connection()),
	}
}

func (qb *QueryBuilder) FindByPk(pk interface{}, retry bool) error {
	doFirst := func() error {
		return qb.DBClient.First(qb.Model, pk).Error
	}

	err := doFirst()

	if retry {
		if !qb.DBClient.InTransaction {
			if mysql.CausedByLostConnection(err) {
				return doFirst()
			}
		}
	}

	return err
}

func (qb *QueryBuilder) FirstByWhere(where func(dbClient *mysql.ConnWrapper) *mysql.ConnWrapper, retry bool) error {
	doFindByWhere := func() error {
		newDB := where(qb.DBClient)
		qb.DBClient = newDB
		return newDB.First(qb.Model).Error
	}

	err := doFindByWhere()

	if retry {
		if !qb.DBClient.InTransaction {
			if mysql.CausedByLostConnection(err) {
				return doFindByWhere()
			}
		}
	}

	return err
}
