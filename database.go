package noaweb

import (
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq" // Importing postgres driver
)

type dbConnections map[string]*sql.DB

// ConnectDB opens a connection to a database. The database connection is stored
// in the main noaweb instance. After connection is open the connection can be retrieved
// by using GetDatabaseConnection(name), where name is on the form user@database.
// connStr is a database connection string such as:
// postgres://user:password@localhost:port/dbname?sslmode=disable
func (i *Instance) ConnectDB(connStr string) error {
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	dbConn.Ping()
	if err != nil {
		return err
	}

	name, _ := connStrToReference(connStr)
	noawebinst.dbConnections[name] = dbConn

	return nil
}

// DisconnectDB closes the database connection. Normally called with a defer right
// after ConnectDB.
func (i *Instance) DisconnectDB(name string) error {
	if val, ok := noawebinst.dbConnections[name]; ok {
		val.Close()
		return nil
	}
	return errors.New("Noaweb.DissconnectDB: Could not find given connection")
}

// GetDatabaseConnection returns a database object which can be used to run
// SQL statements.
func GetDatabaseConnection(name string) (*sql.DB, error) {
	if val, ok := noawebinst.dbConnections[name]; ok {
		return val, nil
	}
	return nil, errors.New("GetDatabaseConnection: Could not find given connection")
}

// Splits name to user and database.
func connStrToReference(connStr string) (string, error) {
	var str []string
	str = strings.Split(connStr, "//")
	str = strings.Split(str[1], ":")
	user := str[0]
	str = strings.Split(str[2], "/")
	str = strings.Split(str[1], "?")
	db := str[0]

	return user + "@" + db, nil
}

// MigrateFolder runs all migrations in a given folder to database given by name.
// Action can be up, down, refresh or false.
// TODO: Add flags
func (i *Instance) MigrateFolder(migrationsPath, name string) error {
	action := i.commandFlags.migrateAction
	fmt.Println("Migrating from folder: " + migrationsPath + " using action: " + action + ".")
	if _, ok := noawebinst.dbConnections[name]; ok {
		if action == "false" {
			return nil
		}

		if !strings.HasPrefix(migrationsPath, "/") {
			migrationsPath = noawebinst.AssetsDir + "/" + migrationsPath
		}

		var stmts []string
		var dat []byte

		walkFn := func(path string, info os.FileInfo, err error) error {
			if action == "down" || action == "refresh" {
				if filepath.Ext(path) == ".drop" {
					dat, err = ioutil.ReadFile(path)
					stmt := string(dat)
					fmt.Println("Adding migration: " + path)
					if err != nil {
						return err
					}
					stmts = append([]string{stmt}, stmts...)
				}
			}
			if action == "up" || action == "refresh" {
				if filepath.Ext(path) == ".create" || filepath.Ext(path) == ".alter" {
					dat, err = ioutil.ReadFile(path)
					stmt := string(dat)
					fmt.Println("Adding migration: " + path)
					if err != nil {
						return err
					}
					stmts = append(stmts, stmt)
				}
			}
			return nil
		}
		err := filepath.Walk(migrationsPath, walkFn)
		if err != nil {
			return err
		}

		for _, stmt := range stmts {

			if err := i.MigrateDB(name, stmt); err != nil {
				return err
			}
		}

		if action == "down" {
			fmt.Println("Program must be run with valid database, exiting..")
			os.Exit(1)
		}

		return nil
	}
	return errors.New("Noaweb.MigrateDB: Could not find given connection")

}

// MigrateDB migrates the database by running the stmtStr query.
func (i *Instance) MigrateDB(name, stmtStr string) error {
	if val, ok := noawebinst.dbConnections[name]; ok {
		stmt, err := val.Prepare(stmtStr)
		if err != nil {
			return err
		}
		_, err = stmt.Exec()
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("Noaweb.MigrateDB: Could not find given connection")
}

// DatabaseFunctions struct used to structure package
type DatabaseFunctions struct{}

// Database variable used to structure package
var Database DatabaseFunctions

// Get returns the database object which can be used when doing database stuff.
func (DatabaseFunctions) Get(name string) (*sql.DB, error) {
	if val, ok := noawebinst.dbConnections[name]; ok {
		return val, nil
	}
	return nil, errors.New("Noaweb.GetDB: Could not find given connection")
}
