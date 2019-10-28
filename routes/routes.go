package routes

import (
	"encoding/json"
	"github.com/evgeniz/test-task-guru-team/cache"
	"github.com/evgeniz/test-task-guru-team/db"
	"log"
	"net/http"
	"time"
)

func SetRoutes()  {
	http.HandleFunc("/user/create", createUser())
	http.HandleFunc("/user/get", getUser())
	http.HandleFunc("/user/deposit", depositUser())
	http.HandleFunc("/transaction", transaction())
}

func createUser() http.HandlerFunc {
	type request struct {
		Id      uint64  `json:"id"`
		Balance float64 `json:"balance"`
		Token   string  `json:"token"`
	}

	type response struct {
		Error string `json:"error"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var req request
		err := decoder.Decode(&req)
		if err != nil {
			log.Fatal(err)
		}

		e := response{}
		if _, ok := cache.UCache.Load(req.Id); ok {
			e.Error = "Id already exists"
		} else {
			cache.UCache.Store(req.Id, &cache.User{req.Balance, req.Token})
			cache.DCache.Store(req.Id, make([]*cache.Deposit, 0))
			cache.TCache.Store(req.Id, make([]*cache.Transaction, 0))
		}
		resp, err := json.Marshal(e)
		_, err = w.Write(resp)
		if err != nil {
			panic(err)
		}
		if e.Error == "" {
			go db.DBCreateUser(req.Id)
		}
	}
}
func getUser() http.HandlerFunc {
	type request struct {
		Id    uint64 `json:"id"`
		Token string `json:"token"`
	}

	type unsucresp struct {
		Error string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var req request
		err := decoder.Decode(&req)
		if err != nil {
			log.Fatal(err)
		}

		u, _ := cache.UCache.Load(req.Id)
		e := unsucresp{}
		if req.Token != u.Token {
			e.Error = "invalid token"
			message, err := json.Marshal(e)
			_, err = w.Write(message)
			if err != nil {
				panic(err)
			}
			return
		}

		resp := db.UserModel{
			Id:      req.Id,
			Balance: u.Balance,
		}

		trs, _ := cache.TCache.Load(req.Id)
		for _, value := range trs {
			switch value.Type {
			case "Bet":
				{
					resp.BetCount++
					resp.BetSum += value.Sum
				}
			case "Win":
				{
					resp.WinCount++
					resp.WinSum += value.Sum
				}
			}
		}

		deps, _ := cache.DCache.Load(req.Id)
		for _, dep := range deps {
			resp.DepSum += dep.ABalance - dep.BBalance
			resp.DepCount++
		}

		rJson, err := json.Marshal(resp)
		_, err = w.Write(rJson)
		if err != nil {
			panic(err)
		}
	}
}

func depositUser() http.HandlerFunc {
	type request struct {
		UserId    uint64  `json:"userId"`
		DepositId uint64  `json:"depositId"`
		Amount    float64 `json:"amount"`
		Token     string  `json:"token"`
	}

	type response struct {
		Error   string  `json:"error"`
		Balance float64 `json:"balance"`
	}

	type unsucresp struct {
		Error string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var req request
		err := decoder.Decode(&req)
		if err != nil {
			log.Fatal(err)
		}

		u, _ := cache.UCache.Load(req.UserId)
		e := unsucresp{}
		if req.Token != u.Token {
			e.Error = "invalid token"
			message, err := json.Marshal(e)
			_, err = w.Write(message)
			if err != nil {
				panic(err)
			}
			return
		}

		deposit := &cache.Deposit{
			Id:       req.DepositId,
			BBalance: u.Balance,
			ABalance: u.Balance + req.Amount,
			Time:     time.Now(),
		}
		depSlice, _ := cache.DCache.Load(req.UserId)
		depSlice = append(depSlice, deposit)
		cache.DCache.Store(req.UserId, depSlice)

		u.Balance += req.Amount

		resp := response{Balance: u.Balance}
		rJson, err := json.Marshal(resp)
		_, err = w.Write(rJson)
		if err != nil {
			panic(err)
		}
		db.UserIDs[req.UserId] = true
	}
}

func transaction() http.HandlerFunc {
	type request struct {
		UserId        uint64  `json:"userId"`
		TransactionId uint64  `json:"transactionId"`
		Type          string  `json:"type"`
		Amount        float64 `json:"amount"`
		Token         string  `json:"token"`
	}

	type response struct {
		Error   string  `json:"error"`
		Balance float64 `json:"balance"`
	}

	type unsucresp struct {
		Error string
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(405), 405)
			return
		}

		decoder := json.NewDecoder(r.Body)
		var req request
		err := decoder.Decode(&req)
		if err != nil {
			log.Fatal(err)
		}

		u, _ := cache.UCache.Load(req.UserId)
		e := unsucresp{}
		if req.Token != u.Token {
			e.Error = "invalid token"
			message, err := json.Marshal(e)
			_, err = w.Write(message)
			if err != nil {
				panic(err)
			}
			return
		}

		if req.Type == "Bet" && u.Balance-req.Amount < 0 {
			e.Error = "invalid transaction: insufficient funds"
			message, err := json.Marshal(e)
			_, err = w.Write(message)
			if err != nil {
				panic(err)
			}
			return
		}

		var BAfter float64
		if req.Type == "Bet" {
			BAfter = u.Balance - req.Amount
		} else {
			BAfter = u.Balance + req.Amount
		}

		t := &cache.Transaction{
			Id:       req.TransactionId,
			Type:     req.Type,
			Sum:      req.Amount,
			BBalance: u.Balance,
			ABalance: BAfter,
			Time:     time.Now(),
		}
		u.Balance = t.ABalance

		tranSlice, _ := cache.TCache.Load(req.UserId)
		tranSlice = append(tranSlice, t)
		cache.TCache.Store(req.UserId, tranSlice)

		resp := response{Balance: u.Balance}
		rJson, err := json.Marshal(resp)
		_, err = w.Write(rJson)
		if err != nil {
			panic(err)
		}
		db.UserIDs[req.UserId] = true
	}
}

func NewUserModel(userId uint64) *db.UserModel {
	var depC uint64
	var depS float64
	deps, _ := cache.DCache.Load(userId)
	for _, d := range deps {
		depS += d.ABalance - d.BBalance
		depC++
	}
	var betC uint64
	var winC uint64
	var betS float64
	var winS float64
	trs, _ := cache.TCache.Load(userId)
	for _, value := range trs {
		switch value.Type {
		case "Bet":
			{
				betC++
				betS += value.Sum
			}
		case "Win":
			{
				winC++
				winS += value.Sum
			}
		}
	}
	u, _ := cache.UCache.Load(userId)
	return &db.UserModel{
		Id:       userId,
		Balance:  u.Balance,
		DepCount: depC,
		DepSum:   depS,
		BetCount: betC,
		BetSum:   betS,
		WinCount: winC,
		WinSum:   winS,
	}
}
