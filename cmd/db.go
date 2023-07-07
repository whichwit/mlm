package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"

	mssql "github.com/microsoft/go-mssqldb"
)

// makeConnStr returns a URL struct so it may be modified by various
// tests before used as a DSN.
// func makeConnStr(host string, database string, instance string) *url.URL {
// 	dsn := os.Getenv("SQLSERVER_DSN")
// 	if len(dsn) > 0 {
// 		parsed, err := url.Parse(dsn)
// 		if err != nil {
// 			log.Fatal("unable to parse SQLSERVER_DSN as URL", err)
// 		}
// 		values := parsed.Query()
// 		if values.Get("log") == "" {
// 			values.Set("log", "127")
// 		}
// 		parsed.RawQuery = values.Encode()
// 		return parsed
// 	}
// 	values := url.Values{}
// 	values.Set("log", "127")
// 	values.Set("database", database)
// 	return &url.URL{
// 		Scheme:   "sqlserver",
// 		Host:     os.Getenv("HOST"),
// 		Path:     os.Getenv("INSTANCE"),
// 		User:     url.UserPassword(os.Getenv("SQLUSER"), os.Getenv("SQLPASSWORD")),
// 		RawQuery: values.Encode(),
// 	}
// }

func makeConnURL() *url.URL {
	query := url.Values{}
	query.Add("database", database)
	query.Add("TrustServerCertificate", "true")
	query.Add("app name", "mlm")

	return &url.URL{
		Scheme:   "sqlserver",
		Host:     server + ":" + strconv.Itoa(port),
		RawQuery: query.Encode(),
	}
}

func InitDb() {
	connString := makeConnURL().String()

	if debug {
		fmt.Printf(" connString:%s\n", connString)
	}

	// Create a new connector object by calling NewConnector
	connector, err := mssql.NewConnector(connString)
	if err != nil {
		log.Println(err)
		return
	}

	// Use SessionInitSql to set any options that cannot be set with the dsn string
	// With ANSI_NULLS set to ON, compare NULL data with = NULL or <> NULL will return 0 rows
	connector.SessionInitSQL = "SET ANSI_NULLS ON"

	db := sql.OpenDB(connector)
	defer db.Close()

	var result int
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// (*Row) Scan should copy data to bitval
	err = db.QueryRowContext(ctx, "select db_name()").Scan(&result)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("result: ", result)

	//   db, err := sql.Open("sqlserver", u.String())
	// connString := fmt.Sprintf("sqlserver://server=%s;user id=%s;password=%s;port=%d", server, port)
	// if *debug {
	// 	fmt.Printf(" connString:%s\n", connString)
	// }
	// conn, err := sql.Open("mssql", connString)
	// if err != nil {
	// 	log.Fatal("Open connection failed:", err.Error())
	// }
	// defer conn.Close()

	// stmt, err := conn.Prepare("select 1, 'abc'")
	// if err != nil {
	// 	log.Fatal("Prepare failed:", err.Error())
	// }
	// defer stmt.Close()

	// row := stmt.QueryRow()
	// var somenumber int64
	// var somechars string
	// err = row.Scan(&somenumber, &somechars)
	// if err != nil {
	// 	log.Fatal("Scan failed:", err.Error())
	// }
	// fmt.Printf("somenumber:%d\n", somenumber)
	// fmt.Printf("somechars:%s\n", somechars)

	// fmt.Printf("bye\n")
}
