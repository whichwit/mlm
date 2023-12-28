/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"database/sql"
	"log"

	mssql "github.com/microsoft/go-mssqldb"
	"github.com/spf13/cobra"
)

var outdir string

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

		q := "SELECT TOP 2 Name, Logic FROM CV3MLM WHERE Active = 1 AND Status = 4"

		rows, err := db.Query(q)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var (
				name  string
				logic string
			)
			if err := rows.Scan(&name, &logic); err != nil {
				log.Fatal(err)
			}
			log.Printf("id %s name is %s\n", name, logic)
		}
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)

	extractCmd.Flags().StringVarP(&outdir, "output-dir", "o", "mlm", "Output directory")
}
