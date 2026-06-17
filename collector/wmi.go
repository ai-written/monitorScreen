package collector

import (
	"github.com/yusufpapurcu/wmi"
)

func wmiQuery(query string, dst interface{}) error {
	return wmi.Query(query, dst)
}
