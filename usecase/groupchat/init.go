package groupchat

import (
	"github.com/lolmourne/go-groupchat/model"
	"github.com/lolmourne/go-groupchat/resource"
)

type Usecase struct {
	dbRsc resource.DBItf
}
type UsecaseItf interface {
	CreateGroupchat(roomName, adminID, description, categoryID string) (*model.Room, error)
	EditGroupchat(name, description, categoryID string) (*model.Room, error)
	JoinRoom(roomID, userID string) error
}

func NewUseCase(dbRsc resource.DBItf) UsecaseItf {
	return &Usecase{
		dbRsc: dbRsc,
	}
}
