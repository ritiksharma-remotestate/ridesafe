package database

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
	_"github.com/lib/pq"
		"github.com/sirupsen/logrus"

	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var(
	Ridesafe *sqlx.DB
)

type SSLMode string
 
const (
	SSLModeEnable SSLMode="enable"
	SSLModeDisable SSLMode="disable"
)

func ConnectAndMigrate(host,port,databaseName,user,password string, sslmode SSLMode) error{
	connStr:= fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",host,port,user,password,databaseName,sslmode)
	DB,err:= sqlx.Open("postgres",connStr)

	if err!= nil{
		return err
	}
	err=DB.Ping()
	if err!=nil{
		return err
	}
	Ridesafe=DB
	return migrateUp(DB)
}
func ShutdownDatabase() error{

	if Ridesafe != nil {
    return Ridesafe.Close()
}
return nil
}

func migrateUp(db *sqlx.DB) error{
	driver, err:= postgres.WithInstance(db.DB, &postgres.Config{})

	if err!=nil{
		return err
	}
	m,err := migrate.NewWithDatabaseInstance("file://database/migrations","postgres",driver)

if err!=nil{
		return err
	}
	if err:= m.Up(); err!= nil && err!= migrate.ErrNoChange{
		return err
	}
	return nil
}
func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := Ridesafe.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %+v", err)
	}
	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				logrus.Errorf("failed to rollback tx: %s", rollBackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			logrus.Errorf("failed to commit tx: %s", commitErr)
		}
	}()
	err = fn(tx)
	return err
}