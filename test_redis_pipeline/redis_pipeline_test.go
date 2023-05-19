package test_redis_pipeline

import (
	"github.com/go-redis/redis"
	"testing"
)

type WriteBuffer struct {
	buf []byte
	pos int
}
type UserItem struct {
	GateNo int
	UserId string
	InTime int64
	Sended int64
}
type UserStatBuffer struct {
	Users  []UserItem
	Count  int
	Buffer WriteBuffer
}

func TestRedisPipeline(t *testing.T) {

	rdb := redis.NewClient(&redis.Options{
		Addr:         "192.168.78.171:6379",
		Password:     "",
		DB:           0,
		MinIdleConns: 10,
		PoolSize:     20,
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		t.Log(err)
	}

	var user_buff UserStatBuffer
	user_buff.Count = 3
	user_buff.Users = append(user_buff.Users, UserItem{
		GateNo: 1,
		UserId: "111",
		InTime: 123,
		Sended: 2,
	})
	user_buff.Users = append(user_buff.Users, UserItem{
		GateNo: 2,
		UserId: "222",
		InTime: 456,
		Sended: 3,
	})
	user_buff.Users = append(user_buff.Users, UserItem{
		GateNo: 4,
		UserId: "333",
		InTime: 789,
		Sended: 5,
	})
	pipe := rdb.Pipeline()
	for n := 0; n < user_buff.Count; n++ {
		pipe.HIncrBy(user_buff.Users[n].UserId, "abc", 1)
	}
	_, err = pipe.Exec()
	if err != nil {
		t.Log(err)
	}
}
