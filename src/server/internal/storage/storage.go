package storage

import "github.com/tarantool/go-tarantool"

type Storage interface {
	KeepAlive(dbc *tarantool.Connection)
}
