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
	"sync"

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

			var wg sync.WaitGroup
			wg.Add(len(files))
			for _, file := range files {
				go func(file string) {
					LoadFromFile(file)
					defer wg.Done()
				}(file)
			}
			wg.Wait()
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
	log.Println("Processing:", path)
	b, err := os.ReadFile(path)
	CheckErr(err)

	_mlm := mlm.New(string(b))

	usr, err := user.Current()
	CheckErr(err)
	username := usr.Username
	if _usr := strings.Split(usr.Username, `\`); len(_usr) > 1 {
		username = _usr[len(_usr)-1]
	}

	// log.Println("user:", username)

	// spew.Dump(_mlm.Name, _mlm.Title, _mlm.Arden, _mlm.Version, _mlm.Date, _mlm.Usage)
	// spew.Dump(_mlm)

	if err = db.Ping(); err != nil {
		log.Fatal("ping error:", err)
	}

	/* reference SP
	ALTER PROC [dbo].[SCMMlmInsPr]

	(
	    @UserID		 VARCHAR(30),
	    @Name        varchar(80),
	    @Description varchar(255),
	    @Logic       varchar(Max),
	    @Status      int,
	    @UsageType   int = 0,
	    @Title       varchar(255) = null,
	    @ArdenVersion varchar(5) = null,
	    @MLMVersion  varchar(5) = null,
	    @MLMDate     datetime = null
	)

	AS
	*/

	// insert mlm
	var rs mssql.ReturnStatus
	db.Exec(
		"SCMMlmInsPr",
		sql.Named("UserID", username),
		sql.Named("Name", _mlm.Name),
		sql.Named("Title", _mlm.Title),
		sql.Named("Description", ""),
		sql.Named("Logic", _mlm.Content),
		sql.Named("Status", 4),
		sql.Named("UsageType", _mlm.Usage),
		sql.Named("ArdenVersion", _mlm.Arden),
		sql.Named("MLMVersion", _mlm.Version),
		sql.Named("MLMDate", _mlm.Date),
		&rs,
	)
	log.Println("SCMMlmInsPr return status:", rs)
	if rs != 0 {
		log.Fatal("SCMMlmInsPr error:", _mlm)
	}
	log.Println("Done:", path)
	log.Println("-------------")
}
