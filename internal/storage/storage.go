package storage

import "errors"


var (
	ErrUuidNotFount = errors.New("uuid not found")
)

type Storage interface {
	SaveMsg(msgToSave string, count int) (uid string, err error)
	GetMsg(uid string) (msg string, cnt int, err error)
	DeleteMsg(uid string) error
	DecreaseCountMsg(uid string) error
	DeleteOldMessages(days int) error
}