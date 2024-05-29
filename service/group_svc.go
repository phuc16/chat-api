package service

import (
	"app/dto"
	"app/entity"
	"app/pkg/trace"
	"app/pkg/utils"
	"context"
)

type GroupService struct {
	GroupRepo IGroupRepo
}

func NewGroupService(groupRepo IGroupRepo) *GroupService {
	return &GroupService{GroupRepo: groupRepo}
}

func (s *GroupService) CreateGroup(ctx context.Context, req *dto.CreateGroupReq) (res any, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	err = s.GroupRepo.SaveGroup(ctx, &req.Group)
	return
}

func (s *GroupService) GetGroupInfo(ctx context.Context, req *dto.GetGroupInfoReq) (res *entity.Group, err error) {
	ctx, span := trace.Tracer().Start(ctx, utils.GetCurrentFuncName())
	defer span.End()

	return s.GroupRepo.FindGroupByID(ctx, req.GroupID)
}
