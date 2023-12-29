/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/spf13/cobra"
)

var outdir string
var q string

// extractCmd represents the extract command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extract MLM from the database",

	Run: func(cmd *cobra.Command, args []string) {
		connector, err := mssql.NewConnector(makeConnURL().String())
		CheckErr(err)

		// Use SessionInitSql to set any options that cannot be set with the dsn string
		// With ANSI_NULLS set to ON, compare NULL data with = NULL or <> NULL will return 0 rows
		connector.SessionInitSQL = "SET ANSI_NULLS ON"

		db = sql.OpenDB(connector)
		defer db.Close()

		rows, err := db.Query(q)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		if !rows.Next() {
			log.Println("No records found")
			return
		}

		// create out directory if it doesn't exist
		if _, err := os.Stat(outdir); os.IsNotExist(err) {
			err = os.MkdirAll(outdir, 0755)
			if err != nil {
				log.Fatal(err)
			}
		}

		for rows.Next() {
			var (
				name  string
				logic string
			)
			if err := rows.Scan(&name, &logic); err != nil {
				log.Fatal(err)
			}
			logic = strings.Replace(logic, "{{{SINGLE-QUOTE}}}", "'", -1)

			// write content of logic to file
			file, err := os.Create(fmt.Sprintf("%s/%s.mlm", outdir, name))
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			_, err = file.WriteString(logic)
			if err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)

	extractCmd.Flags().StringVarP(&outdir, "output-dir", "o", "mlm", "Output directory")
	extractCmd.Flags().StringVarP(&q, "query", "q", "SELECT Name, Logic FROM CV3MLM WHERE Active = 1 AND Status = 4", "Extraction query")

}
