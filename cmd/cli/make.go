package main

import (
	"errors"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
)

func doMake(arg2, arg3, arg4 string) error {

	switch arg2 {
	case "key":
		rnd := cel.RandomString(32)
		color.Yellow("32 character encryption key: %s", rnd)

	case "migration":
		err := checkFroDB()
		if err != nil {
			return err
		}

		if arg3 == "" {
			return errors.New("you must give the migration a name")
		}

		// default to migration type of fizz
		migrationType := "fizz"
		var up, down []byte

		// are doing fizz or sql?
		if arg4 == "fizz" || arg4 == "" {
			up, _ = templateFS.ReadFile("templates/migrations/migration_up.fizz")
			down, _ = templateFS.ReadFile("templates/migrations/migration_down.fizz")
		} else {
			migrationType = "sql"
		}

		// create the migrations for either fizz or sql
		err = cel.CreatePopMigration(up, down, arg3, migrationType)
		if err != nil {
			return err
		}

		// ---------------
		// golang-migrate
		// ---------------

		// fileName := fmt.Sprintf("%d_%s", time.Now().UnixMicro(), arg3)
		// upFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".up.sql"
		// downFile := cel.RootPath + "/migrations/" + fileName + "." + dbType + ".down.sql"

		// err := copyFileFromTemplate("templates/migrations/migration."+dbType+".up.sql", upFile)
		// if err != nil {
		// 	return err
		// }

		// err = copyFileFromTemplate("templates/migrations/migration."+dbType+".down.sql", downFile)
		// if err != nil {
		// 	return err
		// }

	case "auth":
		err := doAuth()
		if err != nil {
			return err
		}

	case "handler":
		if arg3 == "" {
			return errors.New("you must give the handler a name")
		}

		fileName := cel.RootPath + "/handlers/" + strings.ToLower(arg3) + ".go"
		if fileExists(fileName) {
			return errors.New(fileName + " already exists")
		}

		data, err := templateFS.ReadFile("templates/handlers/handler.go.txt")
		if err != nil {
			return err
		}

		handler := string(data)
		handler = strings.ReplaceAll(handler, "$HANDLER_NAME$", strcase.ToCamel(arg3))

		err = os.WriteFile(fileName, []byte(handler), 0644)
		if err != nil {
			return err
		}

	case "model":
		if arg3 == "" {
			return errors.New("you must give the model a name")
		}

		data, err := templateFS.ReadFile("templates/data/model.go.txt")
		if err != nil {
			return err
		}

		model := string(data)
		plur := pluralize.NewClient()

		var modelName = arg3
		var tableName = arg3

		if plur.IsPlural(arg3) {
			modelName = plur.Singular(arg3)
			tableName = strings.ToLower(tableName)
		} else {
			tableName = strings.ToLower(tableName)
		}

		fileName := cel.RootPath + "/data/" + strings.ToLower(modelName) + ".go"
		if fileExists(fileName) {
			return errors.New(fileName + " already exists")
		}

		model = strings.ReplaceAll(model, "$MODEL_NAME$", strcase.ToCamel(modelName))
		model = strings.ReplaceAll(model, "$TABLE_NAME$", tableName)

		err = copyDataToFile([]byte(model), fileName)
		if err != nil {
			return err
		}

	case "session":
		err := doSessionTable()
		if err != nil {
			return err
		}

	case "mail":
		if arg3 == "" {
			return errors.New("you must give the mail template a name")
		}
		htmlMail := cel.RootPath + "/mail/" + strings.ToLower(arg3) + ".html.tmpl"
		plainMail := cel.RootPath + "/mail/" + strings.ToLower(arg3) + ".plain.tmpl"

		err := copyFileFromTemplate("templates/mailer/mail.html.tmpl", htmlMail)
		if err != nil {
			return err
		}

		err = copyFileFromTemplate("templates/mailer/mail.plain.tmpl", plainMail)
		if err != nil {
			return err
		}
	}

	return nil
}
