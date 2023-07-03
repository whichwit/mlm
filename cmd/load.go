/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
	mssql "github.com/microsoft/go-mssqldb"
	"github.com/spf13/cobra"
	"github.com/whichwit/mlm/mlm"
)

var db *sql.DB

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load path",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		if debug {
			fmt.Println("load called")
		}

		path := args[0]
		fileInfo, err := os.Stat(path)
		CheckErr(err)

		var files []string
		if fileInfo.IsDir() {
			files, err = filepath.Glob(path + "/*.mlm")
			CheckErr(err)
		} else if filepath.Ext(path) == ".mlm" {
			files = append(files, path)
		}

		if len(files) > 0 {
			// Create a new connector object by calling NewConnector
			connector, err := mssql.NewConnector(makeConnURL().String())
			CheckErr(err)

			// Use SessionInitSql to set any options that cannot be set with the dsn string
			// With ANSI_NULLS set to ON, compare NULL data with = NULL or <> NULL will return 0 rows
			connector.SessionInitSQL = "SET ANSI_NULLS ON"

			db = sql.OpenDB(connector)
			defer db.Close()

			for _, file := range files {
				LoadFromFile(file)
			}
		} else {
			log.Println("no file to process")
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
}

// load from a single file
func LoadFromFile(path string) {
	log.Println("load from file:", path)
	b, err := os.ReadFile(path)
	CheckErr(err)

	_mlm := mlm.New(string(b))

	usr, err := user.Current()
	CheckErr(err)

	username := usr.Username

	if _usr := strings.Split(usr.Username, `\`); len(_usr) > 1 {
		username = _usr[len(_usr)-1]
	}

	log.Println("user:", username)

	spew.Dump(_mlm.Name, _mlm.Title, _mlm.Arden, _mlm.Version, _mlm.Date, _mlm.Usage)
	log.Println("-------------")

	if err = db.Ping(); err != nil {
		log.Fatal("ping error:", err)
	}

}

// // read a file
// // then use regex to match the pattern
// // return matched results
// func parse(path string) {
// 	// read a file, then print to screen
// 	b, err := os.ReadFile(path)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	data := stripComments(string(b))

// 	// fmt.Printf(" stripComment:%s\n", data)

// 	title := parseFirst(data, `(?siU)maintenance:.*?title:(.*);;`)
// 	fmt.Printf("title: %q\n", title)
// 	name := parseFirst(data, `(?siU)maintenance:.*title:.*;;.*(?:mlmname|filename):(.*);;`)
// 	fmt.Printf("name: %q\n", name)
// 	arden := parseFirst(data, `(?siU)maintenance:.*title:.*;;.*mlmname:.*;;.*arden:(.*);;`)
// 	fmt.Printf("arden: %q\n", arden)
// 	date := parseFirst(data, `(?siU)maintenance:.*title:.*;;.*mlmname:.*;;.*arden:.*;;.*date:(.*);;`)
// 	fmt.Printf("date: %q\n", date)

// 	parseUsage(data)
// }

// // parse mlm title
// func parseTitle(s string) string {
// 	re := regexp.MustCompile(`(?siU)maintenance:.*?title:(.*);;`)
// 	match := re.FindStringSubmatch(s)
// 	if match != nil {
// 		return ""
// 	}
// 	return strings.Trim(match[1], " ")
// }

// // parse mlm name
// func parseName(s string) string {
// 	re := regexp.MustCompile(`(?siU)maintenance:.*title:.*;;.*(?:mlmname|filename):(.*);;`)
// 	match := re.FindStringSubmatch(s)
// 	if match == nil {
// 		return ""
// 	}
// 	name := match[1]
// 	return strings.ToUpper(strings.Trim(name, " "))
// }

// // parse first match of regex
// func parseFirst(text string, regex string) string {
// 	re := regexp.MustCompile(regex)
// 	match := re.FindStringSubmatch(text)
// 	if match != nil {
// 		return ""
// 	}
// 	return strings.Trim(match[1], " ")
// }

// // parse mlm usage
// func parseUsage(s string) int {
// 	// get evoke slot
// 	evoke_slot := regexp.MustCompile(`(?isU);;\s+evoke:(.*);;`).FindStringSubmatch(s)
// 	if evoke_slot == nil {
// 		return -1 // no evoke slot
// 	}

// 	evoked_events := (regexp.MustCompile((`(?isU)([^\s]+)\s*;`))).FindAllStringSubmatch(evoke_slot[1], -1)
// 	if evoked_events == nil {
// 		return 2 // no evoked events
// 	}

// 	for _, evoked_event := range evoked_events {
// 		activate_application_re := regexp.MustCompile(`(?isU)` + evoked_event[1] + `\s*:=\s*event\s*{\s*ActivateApplication User UserInfo.*}`)
// 		if activate_application_re.MatchString(s) {
// 			return 1 // 1 = activate application
// 		}
// 	}

// 	return 0 // event present but not activate application
// }

// // parse mlm title
// // func parseTitle(text string) string {
// // 	re := regexp.MustCompile(`(?siU)maintenance:.*?title:(.*);;`)
// // 	match := re.FindStringSubmatch(text)
// // 	if match == nil {
// // 		return ""
// // 	}
// // 	return match[1]
// // }

// // remove all c++-style comments from input using regex
// func stripComments(text string) string {
// 	return stripLineComments(stripBlockComments(text))
// }

// // remove all c++-style line comments from input using regex
// func stripLineComments(text string) string {
// 	// remove all c++-style comments from data using regex
// 	re := regexp.MustCompile(`//.*`)
// 	return re.ReplaceAllString(text, "")
// }

// // remove all c++-style block comments from input using regex
// func stripBlockComments(text string) string {
// 	// remove block level comments from text using regular expression
// 	re := regexp.MustCompile(`(?sU)(/\*.*\*/)`)
// 	return re.ReplaceAllString(text, "")
// }

// // fmt.Printf(" files:%s\n", files)

// // connString := fmt.Sprintf("server=%s;database=%s;port=%d;trusted+connection=true", server, database, port)

// // fmt.Printf(" connString:%s\n", connString)

// // iterate through a folder, read content of each file, then print to screen

// // read a file, then print to screen

// // and insert into the database

// // db, err := sql.Open("mssql", connString)
// // if err != nil {
// // 	log.Fatal("Open connection failed:", err.Error())
// // }
// // defer db.Close()

// // // insert
// // insertSql := "insert into test (name) values ('test')"
// // insert, err := db.Query(insertSql)
// // if err != nil {
// // 	log.Fatal("Insert failed:", err.Error())
// // }

// // defer insert.Close()

// // conn, err := sql.Open("mssql", connString)
// // if err != nil {
// // 	log.Fatal("Open connection failed:", err.Error())
// // }
// // defer conn.Close()
