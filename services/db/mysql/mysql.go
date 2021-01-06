package mysql

import (
	"database/sql"
	"github.com/gookit/config/v2"
	"github.com/luoxiaojun1992/go-skeleton/services/helper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"strings"
	"sync"
	"time"
)

const DEFAULT_CONNECTION = "default"

type ConnWrapper struct {
	*gorm.DB
	InTransaction bool
}

func (cw *ConnWrapper) BeginTx(opts ...*sql.TxOptions) *ConnWrapper {
	newConn := &ConnWrapper{cw.Begin(opts...), true}
	cw.InTransaction = true
	return newConn
}

func (cw *ConnWrapper) CommitTx(originCw *ConnWrapper) *ConnWrapper {
	cw.Commit()
	cw.InTransaction = false
	originCw.InTransaction = false
	return cw
}

func (cw *ConnWrapper) Txn(fc func(tx *ConnWrapper) error, opts ...*sql.TxOptions) error {
	errTxn := cw.Transaction(func(tx *gorm.DB) error {
		cw.InTransaction = true
		return fc(&ConnWrapper{tx, true})
	}, opts...)
	cw.InTransaction = false
	return errTxn
}

var Clients map[string]*gorm.DB

var dbCreateLock sync.Mutex

func init() {
	Clients = make(map[string]*gorm.DB)
}

func Setup() {
	connsConfig := make(map[string]interface{})
	errGetConnsConfig := config.MapStruct("db.connections", &connsConfig)
	helper.CheckErrThenAbort("failed to get db connections config", errGetConnsConfig)
	for connection, _ := range connsConfig {
		Create(connection)
	}
}

func Connection(connection string) *ConnWrapper {
	if len(connection) == 0 {
		defaultConnection := config.String("db.defaultConnection", DEFAULT_CONNECTION)
		connection = defaultConnection
	}

	if db, ok := Clients[connection]; ok {
		return &ConnWrapper{db, false}
	} else {
		log.Panicln("failed to get db [" + connection + "] connection: not existed")
	}

	return nil
}

func DefaultConnection() *ConnWrapper {
	return Connection("")
}

func Create(connection string) *gorm.DB {
	if len(connection) == 0 {
		defaultConnection := config.String("db.defaultConnection", DEFAULT_CONNECTION)
		connection = defaultConnection
	}

	dbCreateLock.Lock()
	defer dbCreateLock.Unlock()

	if db, ok := Clients[connection]; ok {
		return db
	}

	//Fetching config & creating db client & caching db client
	host := config.String("db.connections." + connection + ".host")
	port := config.String("db.connections." + connection + ".port")
	username := config.String("db.connections." + connection + ".username")
	password := config.String("db.connections." + connection + ".password")
	dbName := config.String("db.connections." + connection + ".dbName")
	charset := config.String("db.connections." + connection + ".charset")
	tablePrefix := config.String("db.connections." + connection + ".table_prefix")

	newDB, errOpenDB := gorm.Open(
		mysql.Open(username+":"+password+"@tcp("+host+":"+port+")/"+dbName+"?charset="+charset),
		&gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			DryRun:                 false,
			DisableAutomaticPing:   false,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   tablePrefix,
				SingularTable: true,
			},
		})

	helper.CheckErrThenAbort("failed to connect mysql ["+connection+"]", errOpenDB)

	sqlDriver, errSqlDriver := newDB.DB()
	helper.CheckErrThenAbort("failed to connect mysql ["+connection+"]", errSqlDriver)

	sqlDriver.SetMaxIdleConns(config.Int("db.connections."+connection+".maxIdleConns", 50))
	sqlDriver.SetMaxOpenConns(config.Int("db.connections."+connection+".maxOpenConns", 100))

	lifeTime, errParseDuration := time.ParseDuration(config.String("db.connections."+connection+".connMaxLifetime", "300s"))
	helper.CheckErrThenAbort("failed to parse db ["+connection+"] conn lifetime", errParseDuration)
	sqlDriver.SetConnMaxLifetime(lifeTime)

	debug := config.String("db.connections."+connection+".debug", "0")
	if debug == "1" {
		newDB.Debug()
	}

	Clients[connection] = newDB

	return newDB
}

func CreateDefault() *gorm.DB {
	return Create("")
}

func CloseClients() {
	for connection, client := range Clients {
		sqlDriver, errSqlDriver := client.DB()
		helper.CheckErrThenAbort("failed to close mysql ["+connection+"] client", errSqlDriver)
		errCloseClient := sqlDriver.Close()
		helper.CheckErrThenAbort("failed to close mysql ["+connection+"] client", errCloseClient)
	}
}

func CausedByLostConnection(err error) bool {
	if !helper.CheckErr(err) {
		return false
	}

	if strings.Contains(err.Error(), "invalid connection") {
		return true
	}

	if strings.Contains(err.Error(), "sql: database is closed") {
		return true
	}

	if strings.Contains(err.Error(), "closing bad idle connection: EOF") {
		return true
	}

	if strings.Contains(err.Error(), "connect: network is unreachable") {
		return true
	}

	if strings.Contains(err.Error(), "connect: operation timed out") {
		return true
	}

	if strings.Contains(err.Error(), "no such host") {
		return true
	}

	if strings.Contains(err.Error(), "connect: connection refused") {
		return true
	}

	if strings.Contains(err.Error(), "closing bad idle connection: connection reset by peer") {
		return true
	}

	if strings.Contains(err.Error(), "server has gone away") {
		return true
	}

	if strings.Contains(err.Error(), "no connection to the server") {
		return true
	}

	if strings.Contains(err.Error(), "Lost connection") {
		return true
	}

	if strings.Contains(err.Error(), "is dead or not enabled") {
		return true
	}

	if strings.Contains(err.Error(), "Error while sending") {
		return true
	}

	if strings.Contains(err.Error(), "decryption failed or bad record mac") {
		return true
	}

	if strings.Contains(err.Error(), "server closed the connection unexpectedly") {
		return true
	}

	if strings.Contains(err.Error(), "SSL connection has been closed unexpectedly") {
		return true
	}

	if strings.Contains(err.Error(), "Error writing data to the connection") {
		return true
	}

	if strings.Contains(err.Error(), "Resource deadlock avoided") {
		return true
	}

	if strings.Contains(err.Error(), "Transaction() on null") {
		return true
	}

	if strings.Contains(err.Error(), "child connection forced to terminate due to client_idle_limit") {
		return true
	}

	if strings.Contains(err.Error(), "query_wait_timeout") {
		return true
	}

	if strings.Contains(err.Error(), "reset by peer") {
		return true
	}

	if strings.Contains(err.Error(), "Physical connection is not usable") {
		return true
	}

	if strings.Contains(err.Error(), "TCP Provider: Error code 0x68") {
		return true
	}

	if strings.Contains(err.Error(), "Name or service not known") {
		return true
	}

	return false
}
