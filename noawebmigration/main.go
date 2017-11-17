package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	migrationType := flag.String("type", "", "Which file to be created ( .create | .drop | .alter )")
	flag.Parse()

	file := time.Now().Format("20060102_150405_") + strings.Join(os.Args[1:], "_")
	if strings.HasSuffix(os.Args[1], "/") {
		file = os.Args[1] + time.Now().Format("20060102_150405_") + strings.Join(os.Args[2:], "_")
	}

	switch *migrationType {
	case "create":
		fmt.Println("Creating " + file + ".create")
		createMigration(file, "create")
	case "drop":
		fmt.Println("Creating " + file + ".drop")
		createMigration(file, "drop")
	case "alter":
		fmt.Println("Creating " + file + ".alter")
		createMigration(file, "alter")
	case "":
		fmt.Println("Creating " + file + ".create / .drop")
		createMigration(file, "create")
		createMigration(file, "drop")
	default:
		fmt.Printf("Unknown migration type, see --help for more information.")
	}
}

func createMigration(filepath, ext string) error {
	os.OpenFile(filepath+"."+ext, os.O_RDONLY|os.O_CREATE, 0666)
	return nil
}
