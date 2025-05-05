package enter

import (
	"time"
)

type TicketInfo struct {
	AccountServerId int64
	YostarUID       int64
	Ticket          string
	EndTime         time.Time
}

// AddEnterTicket 添加EnterTicket
func AddEnterTicket(accountServerId, yostarUID int64, ticker string) bool {
	e := getEnterSet()

	info := &TicketInfo{
		AccountServerId: accountServerId,
		YostarUID:       yostarUID,
		Ticket:          ticker,
		EndTime:         time.Now().Add(30 * time.Minute), // 30分钟有效期
	}
	if e.EnterTicketMap == nil {
		e.EnterTicketMap = make(map[string]*TicketInfo)
	}
	e.EnterTicketMap[ticker] = info

	return true
}

// DelEnterTicket 删除Ticket,如果此Ticket存在返回true,不存在返回false
func DelEnterTicket(ticker string) bool {
	e := getEnterSet()
	if e.EnterTicketMap == nil {
		e.EnterTicketMap = make(map[string]*TicketInfo)
	}
	if _, ok := e.EnterTicketMap[ticker]; ok {
		delete(e.EnterTicketMap, ticker)
		return true
	}
	return false
}

// GetEnterTicketInfo 通过ticker获取登录信息
func GetEnterTicketInfo(ticker string) *TicketInfo {
	e := getEnterSet()
	if info, ok := e.EnterTicketMap[ticker]; ok {
		if time.Now().After(info.EndTime) {
			DelEnterTicket(ticker)
		}
		return info
	}
	return nil
}

// GetAllEnterTicketInfo 获取全部登录信息
func GetAllEnterTicketInfo() map[string]*TicketInfo {
	allTicketInfo := make(map[string]*TicketInfo)
	e := getEnterSet()
	for k, v := range e.EnterTicketMap {
		allTicketInfo[k] = v
	}
	return allTicketInfo
}
