package main

import (
	"fmt"
	"time"
)

func doAuth() error {
	// migrations
	dbType := cel.DB.DataType
	fileName := fmt.Sprintf("%d_create_auth_tables", time.Now().UnixMicro())
	upFile := cel.RootPath + "/migrations/" + fileName + ".up.sql"
	downFile := cel.RootPath + "/migrations/" + fileName + ".down.sql"

	err := copyFileFromTemplate("templates/migrations/auth_tables."+dbType+".sql", upFile)
	if err != nil {
		return err
	}

	err = copyDataToFile([]byte("drop table if exists users cascade;\ndrop table if exists tokens cascade;\ndrop table if exists remember_tokens cascade;"), downFile)
	if err != nil {
		return err
	}

	// run migration
	err = doMigrate("up", "")
	if err != nil {
		return err
	}

	// copy files
	err = copyFileFromTemplate("templates/data/user.go.txt", cel.RootPath+"/data/user.go")
	if err != nil {
		return err
	}

	err = copyFileFromTemplate("templates/data/token.go.txt", cel.RootPath+"/data/token.go")
	if err != nil {
		return err
	}

	return nil
}
