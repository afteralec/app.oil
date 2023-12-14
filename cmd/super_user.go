/*
Copyright © 2023 Alec DuBois <alec@petrichormud.com>
*/
package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"petrichormud.com/app/internal/password"
	"petrichormud.com/app/internal/queries"
	"petrichormud.com/app/internal/username"
)

var superUserCmd = &cobra.Command{
	Use:   "super-user",
	Short: "Seed the database with a super-user.",
	Long:  `Seed the database with a super-user.`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		u, err := cmd.Flags().GetString("username")
		if err != nil {
			return err
		}
		pw, err := cmd.Flags().GetString("password")
		if err != nil {
			return err
		}
		dbURL, err := cmd.Flags().GetString("db-url")
		if err != nil {
			return err
		}

		if !username.IsValid(u) {
			return errors.New("please enter a valid username")
		}

		if !password.IsValid(pw) {
			return errors.New("please enter a valid password")
		}

		db, err := sql.Open("mysql", dbURL)
		if err != nil {
			return err
		}

		// TODO: Move this database setup to its own function
		_, err = db.Exec("SET GLOBAL local_infile=true;")
		if err != nil {
			return err
		}

		if err = db.Ping(); err != nil {
			return errors.New("error while pinging db")
		}

		q := queries.New(db)
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()
		qtx := q.WithTx(tx)

		pwHash, err := password.Hash(pw)
		if err != nil {
			return errors.New("error while hashing password")
		}

		_, err = qtx.CreatePlayer(context.Background(), queries.CreatePlayerParams{
			Username: u,
			PwHash:   pwHash,
		})
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("super-user called with %s and %s.", u, pw)
		fmt.Println(msg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(superUserCmd)

	superUserCmd.Flags().StringP("db-url", "d", "root:pass@127.0.0.1/test", "The URL for the database.")
	superUserCmd.Flags().StringP("username", "u", "", "The username for the new user.")
	superUserCmd.Flags().StringP("password", "p", "", "The password for the user.")
	superUserCmd.MarkFlagRequired("username")
	superUserCmd.MarkFlagRequired("password")
}
