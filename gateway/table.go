package gateway

import "sync"

var tables table

// map: fd<> conn
type table struct {
	did2conn sync.Map
}

func InitTables() {
	tables = table{
		did2conn: sync.Map{},
	}
}
