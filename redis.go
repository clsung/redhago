package redhago

import (
	"encoding/json"
	"time"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/golang/glog"
)

type Pool struct {
	pool *redigo.Pool
}

func NewPool(addr, pass string, maxIdle, idleTimeout int) *Pool {
	pool := &redigo.Pool{
		MaxIdle: maxIdle,
		// if idleTimeout is zero, idle connections are not closed.
		IdleTimeout: time.Duration(idleTimeout) * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			if pass != "" {
				if _, err := c.Do("AUTH", pass); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			glog.Info("TestOnBorrow begin")
			_, err := c.Do("EXISTS", "foobar")
			glog.Info("TestOnBorrow end")
			return err
		},
	}
	glog.Infof("create redis pool %v with maxIdleConn: %d", pool, maxIdle)
	return &Pool{pool: pool}
}

func (p *Pool) SetDb(db string) error {
	conn := p.pool.Get()
	if err := conn.Send("select", db); err != nil {
		glog.Errorf("select %v failed: %v", db, err)
		return err
	}
	glog.Infof("select %v success", db)
	return nil
}

func (p *Pool) Get(key string) ([]byte, error) {
	conn := p.pool.Get()

	r, err := redigo.Bytes(conn.Do("GET", key))
	if err != nil {
		glog.Errorf("get %s failed: %v", key, err)
	}
	return r, err
}

func (p *Pool) Setex(key string, obj interface{}, exp int) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	conn := p.pool.Get()

	_, err = conn.Do("SETEX", key, exp, data)
	if err != nil {
		glog.Errorf("setex %s failed: %v", key, err)
	}
	return err
}

func (p *Pool) Set(key string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	conn := p.pool.Get()

	_, err = conn.Do("SET", key, data)
	if err != nil {
		glog.Errorf("set %s failed: %v", key, err)
	}
	return err
}

func (p *Pool) Lpush(key string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	conn := p.pool.Get()

	_, err = conn.Do("LPUSH", key, data)
	if err != nil {
		glog.Errorf("lpush %s failed: %v", key, err)
	}
	return err
}
