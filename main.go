package main

import (
	"log"
	"os"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/urfave/cli"

	"github.com/WhoSV/codestack-api/database"
	"github.com/WhoSV/codestack-api/model"
	"github.com/WhoSV/codestack-api/repository"
	"github.com/WhoSV/codestack-api/server"
)

func main() {
	app := cli.NewApp()
	app.Name = "odestack-api"
	app.Usage = "CodeStack API Server"
	app.Version = "1.2.0"
	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Start API server.",
			Action: func(c *cli.Context) {
				println("action:", "run")

				// Run the App here
				server.ExecuteServer()
			},

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config-path",
					Value: "config.json",
					Usage: "Path to the config file",
				},
			},
		},
		{
			Name:    "createadminuser",
			Aliases: []string{"c"},
			Usage:   "Create a new admin user.",
			Action: func(c *cli.Context) {
				println("action:", "create admin user")

				email := c.String("email")
				password := c.String("password")
				role := c.String("admin")
				fullName := c.String("full_name")

				// Validate required options.
				if len(email) == 0 {
					log.Fatal("Email can not be empty.")
				}
				if len(password) == 0 {
					log.Fatal("Password can not be empty.")
				}
				// set default value if empty
				if len(role) == 0 {
					role = model.RoleAdmin
				}
				if len(fullName) == 0 {
					log.Fatal("fullName can not be empty.")
				}

				database.ConnectDB()
				defer database.CloseDB()
				person, err := repository.CreateNewPerson(email, password, role, fullName)
				if err != nil {
					log.Fatal("Can not create user:", err.Error())
					return
				}
				log.Println("User:", person.Email, "created successfuly!")

			},

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "config-path",
					Value: "config.json",
					Usage: "Path to the config file",
				},
				cli.StringFlag{
					Name:  "email",
					Value: "",
					Usage: "User's email address",
				},
				cli.StringFlag{
					Name:  "password",
					Value: "",
					Usage: "User's password",
				},
				cli.StringFlag{
					Name:  "role",
					Value: "",
					Usage: "User's role",
				},
				cli.StringFlag{
					Name:  "full_name",
					Value: "",
					Usage: "User's full name",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
