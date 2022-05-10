package redis

//redis key

const (
	Prefix                 = "bluebell:"
	KeyPostTimeZSet        = "post:time"   //zset;帖子以及发帖时间
	KeyPostScoreZSet       = "post:score"  //zset;帖子及投票分数
	KeyPostVotedZSetPrefix = "post:voted:" //zset;记录用户及投票类型;参数是post id

	KeyCommunitySetPrefix = "community:" //保存每个分区下的帖子
)

func getRedisKey(key string) string {
	return Prefix + key
}
