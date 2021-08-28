package postgres

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	rr "github.com/artofey/media-tel-test"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var dbCreate string = `
DO
$do$
BEGIN
   IF EXISTS (SELECT FROM pg_database WHERE datname = '%s') THEN
      RAISE NOTICE 'Database already exists';  -- optional
   ELSE
      PERFORM dblink_exec('dbname=' || current_database()  -- current db
                        , 'CREATE DATABASE %s');
   END IF;
END
$do$;
`

type PostgresDB struct {
	db *sqlx.DB
}

func New(ip, port, dbName, user, passwd string) (*PostgresDB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s sslmode=disable",
		ip, port, user, passwd)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("New Error: %v", err)
	}
	db.MustExec(fmt.Sprintf(dbCreate, dbName, dbName))
	db.Close()
	db, err = sqlx.Connect("postgres", dsn+fmt.Sprintf(" dbname=%s", dbName))
	if err != nil {
		return nil, fmt.Errorf("New Error: %v", err)
	}
	return &PostgresDB{
		db: db,
	}, nil
}

func (db *PostgresDB) CreateTable(name string, records rr.Records) error {
	cr := "CREATE TABLE IF NOT EXISTS %s (%s);"
	var columns []string
	for _, h := range records.Headers {
		columns = append(columns, fmt.Sprintf("\"%s\" text", h))
	}
	fields := strings.Join(columns, ",")
	query := fmt.Sprintf(cr, strings.Split(name, ".")[0], fields)
	// fmt.Println(query)
	_, err := db.db.Exec(query)
	if err != nil {
		return fmt.Errorf("exec query: \n%v\nerror: %v", query, err)
	}
	return nil
}

func (db *PostgresDB) InsertRecordsToTable(tableName string, records rr.Records) error {
	log.Println("start insert records", records)
	query1 := fmt.Sprintf("INSERT INTO %s ", tableName)
	newHeaders := make([]string, 0, len(records.Headers))
	for _, h := range records.Headers {
		newHeaders = append(newHeaders, fmt.Sprintf("\"%s\"", h))
	}
	query2 := fmt.Sprintf("( %s ) VALUES ", strings.Join(newHeaders, ","))

	tx := db.db.MustBegin()
	for _, rec := range records.R {
		var vals []string
		for i := range rec {
			vals = append(vals, fmt.Sprintf("$%s", strconv.Itoa(i+1)))
		}
		query3 := fmt.Sprintf("( %s )", strings.Join(vals, ","))
		resQuery := query1 + query2 + query3
		fmt.Println(resQuery)
		s := make([]interface{}, len(rec))
		for i, v := range rec {
			s[i] = v
		}
		_, err := tx.Exec(resQuery, s...)
		if err != nil {
			return fmt.Errorf("exec query: \n%v\nerror: %v", resQuery, err)
		}

	}
	return tx.Commit()
}
