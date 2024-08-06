package main

import (
	"database/sql"
	"log"

	"github.com/TheRanomial/GoEcomApi/cmd/api"
	"github.com/TheRanomial/GoEcomApi/configs"
	"github.com/TheRanomial/GoEcomApi/db"
	"github.com/go-sql-driver/mysql"
)

func main(){
	cfg := mysql.Config{
		User:                 configs.Envs.DBUser,
		Passwd:               configs.Envs.DBPassword,
		Addr:                 configs.Envs.DBAddress,
		DBName:               configs.Envs.DBName,
		Net:                  "tcp",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	db, err := db.NewMySQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server:=api.NewAPIServer(":8080",db)

	if err:=server.Run(); err!=nil{
		log.Fatal(err)
	}

}

func initStorage(db *sql.DB){
	err:=db.Ping();
	if err!=nil{
		log.Fatal(err)
	}
	log.Println("DB successfully connected")
}