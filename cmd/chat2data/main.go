package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/byebyebruce/chat2data/cmd"
	"github.com/byebyebruce/chat2data/datachain"
	"github.com/joho/godotenv"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/tools/sqldatabase/mysql"
	"github.com/tmc/langchaingo/tools/sqldatabase/postgresql"
	"github.com/tmc/langchaingo/tools/sqldatabase/sqlite3"
)

var (
	mysqlDSN    = flag.String("mysql", "", "mysql dsn (e.g. user:pwd@tcp(localhost:3306)/test)")
	sqlite3DSN  = flag.String("sqlite3", "", "sqlite3 dsn (e.g. test.db)")
	pgxDSN      = flag.String("postgre", "", "postgre dsn (e.g. postgres://db_user:mysecretpassword@localhost:5438/test?sslmode=disable)")
	csv         = flag.String("csv", "", "csv dir or file")
	useAllTable = flag.Bool("all", true, "use all table or choose by question")
)

func main() {
	flag.Parse()

	godotenv.Overload()

	llm, err := openai.NewChat()
	if err != nil {
		log.Fatalf("openai err: %s", err)
	}

	var chain *datachain.DataChain
	if *sqlite3DSN != "" {
		chain, err = datachain.New(llm, sqlite3.EngineName, *sqlite3DSN, *useAllTable)
	} else if *mysqlDSN != "" {
		chain, err = datachain.New(llm, mysql.EngineName, *mysqlDSN, *useAllTable)
	} else if *pgxDSN != "" {
		chain, err = datachain.New(llm, postgresql.EngineName, *pgxDSN, *useAllTable)
	} else if *csv != "" {
		dbFile := path.Join(os.TempDir(), "chat2data.db")
		os.Remove(dbFile)
		defer os.Remove(dbFile)
		err = cmd.Load(dbFile, *csv)
		if err != nil {
			log.Fatalf("load csv err: %s", err)
		}
		chain, err = datachain.New(llm, sqlite3.EngineName, dbFile, *useAllTable)
	} else {
		log.Fatalf("no dsn")
	}
	if err != nil {
		log.Fatalf("open database err: %s", err)
	}

	defer chain.Close()

	if err := cmd.CLI(chain); err != nil {
		fmt.Println(err)
	}
}
