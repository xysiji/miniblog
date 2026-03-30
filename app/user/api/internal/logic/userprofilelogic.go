package logic

import (
	"context"
	"encoding/json"
	"fmt"

	"miniblog/app/relation/rpc/relationclient"
	"miniblog/app/user/api/internal/svc"
	"miniblog/app/user/api/internal/types"
	"miniblog/app/user/rpc/userclient"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserProfileLogic {
	return &UserProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserProfileLogic) UserProfile(req *types.UserProfileReq) (resp *types.UserProfileResp, err error) {
	// 1. 获取当前登录用户的 ID
	var currentUserId int64
	if uidVal := l.ctx.Value("userId"); uidVal != nil {
		switch v := uidVal.(type) {
		case json.Number:
			currentUserId, _ = v.Int64()
		case float64:
			currentUserId = int64(v)
		case string:
			fmt.Sscanf(v, "%d", &currentUserId)
		}
	}

	targetUserId := req.TargetUserId

	// 2. 呼叫 User RPC：获取目标用户的基本信息
	userInfoResp, err := l.svcCtx.UserRpc.UserInfo(l.ctx, &userclient.UserInfoRequest{
		UserId: targetUserId,
	})
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败: %v", err)
	}

	// 3. 呼叫 Relation RPC：获取目标用户的关注数和粉丝数
	countResp, err := l.svcCtx.RelationRpc.GetRelationCount(l.ctx, &relationclient.RelationCountReq{
		TargetUserId: targetUserId,
	})
	// 容错处理：即使计数服务暂时不可用，也不应该让整个主页崩溃，可以默认为0
	var followingCount, followerCount int64
	if err == nil {
		followingCount = countResp.FollowingCount
		followerCount = countResp.FollowerCount
	} else {
		l.Logger.Errorf("获取关系计数失败: %v", err)
	}

	// 4. 呼叫 Relation RPC：检查当前登录用户是否关注了该目标用户
	isFollow := false
	if currentUserId > 0 && currentUserId != targetUserId {
		followCheckResp, err := l.svcCtx.RelationRpc.CheckIsFollow(l.ctx, &relationclient.CheckIsFollowReq{
			CurrentUserId: currentUserId,
			TargetUserId:  targetUserId,
		})
		if err == nil {
			isFollow = followCheckResp.IsFollow
		}
	}

	// 5. 数据拼装与返回
	return &types.UserProfileResp{
		UserId:         userInfoResp.UserId,
		Username:       userInfoResp.Username,
		Avatar:         userInfoResp.Avatar,
		Bio:            userInfoResp.Bio,
		FollowingCount: followingCount,
		FollowerCount:  followerCount,
		IsFollow:       isFollow,
	}, nil
}
