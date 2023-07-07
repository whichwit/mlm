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

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		connector, err := mssql.NewConnector(makeConnURL().String())
		CheckErr(err)

		// Use SessionInitSql to set any options that cannot be set with the dsn string
		// With ANSI_NULLS set to ON, compare NULL data with = NULL or <> NULL will return 0 rows
		connector.SessionInitSQL = "SET ANSI_NULLS ON"

		db = sql.OpenDB(connector)
		defer db.Close()

		err = db.Ping()
		CheckErr(err)

		log.Println("Connected!", makeConnURL().String())
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// validateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// validateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
