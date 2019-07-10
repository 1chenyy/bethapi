package db

import (
	"bethapi/models"
	"sync"
)

type PriceDBUtil struct {
	PAM models.PriceAndMakcap
	History models.HistoryPrice
	rw sync.RWMutex
}

func NewPriceDBUtil()*PriceDBUtil{
	return &PriceDBUtil{
		rw:sync.RWMutex{},
		PAM:models.PriceAndMakcap{},
		History:models.HistoryPrice{},
	}
}

func (u *PriceDBUtil)GetPriceAndMakcapFromDb()models.PriceAndMakcap{
	u.rw.RLock()
	defer u.rw.RUnlock()
	return u.PAM
}

func (u *PriceDBUtil)SetPriceAndMakcap(pam models.PriceAndMakcap){
	u.rw.Lock()
	defer u.rw.Unlock()
	if pam.Mktcap=="" || pam.Price=="" {
		return
	}
	u.PAM = pam
}

func (u *PriceDBUtil)GetHistoryPriceFromDb()models.HistoryPrice{
	u.rw.RLock()
	defer u.rw.RUnlock()
	return u.History
}

func (u *PriceDBUtil)SetHistoryPrice(history models.HistoryPrice){
	u.rw.Lock()
	defer u.rw.Unlock()
	u.History = history
}
