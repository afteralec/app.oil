/*
Copyright © 2023 Alec DuBois <alec@petrichormud.com>
*/
package cmd

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/config"
	"petrichormud.com/app/internal/service"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	Long:  `Run the application`,
	Run: func(_ *cobra.Command, _ []string) {
		i := service.NewInterfaces()
		defer i.Close()

		a := fiber.New(config.Fiber(i.Templates))

		app.Middleware(a, &i)
		app.Handlers(a, &i)
		app.Static(a)

		log.Fatal(a.Listen(":8008"))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
