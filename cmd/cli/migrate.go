package main

// func doMigrate(arg2, arg3 string) error {
// 	dsn := getDSN()

// 	// run the migration command
// 	switch arg2 {
// 	case "up":
// 		err := cel.MigrateUp(dsn)
// 		if err != nil {
// 			return err
// 		}
// 	case "down":
// 		if arg3 == "all" {
// 			err := cel.MigrateDown(dsn)
// 			if err != nil {
// 				return err
// 			}
// 		} else {
// 			err := cel.Steps(-1, dsn)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	case "reset":
// 		err := cel.MigrateDown(dsn)
// 		if err != nil {
// 			return err
// 		}

// 		err = cel.MigrateUp(dsn)
// 		if err != nil {
// 			return err
// 		}
// 	default:
// 		showHelp()
// 	}

// 	return nil
// }

func doMigrate(arg2, arg3 string) error {
	err := checkFroDB()
	if err != nil {
		return err
	}

	tx, err := cel.PopConnect()
	if err != nil {
		return err
	}
	defer tx.Close()

	// run the migration command
	switch arg2 {
	case "up":
		err := cel.PopMigrationUp(tx)
		if err != nil {
			return err
		}
	case "down":
		if arg3 == "all" {
			err := cel.PopMigrateDown(tx, -1)
			if err != nil {
				return err
			}
		} else {
			err := cel.PopMigrateDown(tx, 1)
			if err != nil {
				return err
			}
		}
	case "reset":
		err := cel.PopMigrationReset(tx)
		if err != nil {
			return err
		}

	case "status":
		err := cel.PopMigrationStatus(tx)
		if err != nil {
			return err
		}
	default:
		showHelp()
	}

	return nil
}
