// redis
package database

import (
	//	"fmt"

	"github.com/garyburd/redigo/redis"
)

const (
	redisUrl      = ""
	redisPassword = ""
)

type RedisConn struct {
	conn redis.Conn
}

//link
func (rc *RedisConn) redisDial() (err error) {

	rc.conn, err = redis.Dial("tcp", redisUrl)
	if err == nil {
		_, err = rc.conn.Do("AUTH", redisPassword)
	}
	return
}

//close
func (rc *RedisConn) redisClose() {
	rc.conn.Close()
}

//insert
func (rc *RedisConn) InsertData(key string, value interface{}, expiredTime int) (err error) {

	defer rc.redisClose()

	err = rc.redisDial()
	if err != nil {
		return
	}
	_, err = rc.conn.Do("SET", key, value, "EX", expiredTime)

	return err
}

func (rc *RedisConn) ExitsData(key string) (bool, error) {

	defer rc.redisClose()

	err := rc.redisDial()
	if err != nil {
		return false, err
	}

	return redis.Bool(rc.conn.Do("EXISTS", key))
}

func (rc *RedisConn) GetData(key string) ([]byte, error) {
	defer rc.redisClose()

	err := rc.redisDial()
	if err != nil {
		return nil, err
	}

	return redis.Bytes(rc.conn.Do("GET", key))
}

func (rc *RedisConn) DeleteData(key string) error {
	defer rc.redisClose()

	err := rc.redisDial()
	if err != nil {
		return err
	}

	_, err = rc.conn.Do("DEL", key)

	return err
}
