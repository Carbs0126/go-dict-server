package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"strconv"
)

const ServerPort = 4000
const DatabaseNameAK = "dict_a_k.db"
const DatabaseNameLZ = "dict_l_z.db"

type application struct {
	infoLog      *log.Logger
	errorLog     *log.Logger
	dbak         *sql.DB
	dblz         *sql.DB
	searchStmtAK *sql.Stmt
	searchStmtLZ *sql.Stmt
}

func main() {
	// 1. 解析外部参数
	var port string
	flag.StringVar(&port, "port", strconv.Itoa(ServerPort), "HTTP server port")
	flag.Parse()
	// 2. 初始化log
	// 第一个参数，可以是file，如果是file，则会把info输出到file中，而不用在启动参数中明确指定重定向
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	// 3. 初始化数据库
	dbak, searchStmtAK := initSQLite3DB(DatabaseNameAK)
	dblz, searchStmtLZ := initSQLite3DB(DatabaseNameLZ)
	defer dbak.Close()
	defer dblz.Close()
	// 4. 初始化application
	app := new(application)
	app.infoLog = infoLog
	app.errorLog = errorLog
	app.dbak = dbak
	app.searchStmtAK = searchStmtAK
	app.dblz = dblz
	app.searchStmtLZ = searchStmtLZ
	// 5. 初始化服务器
	infoLog.Println("Starting Server")
	srv := &http.Server{
		Addr:     ":" + port,
		ErrorLog: errorLog,
		Handler:  app.routes()}
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}

func initSQLite3DB(dbpath string) (*sql.DB, *sql.Stmt) {
	db, err := sql.Open("sqlite3", dbpath)
	checkError(err, "initSQLite3DB 1 :")
	err = db.Ping()
	checkError(err, "initSQLite3DB 2 :")
	strQuery := fmt.Sprintf("SELECT translation FROM dict where word=? LIMIT 1")
	searchStmt, err := db.Prepare(strQuery)
	checkError(err, "initSQLite3DB 3 :")
	return db, searchStmt
}

func checkError(err error, msg string) {
	if err != nil {
		fmt.Println(msg, err)
		panic(err)
	}
}
