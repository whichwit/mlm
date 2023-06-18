package cmd

import (
	"fmt"
	"net/url"

	_ "github.com/microsoft/go-mssqldb"
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

func Init() {
	query := url.Values{}
	query.Add("app name", "MyAppName")

	u := &url.URL{
		Scheme: "sqlserver",
		User:   nil,
		Host:   fmt.Sprintf("%s:%d", server, port),
		// Path:  instance, // if connecting to an instance instead of a port
		RawQuery: query.Encode(),
	}
	fmt.Println(u.String())
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
