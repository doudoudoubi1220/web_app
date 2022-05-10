package redis

import (
	"strconv"
	"time"
	"web_app/models"

	"github.com/go-redis/redis"
)

func GetPostIDsInOrder(p *models.ParamPostList) ([]string, error) {
	//从redis获取id

	//1.根据用户请求中携带的order参数确定要查询的redis key
	key := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		key = getRedisKey(KeyPostScoreZSet)
	}
	//2.确定查询的索引起始点
	return getIDsFormKey(key, p.Page, p.Size)
	//start := (p.Page - 1) * p.Size
	//end := start + p.Size - 1
	////3.ZREVRANGE按分数从大到小的顺序查询指定数量的元素
	//return client.ZRevRange(key, start, end).Result()
}

//GetPostVoteData根据ids查询每篇帖子的投票赞成票的数据
func GetPostVoteData(ids []string) (data []int64, err error) {

	//for _, id := range ids {
	//	key := getRedisKey(KeyPostVotedZSetPrefix + id)
	//	//查找key中分数是1的元素的数量->统计每篇帖子的赞成票的数量
	//	v := client.ZCount(key, "1", "1").Val()
	//	data = append(data, v)
	//}
	//keys := make([]string, 0, len(ids))

	//使用pipeline一次发送多条命令减少RTT
	pipeline := client.Pipeline()
	for _, id := range ids {
		key := getRedisKey(KeyPostVotedZSetPrefix + id)
		pipeline.ZCount(key, "1", "1")
	}
	cmders, err := pipeline.Exec()
	if err != nil {
		return nil, err
	}
	for _, cmder := range cmders {
		v := cmder.(*redis.IntCmd).Val()
		data = append(data, v)

	}

	return
}

func getIDsFormKey(key string, page, size int64) ([]string, error) {

	start := (page - 1) * size
	end := start + size - 1
	//3.ZREVRANGE按分数从大到小的顺序查询指定数量的元素
	return client.ZRevRange(key, start, end).Result()
}

func GetCommunityPostIDsInOrder(p *models.ParamPostList) ([]string, error) {

	orderKey := getRedisKey(KeyPostTimeZSet)
	if p.Order == models.OrderScore {
		orderKey = getRedisKey(KeyPostScoreZSet)
	}
	//使用zinterstore把分区的帖子set与分数的zset生成一个新的zset
	//针对新的zset按之前的逻辑取数据 orderKey string, communityID, page, size int64

	//社区的key
	ckey := getRedisKey(KeyCommunitySetPrefix + strconv.Itoa(int(p.CommunityID)))

	//利用缓存key减少zinterstore的执行次数
	key := orderKey + strconv.Itoa(int(p.CommunityID))
	if client.Exists(orderKey).Val() < 1 {
		//不存在，需要计算
		pipeline := client.Pipeline()
		pipeline.ZInterStore(key, redis.ZStore{
			Aggregate: "MAX",
		}, ckey, orderKey)
		pipeline.Expire(key, 60*time.Second)
		_, err := pipeline.Exec()
		if err != nil {
			return nil, err
		}
	}
	return getIDsFormKey(key, p.Page, p.Size)
	//从redis获取id
	//1.根据用户请求中携带的order参数确定要查询的redis key
	//key := getRedisKey(KeyPostTimeZSet)
	//if p.Order == models.OrderScore {
	//	key = getRedisKey(KeyPostScoreZSet)
	//}
	////2.确定查询的索引起始点
	//start := (p.Page - 1) * p.Size
	//end := start + p.Size - 1
	////3.ZREVRANGE按分数从大到小的顺序查询指定数量的元素
	//return client.ZRevRange(key, start, end).Result()

}
