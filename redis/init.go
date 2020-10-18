package redis

import (
	"fmt"
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

func newPool(hostname, port, password string) *redis.Pool {
	pool := new(redis.Pool)
	pool.MaxIdle = 7
	pool.MaxActive = 29 // max number of connections
	const idleTimeout = 3 * time.Second
	pool.IdleTimeout = idleTimeout
	pool.Wait = true
	pool.MaxConnLifetime = 0 // never close connection based on age
	pool.Dial = func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", hostname+":"+port)
		if err != nil {
			return nil, fmt.Errorf("connecting to Redis: %w", err)
		}
		if password != "" {
			reply, err := c.Do("AUTH", password)
			if err != nil {
				c.Close()
				return nil, fmt.Errorf("connecting to Redis: %w", err)
			}
			err = CheckOKString(reply)
			if err != nil {
				c.Close()
				return nil, fmt.Errorf("connecting to Redis: %w", err)
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
		return fmt.Errorf("ping Redis: %w", err)
	}
	s, err := CheckString(reply)
	if err != nil {
		return fmt.Errorf("ping Redis: %w", err)
	}
	if s != "PONG" {
		return fmt.Errorf("ping Redis: %w", err)
	}
	return nil
}
