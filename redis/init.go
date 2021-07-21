package redis

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/gomodule/redigo/redis"
)

// DB represents a pool of connections to Redis.
type DB struct {
	*redis.Pool
}

// NewDB creates a pool of connections to Redis and pings Redis.
func NewDB(hostname, port, password string) (db *DB, err error) {
	pool := newPool(hostname, port, password)
	db = &DB{pool}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("cannot create Redis database connections: %w", err)
	}
	return db, nil
}

var (
	ErrConnection = errors.New("cannot connect to Redis")
	ErrAuth       = errors.New("authentication failed")
	ErrNotPong    = errors.New("message received is not expected PONG")
)

func newPool(host, port, password string) *redis.Pool {
	address := net.JoinHostPort(host, port)
	pool := new(redis.Pool)
	pool.MaxIdle = 7
	pool.MaxActive = 29 // max number of connections
	const idleTimeout = 3 * time.Second
	pool.IdleTimeout = idleTimeout
	pool.Wait = true
	pool.MaxConnLifetime = 0 // never close connection based on age
	pool.Dial = func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", address)
		if err != nil {
			return nil, fmt.Errorf("%w: at address %s: %s",
				ErrConnection, address, err)
		}
		if password != "" {
			reply, err := c.Do("AUTH", password)
			if err != nil {
				_ = c.Close()
				return nil, fmt.Errorf("%w: %s", ErrAuth, err)
			}
			err = CheckOKString(reply)
			if err != nil {
				_ = c.Close()
				return nil, fmt.Errorf("%w: %s", ErrAuth, err)
			}
		}
		return c, nil
	}
	return pool
}

// Ping pings the Redis database and returns an error if it fails.
func (db *DB) Ping() error {
	c := db.Get()
	defer c.Close()
	reply, err := c.Do("PING")
	if err != nil {
		return fmt.Errorf("%w: %s", ErrDoCommand, err)
	}

	s, err := redis.String(reply, nil)
	if err != nil {
		return err
	} else if s != "PONG" {
		return fmt.Errorf("%w: %s", ErrNotPong, s)
	}
	return nil
}
