package main

import (
	"github.com/evgeniz/test-task-guru-team/db"
	"github.com/evgeniz/test-task-guru-team/routes"
	"log"
	"net/http"
	"time"
)

func main() {
	db.NewDB("root:root@tcp(localhost:3306)/guru_team")
	defer db.DB.Close()

	routes.SetRoutes()

	go func() {
		for {
			time.Sleep(10 * time.Second)
			db.UIds.Mx.Lock()
			lenIds := len(db.UIds.Cache)
			db.UIds.Mx.Unlock()
			if lenIds > 0 {
				for k, v := range db.UIds.Cache {
					if v == true {
						db.DBUpdateUser(routes.NewUserModel(k))
						db.UIds.Mx.Lock()
						delete(db.UIds.Cache, k)
						db.UIds.Mx.Unlock()
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
