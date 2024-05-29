package dto

import "app/entity"

type CreateGroupDTO struct {
	ChatName string
	Owner    entity.PersonInfo
	Members  []entity.PersonInfo
	Avatar   string
}
