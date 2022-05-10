package logic

import (
	"strconv"
	"web_app/dao/redis"
	"web_app/models"

	"go.uber.org/zap"
)

//基于用户投票的相关算法：https://www.ruanyifeng.com/blog/algorithm/

//投票功能：
//1.用户投票的数据
//2.

//PostVote 为帖子投票的函数
func PostVote(userID int64, p *models.ParamVoteData) error {
	//1.判断投票限制
	zap.L().Debug("VoteForPost",
		zap.Int64("userID", userID),
		zap.String("postID", p.PostID),
		zap.Int8("direction", p.Direction))
	return redis.VoteForPost(strconv.Itoa(int(userID)), p.PostID, float64(p.Direction))
	//2.更新帖子的分数
	//3.记录用户为该帖子的投票数据
}
