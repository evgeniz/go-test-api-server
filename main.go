package main

import (
	"log"
	"net/http"
	"github.com/evgeniz/test-task-guru-team/db"
	"github.com/evgeniz/test-task-guru-team/routes"
	"time"
)

func main() {
	db.NewDB("root:root@tcp(localhost:3306)/guru_team")
	defer db.DB.Close()

	routes.SetRoutes()

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if len(db.UserIDs) > 0 {
				for k, v := range db.UserIDs {
					if v == true {
						db.DBUpdateUser(routes.NewUserModel(k))
						delete(db.UserIDs, k)
					}
				}
			}
		}
	}()

	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
