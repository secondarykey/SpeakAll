package db

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"strings"
)

const schemaVersion = "0.2"

type schemaError struct {
	code    int
	message string
}

type RecordData *string
type Record []RecordData
type FlexRows struct {
	Columns []string
	Records []Record
}

func (s *schemaError) Error() string {
	return fmt.Sprintf("%d:%s", s.code, s.message)
}

func NewSchemaError(code int, msg string) *schemaError {
	return &schemaError{
		code:    code,
		message: msg,
	}
}

var inst *sql.DB

func checkSchemaVersion(path, ver string) (string, *schemaError) {

	//%s be
	pArr := strings.Split(path, "%s")
	if len(pArr) != 2 {
		return "", NewSchemaError(-1, "Error:database path is '%s' requid["+path+"]")
	}

	rpath := fmt.Sprintf(path, schemaVersion)
	//exist database file
	_, err := os.Stat(rpath)
	//call version check
	if ver == schemaVersion || ver == "test" {
		if err == nil {
			return rpath, nil
		}
		return rpath, NewSchemaError(0, "Create database")
	}

	if err == nil {
		return rpath, nil
	}
	//code 0
	return rpath, NewSchemaError(0, "Warning:Program version,TOML file version")
}

func Listen(path, version string) error {

	var err error
	rp, scErr := checkSchemaVersion(path, version)

	cFlag := true
	if scErr != nil {
		if scErr.code == 0 {
			log.Println(scErr.Error() + "[" + path + "]")
			cFlag = false
		} else {
			return scErr
		}
	}

	inst, err = sql.Open("sqlite3", rp)
	if err != nil {
		return err
	}

	if cFlag {
		return nil
	}

	err = deleteTables()
	if err != nil {
		return err
	}

	err = createInitTables()
	if err != nil {
		return err
	}

	return insertInitTable()
}

func createInitTables() error {
	err := createUserTable()
	if err != nil {
		return err
	}
	err = createRoleTable()
	if err != nil {
		return err
	}
	err = createUserRoleTable()
	if err != nil {
		return err
	}
	err = createMessageTable()
	if err != nil {
		return err
	}
	err = createCategoryTable()
	if err != nil {
		return err
	}
	err = createMemoTable()
	if err != nil {
		return err
	}
	return nil
}

func insertInitTable() error {

	tx, err := inst.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	pwd := CreateMD5("password")
	rslt, err := InsertUser(tx, "SpeakAll Administrator", "admin@localhost", pwd)
	if err != nil {
		return err
	}
	userId, _ := rslt.LastInsertId()

	rslt, err = insertRole(tx, "Administrator", "Admin")
	if err != nil {
		return err
	}
	rslt, err = insertRole(tx, "Category Chairman", "Chairman")
	if err != nil {
		return err
	}

	rslt, err = insertRole(tx, "Speaker", "Speaker")
	if err != nil {
		return err
	}

	rslt, err = InsertUserRole(tx, int(userId), "Admin")
	if err != nil {
		return err
	}

	rslt, err = InsertUserRole(tx, int(userId), "Chairman")
	if err != nil {
		return err
	}

	rslt, err = InsertUserRole(tx, int(userId), "Speaker")
	if err != nil {
		return err
	}

	return tx.Commit()
}

func deleteTables() error {
	err := deleteUserTable()
	if err != nil {
		return err
	}
	err = deleteRoleTable()
	if err != nil {
		return err
	}
	err = deleteUserRoleTable()
	if err != nil {
		return err
	}
	err = deleteMessageTable()
	if err != nil {
		return err
	}
	err = deleteCategoryTable()
	if err != nil {
		return err
	}
	err = deleteMemoTable()
	if err != nil {
		return err
	}
	return nil
}

func Exec(sql string) (sql.Result, error) {
	return inst.Exec(sql)
}

func Query(sql string) (*FlexRows, error) {

	rows, err := inst.Query(sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	flexRows := &FlexRows{}
	flexRows.Columns, err = rows.Columns()
	if err != nil {
		return nil, err
	}

	records := make([]Record, 0)

	for rows.Next() {
		rec := make([]RecordData, len(flexRows.Columns))
		for i, _ := range rec {
			rec[i] = toPtr("")
		}
		rows.Scan(rec)
		records = append(records, rec)
	}

	flexRows.Records = records
	log.Println(records)

	return flexRows, err
}

func toPtr(s string) *string {
	return &s
}
