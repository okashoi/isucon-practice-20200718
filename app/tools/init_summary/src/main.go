package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dbConnPoolSize = 10
)

type Config struct {
	Database struct {
		Dbname   string `json:"dbname"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
	} `json:"database"`
}

type Memo struct {
	Id        int
	User      int
	Content   string
	IsPrivate int
	CreatedAt string
	UpdatedAt string
	Username  string
	Summary   string
}

type Memos []*Memo

func getFirstLine(memo *Memo) string {
	return strings.Split(memo.Content, "\n")[0]
}

var (
	dbConnPool chan *sql.DB
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	env := os.Getenv("ISUCON_ENV")
	if env == "" {
		env = "local"
	}
	config := loadConfig("../../../config/" + env + ".json")
	db := config.Database
	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8",
		db.Username, db.Password, db.Host, db.Port, db.Dbname,
	)
	log.Printf("db: %s", connectionString)

	dbConnPool = make(chan *sql.DB, dbConnPoolSize)
	for i := 0; i < dbConnPoolSize; i++ {
		conn, err := sql.Open("mysql", connectionString)
		if err != nil {
			log.Panicf("Error opening database: %v", err)
		}
		dbConnPool <- conn
		defer conn.Close()
	}

	migration()

}

func loadConfig(filename string) *Config {
	log.Printf("loading config file: %s", filename)
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	var config Config
	err = json.Unmarshal(f, &config)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	return &config
}

func migration() {
	dbConn := <-dbConnPool
	defer func() {
		dbConnPool <- dbConn
	}()

	log.Println("start load initial date")
	rows, err := dbConn.Query("SELECT id, content FROM memos")
	if err != nil {
		log.Panic(err)
		return
	}
	memos := make(Memos, 0)
	for rows.Next() {
		memo := Memo{}
		rows.Scan(&memo.Id, &memo.Content)
		memo.Summary = getFirstLine(&memo)
		memos = append(memos, &memo)
	}
	rows.Close()
	log.Println("finish load initial date")

	db, _ := dbConn.Begin()
	stmt, err := db.Prepare("UPDATE memos SET summary=? WHERE id=?")
	defer stmt.Close()
	for _, memo := range memos {
		if _, err := stmt.Exec(memo.Summary, memo.Id); err != nil {
			log.Panic(err)
			return
		}
	}
	db.Commit()
}
