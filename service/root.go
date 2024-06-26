package service

import (
	"chat_server_golang/repository"
	"chat_server_golang/types/schema"
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Service struct {
	repository  *repository.Repository

}

func NewService(repository  *repository.Repository) *Service{
	s := &Service{repository: repository}

	return s
}

func (s *Service) PublishServerStatusEvent(ip string, status bool ) {
		type ServerInfoEvent struct {
			IP string
			Status bool
		}

		e := &ServerInfoEvent{IP: ip, Status: status}
		ch := make(chan kafka.Event)

		if v, err := json.Marshal(e); err != nil {
			log.Println("Failed To Marshal", "err", err)
		} else if result, err := s.PublishEvent("test-topic",v, ch); err != nil {
			log.Println("Failed To Marshal", "err", err)
		} else {
			log.Println("Success To Send", result)
		}
}


func (s *Service)  PublishEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error) {
	return s.repository.Kafka.PublishEvent(topic, value, ch)
}

func (s *Service) ServerSet(ip string, available bool) error {
	if err := s.repository.ServerSet(ip, available); err != nil {
		log.Println("failed to ServerSet:","ip",ip, "available", available)
		return err
	} else {
		return nil
	}
}
func (s *Service) InsertChatting(user, message, roomName string){

	
	if  err := s.repository.InsertChatting(user, message, roomName); err !=nil {
		//TODO: 
		log.Println("fail to chat",err)
	}
}

func (s *Service) EnterRoom(roomName string) ([]*schema.Chat, error) {
	log.Println("roomName:",roomName)
	if res, err := s.repository.GetChatList(roomName); err !=nil {
		log.Println("..Failed to get chat List", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}

}

func (s *Service) MakeRoom(name string) error{

	log.Println("name=", name)
	if  err := s.repository.MakeRoom(name); err !=nil {
		log.Println("..Failed to make room ", "err", err.Error())
		return err
	} else {
		return nil
	}
	
}

func (s *Service) RoomList() ([]*schema.Room, error) {
	if res, err := s.repository.RoomList(); err !=nil {
		log.Println("..Failed to get all room List", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}

}
func (s *Service) Room(name string)(*schema.Room, error)  {
	
	if res, err := s.repository.Room(name); err !=nil {
		log.Println("..Failed to get room", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}
//GetChatList