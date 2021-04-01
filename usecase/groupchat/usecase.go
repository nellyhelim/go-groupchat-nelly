package groupchat

import (
	"github.com/lolmourne/go-groupchat/model"
)

func (u *Usecase) CreateGroupchat(roomName, adminID, description, categoryID string) (*model.Room, error) {
	room, err := u.dbRsc.CreateRoom(roomName, adminID, description, categoryID)

	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (u *Usecase) EditGroupchat(name, description, categoryID string) (*model.Room, error) {
	room, err := u.dbRsc.EditGroupchat(name, description, categoryID)

	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (u *Usecase) JoinRoom(roomID, userID string) error {
	err := u.dbRsc.AddRoomParticipant(roomID, userID)

	if err != nil {
		return err
	}
	return nil
}
