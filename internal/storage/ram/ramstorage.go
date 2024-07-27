package ramstorage

import (
	"fmt"
	"sync"
	"time"

	"github.com/anosovs/datalink/internal/models"
	"github.com/anosovs/datalink/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	mx sync.Mutex
	messages []models.Message
}

func New() (*Storage) {
	return &Storage{
		messages: []models.Message{},
	}
}


func (s *Storage) SaveMsg(msgToSave string, count int) (uid string, err error) {
	uid = uuid.New().String()

	msg := models.Message{
		Uuid: uid,
		Message: msgToSave,
		Count: int8(count),
		CreatedAt: time.Now(),		
	}
	fmt.Println(msg)
	
	s.messages = append(s.messages, msg)
	return uid, nil
}

func (s *Storage) GetMsg(uid string) (msg string, cnt int, err error) {
	_, m, err := s.getMessageModelByUid(uid)
	if err != nil {
		return "", 0, err
	}

	return m.Message, int(m.Count), nil
}

func (s *Storage) DeleteMsg(uid string) error {
	index, _, err := s.getMessageModelByUid(uid)
	if err != nil {
		return err
	}

	s.messages = append(s.messages[:index], s.messages[index+1:]...)

	return nil
}


func (s *Storage) DecreaseCountMsg(uid string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	i, m, err := s.getMessageModelByUid(uid)
	if err != nil {
		return err
	}
	s.messages[i].Count = m.Count - 1

	return nil
}

func (s *Storage) DeleteOldMessages(days int) error{
	availableUntil := time.Now().Add(-time.Hour * 24 * time.Duration(days))
	for _, m := range s.messages {
		if (m.CreatedAt.Unix() < availableUntil.Unix()) {
			s.DeleteMsg(m.Uuid)
		}
	}
	return nil
}


func (s *Storage) getMessageModelByUid(uid string) (index int, msg models.Message, err error) {
	for i, m := range s.messages {
		if (m.Uuid == uid) {
			msg = m
			index = i
		}
	}
	if msg.Uuid == "" {
		return 0, models.Message{}, storage.ErrUuidNotFount
	}

	return index, msg, nil
}