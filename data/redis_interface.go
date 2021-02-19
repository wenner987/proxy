package data

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"proxy/configure"
	"strconv"
	"strings"
	"sync"
)

var redisCli *redis.Client
var redisOnce = &sync.Once{}

func initRedis() {
	builder := strings.Builder{}
	builder.WriteString(configure.GetConfigInstance().Redis.Ip)
	builder.WriteString(":")
	builder.WriteString(strconv.Itoa(configure.GetConfigInstance().Redis.Port))
	redisAddr := builder.String()
	log.Printf("Redis ip: %s, port: %d", configure.GetConfigInstance().Redis.Ip,
		configure.GetConfigInstance().Redis.Port)
	log.Println("redis init addr:", redisAddr)
	passWord := configure.GetConfigInstance().Redis.Passwd
	redisCli = redis.NewClient(
		&redis.Options{
			Addr: redisAddr,
			Password: passWord})
}

func UpdateZSetInfo(score, set int, ip string) {
	redisOnce.Do(initRedis)
	oldScore := getOldScore(set, ip)
	if oldScore < -20 && score < 0 || oldScore > 100 && score > 0 {
		return
	}
	newScore, err := redisCli.ZIncrBy(context.Background(), strconv.Itoa(set), float64(score), ip).Result()
	if err != nil {
		log.Println("Update score fail:", set, " ", ip)
		log.Println("err:", err)
	} else {
		log.Printf("Update score succ: %d, %s new score: %f", set, ip, newScore)
	}
}

func GetIpFromZSet(set int) string {
	redisOnce.Do(initRedis)
	score, _ := redisCli.ZRevRange(context.Background(), strconv.Itoa(set), 0, 0).Result()
	return score[0]
}

func getOldScore(set int, ip string) float64 {
	oldScore, _ := redisCli.ZScore(context.Background(), strconv.Itoa(set), ip).Result()
	log.Printf("Old score: %f", oldScore)
	return oldScore
}