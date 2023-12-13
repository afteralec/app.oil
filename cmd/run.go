/*
Copyright © 2023 Alec DuBois <alec@petrichormud.com>
*/
package cmd

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	fiber "github.com/gofiber/fiber/v2"
	html "github.com/gofiber/template/html/v2"
	"github.com/spf13/cobra"

	"petrichormud.com/app/internal/app"
	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/shared"
	"petrichormud.com/app/web"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the application",
	Long:  `Run the application`,
	Run: func(_ *cobra.Command, _ []string) {
		i := shared.SetupInterfaces()
		defer i.Close()

		views := html.NewFileSystem(http.FS(web.ViewsFS), ".html")
		a := fiber.New(configs.Fiber(views))

		app.Middleware(a, &i)
		app.Handlers(a, &i)
		app.Static(a)

		log.Fatal(a.Listen(":8008"))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
