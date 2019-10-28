package cache

import (
	"sync"
	"time"
)

type User struct {
	Balance float64
	Token   string
}

type Deposit struct {
	Id       uint64
	BBalance float64
	ABalance float64
	Time     time.Time
}

type Transaction struct {
	Id       uint64
	Type     string
	Sum      float64
	BBalance float64
	ABalance float64
	Time     time.Time
}

var (
	UCache = &Users{cache: make(map[uint64]*User)}
	DCache = &Deposits{cache: make(map[uint64][]*Deposit)}
	TCache = &Transactions{cache: make(map[uint64][]*Transaction)}
)

type Users struct {
	mx    sync.RWMutex
	cache map[uint64]*User
}

type Deposits struct {
	mx    sync.RWMutex
	cache map[uint64][]*Deposit
}

type Transactions struct {
	mx    sync.RWMutex
	cache map[uint64][]*Transaction
}

func (u *Users) Load(userId uint64) (*User, bool) {
	u.mx.RLock()
	defer u.mx.RUnlock()
	val, ok := u.cache[userId]
	return val, ok
}

func (u *Users) Store(userId uint64, user *User) {
	u.mx.Lock()
	defer u.mx.Unlock()
	u.cache[userId] = user
}

func (d *Deposits) Load(userId uint64) ([]*Deposit, bool) {
	d.mx.RLock()
	defer d.mx.RUnlock()
	val, ok := d.cache[userId]
	return val, ok
}

func (d *Deposits) Store(userId uint64, deps []*Deposit) {
	d.mx.Lock()
	defer d.mx.Unlock()
	d.cache[userId] = deps
}

func (t *Transactions) Load(userId uint64) ([]*Transaction, bool) {
	t.mx.RLock()
	defer t.mx.RUnlock()
	val, ok := t.cache[userId]
	return val, ok
}

func (t *Transactions) Store(userId uint64, deps []*Transaction) {
	t.mx.Lock()
	defer t.mx.Unlock()
	t.cache[userId] = deps
}
