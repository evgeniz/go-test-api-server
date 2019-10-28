package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var DB *sql.DB

type UserModel struct {
	Id       uint64  `json:"id"`
	Balance  float64 `json:"balance"`
	DepCount uint64  `json:"depositCount"`
	DepSum   float64 `json:"depositSum"`
	BetCount uint64  `json:"betCount"`
	BetSum   float64 `json:"betSum"`
	WinCount uint64  `json:"winCount"`
	WinSum   float64 `json:"winSum"`
}

var UserIDs = make(map[uint64]bool)

func DBCreateUser(userId uint64) {
	stmt, err := DB.Prepare("INSERT users SET id=?, balance=?, deposit_count=?, deposit_sum=?, bet_count=?, bet_sum=?, win_count=?, win_sum=?")
	if err != nil {
		log.Fatal("DB.Prepare: ", err)
	}
	_, err = stmt.Exec(userId, 0.0, 0, 0.0, 0, 0.0, 0, 0.0)
	if err != nil {
		log.Fatal("stmt.Exec: ", err)
	}
}

func DBUpdateUser(u *UserModel) {
	stmt, err := DB.Prepare("UPDATE users SET balance=?, deposit_count=?, deposit_sum=?, bet_count=?, bet_sum=?, win_count=?, win_sum=? WHERE id=?")
	if err != nil {
		log.Fatal("DB.Prepare: ", err)
	}

	_, err = stmt.Exec(u.Balance, u.DepCount, u.DepSum, u.BetCount, u.BetSum, u.WinCount, u.WinSum, u.Id)
	if err != nil {
		log.Fatal("stmt.Exec: ", err)
	}
}

func NewDB(dbURL string) {
	dataBase, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal("NewDB: ", err)
	}

	DB = dataBase

	if err := DB.Ping(); err != nil {
		log.Fatal("NewDB: ", err)
	}
}