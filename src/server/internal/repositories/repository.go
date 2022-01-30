package repositories

import "github.com/tarantool/go-tarantool"

type Repository interface {
	KeepAlive(dbc *tarantool.Connection)
}
