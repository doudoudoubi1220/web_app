package logic

import (
	"web_app/dao/mysql"
	"web_app/dao/redis"
	"web_app/models"
	"web_app/pkg/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	//1.生成post id
	p.ID = snowflake.GenID()
	//2.保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	//3.返回
	return
}

func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {

	//查询并组合我们接口想用的数据
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("mysql.GetPostByID(pid) failed", zap.Error(err))
		return
	}

	//根据作者id查询作者信息
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("mysql.GetUserByID(pid) failed",
			zap.Int64("author_id", post.AuthorID),
			zap.Error(err))
		return
	}

	community, err := mysql.GetCommunityDetailByID(post.CommunityID)
	if err != nil {
		zap.L().Error("mysql.GetCommunityDetailByID(pid) failed",
			zap.Int64("community_id", post.CommunityID),
			zap.Error(err))
		return
	}
	//拼接数据
	data = &models.ApiPostDetail{
		AuthorName:      user.Username,
		Post:            post,
		CommunityDetail: community,
	}
	return

}

func GetPostList(page, size int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(page, size)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))

	for _, post := range posts {

		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(pid) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)

	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//去redis查询id列表
	ids, err := redis.GetPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("GetPostIDsInOrder return 0 data")
		return
	}

	zap.L().Debug("GetPostList2", zap.Any("ids", ids))

	//3.根据id去MySQL数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	//提前查询好每篇帖子的投票数据
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	//将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {

		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(pid) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)

	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	//去redis查询id列表
	ids, err := redis.GetCommunityPostIDsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		zap.L().Warn("GetPostIDsInOrder return 0 data")
		return
	}

	zap.L().Debug("GetPostList2", zap.Any("ids", ids))

	//3.根据id去MySQL数据库查询帖子详细信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}

	//提前查询好每篇帖子的投票数据
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	//将帖子的作者及分区信息查询出来填充到帖子中
	for idx, post := range posts {

		//根据作者id查询作者信息
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("mysql.GetUserByID(post.AuthorID) failed",
				zap.Int64("author_id", post.AuthorID),
				zap.Error(err))
			continue
		}

		//根据社区id查询社区详细信息
		community, err := mysql.GetCommunityDetailByID(post.CommunityID)
		if err != nil {
			zap.L().Error("mysql.GetCommunityDetailByID(pid) failed",
				zap.Int64("community_id", post.CommunityID),
				zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.Username,
			VoteNum:         voteData[idx],
			Post:            post,
			CommunityDetail: community,
		}
		data = append(data, postdetail)

	}
	return
}

// GetPostListNew 将两个帖子列表查询逻辑合二为一的函数
func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 根据请求参数的不同，执行不同的逻辑
	if p.CommunityID == 0 {
		// 查询所有帖子列表
		data, err = GetPostList2(p)
	} else {
		// 根据社区id查询帖子列表
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		zap.L().Error("GetPostListNew failed", zap.Error(err))
		return nil, err
	}
	return
}
