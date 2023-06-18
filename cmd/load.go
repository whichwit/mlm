/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"github.com/whichwit/mlm/mlm"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("load called")

		workingDirectory, err := os.Getwd() // get current working directory
		if err != nil {
			log.Fatal(err)
		}
		CheckErr(err)

		fmt.Printf(" dir:%s\n", workingDirectory)

		files, err := filepath.Glob(workingDirectory + "/../ascension/MISAG/mlm" + "/*.mlm")

		if err != nil {
			log.Fatal(err)
		}

		// for _, file := range files {
		// 	// fmt.Printf(" file:%s\n", file)
		// 	if strings.HasSuffix(file, "UTIL_STRING_PARSE.mlm") || strings.HasSuffix(file, "STD_MOBILE_ALERT_NOTIFICATION.mlm") || strings.HasSuffix(file, `SM_TEST_PATIENT_MISMATCH.mlm`) {

		// 		parse(file)
		// 		print("\n")
		// 	}
		// }

		// fmt.Println("------------------")
		for _, file := range files {
			LoadFromFile(file)

		}

	},
}

func init() {
	rootCmd.AddCommand(loadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadCmd.Flags().IntVarP(&port, "port", "", 1433, "the database port")
	// loadCmd.Flags().StringVarP(&server, "server", "", "localhost", "the database server")
	// loadCmd.Flags().StringVarP(&database, "database", "", "localhost", "the database name")

}

// load from a single file
func LoadFromFile(path string) {
	b, err := os.ReadFile(path)
	CheckErr(err)

	_mlm := mlm.New(string(b))

	// name := parseName(data)
	// usage := parseUsage(data)
	// fmt.Printf("%d %s  %s  %s\n", e.Usage, e.Title, e.Name, file)
	spew.Dump(_mlm.Usage, _mlm.Title, _mlm.Name, path, "-------------")
}

// read a file
// then use regex to match the pattern
// return matched results
func parse(path string) {
	// read a file, then print to screen
	b, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	data := stripComments(string(b))

	// fmt.Printf(" stripComment:%s\n", data)

	title := parseFirst(data, `(?siU)maintenance:.*?title:(.*);;`)
	fmt.Printf("title: %q\n", title)
	name := parseFirst(data, `(?siU)maintenance:.*title:.*;;.*(?:mlmname|filename):(.*);;`)
	fmt.Printf("name: %q\n", name)
	arden := parseFirst(data, `(?siU)maintenance:.*title:.*;;.*mlmname:.*;;.*arden:(.*);;`)
	fmt.Printf("arden: %q\n", arden)
	date := parseFirst(data, `(?siU)maintenance:.*title:.*;;.*mlmname:.*;;.*arden:.*;;.*date:(.*);;`)
	fmt.Printf("date: %q\n", date)

	parseUsage(data)
}

// parse mlm title
func parseTitle(s string) string {
	re := regexp.MustCompile(`(?siU)maintenance:.*?title:(.*);;`)
	match := re.FindStringSubmatch(s)
	if match != nil {
		return ""
	}
	return strings.Trim(match[1], " ")
}

// parse mlm name
func parseName(s string) string {
	re := regexp.MustCompile(`(?siU)maintenance:.*title:.*;;.*(?:mlmname|filename):(.*);;`)
	match := re.FindStringSubmatch(s)
	if match == nil {
		return ""
	}
	name := match[1]
	return strings.ToUpper(strings.Trim(name, " "))
}

// parse first match of regex
func parseFirst(text string, regex string) string {
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(text)
	if match != nil {
		return ""
	}
	return strings.Trim(match[1], " ")
}

// parse mlm usage
func parseUsage(s string) int {
	// get evoke slot
	evoke_slot := regexp.MustCompile(`(?isU);;\s+evoke:(.*);;`).FindStringSubmatch(s)
	if evoke_slot == nil {
		return -1 // no evoke slot
	}

	evoked_events := (regexp.MustCompile((`(?isU)([^\s]+)\s*;`))).FindAllStringSubmatch(evoke_slot[1], -1)
	if evoked_events == nil {
		return 2 // no evoked events
	}

	for _, evoked_event := range evoked_events {
		activate_application_re := regexp.MustCompile(`(?isU)` + evoked_event[1] + `\s*:=\s*event\s*{\s*ActivateApplication User UserInfo.*}`)
		if activate_application_re.MatchString(s) {
			return 1 // 1 = activate application
		}
	}

	return 0 // event present but not activate application
}

// parse mlm title
// func parseTitle(text string) string {
// 	re := regexp.MustCompile(`(?siU)maintenance:.*?title:(.*);;`)
// 	match := re.FindStringSubmatch(text)
// 	if match == nil {
// 		return ""
// 	}
// 	return match[1]
// }

// remove all c++-style comments from input using regex
func stripComments(text string) string {
	return stripLineComments(stripBlockComments(text))
}

// remove all c++-style line comments from input using regex
func stripLineComments(text string) string {
	// remove all c++-style comments from data using regex
	re := regexp.MustCompile(`//.*`)
	return re.ReplaceAllString(text, "")
}

// remove all c++-style block comments from input using regex
func stripBlockComments(text string) string {
	// remove block level comments from text using regular expression
	re := regexp.MustCompile(`(?sU)(/\*.*\*/)`)
	return re.ReplaceAllString(text, "")
}

// fmt.Printf(" files:%s\n", files)

// connString := fmt.Sprintf("server=%s;database=%s;port=%d;trusted+connection=true", server, database, port)

// fmt.Printf(" connString:%s\n", connString)

// iterate through a folder, read content of each file, then print to screen

// read a file, then print to screen

// and insert into the database

// db, err := sql.Open("mssql", connString)
// if err != nil {
// 	log.Fatal("Open connection failed:", err.Error())
// }
// defer db.Close()

// // insert
// insertSql := "insert into test (name) values ('test')"
// insert, err := db.Query(insertSql)
// if err != nil {
// 	log.Fatal("Insert failed:", err.Error())
// }

// defer insert.Close()

// conn, err := sql.Open("mssql", connString)
// if err != nil {
// 	log.Fatal("Open connection failed:", err.Error())
// }
// defer conn.Close()
