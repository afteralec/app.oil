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
	"petrichormud.com/app/internal/permission"
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

		// TODO: Move this database setup to its own function so it's repeatable
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

		result, err := qtx.CreatePlayer(context.Background(), queries.CreatePlayerParams{
			Username: u,
			PwHash:   pwHash,
		})
		if err != nil {
			return err
		}

		pid, err := result.LastInsertId()
		if err != nil {
			return err
		}

		if err := qtx.CreatePlayerPermissionIssuedChangeHistory(context.Background(), queries.CreatePlayerPermissionIssuedChangeHistoryParams{
			PID:        pid,
			IPID:       pid,
			Permission: permission.PlayerAssignAllPermissions,
		}); err != nil {
			return err
		}

		if err := qtx.CreatePlayerPermission(context.Background(), queries.CreatePlayerPermissionParams{
			PID:        pid,
			IPID:       pid,
			Permission: permission.PlayerAssignAllPermissions,
		}); err != nil {
			return err
		}

		if err = tx.Commit(); err != nil {
			return err
		}

		msg := fmt.Sprintf("User %s created and seeded with the root permission.", u)
		fmt.Println(msg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(superUserCmd)

	superUserCmd.Flags().StringP("db-url", "d", "root:pass@/test", "The URL for the database.")
	superUserCmd.Flags().StringP("username", "u", "", "The username for the new user.")
	superUserCmd.Flags().StringP("password", "p", "", "The password for the user.")
	superUserCmd.MarkFlagRequired("username")
	superUserCmd.MarkFlagRequired("password")
}
