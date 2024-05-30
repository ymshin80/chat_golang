package repository

import (
	"chat_server_golang/config"
	"chat_server_golang/types/schema"
	"database/sql"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Repository struct {
	cfg *config.Config
	db *sql.DB
}
const (
	room = "chatting.room"
	chat = "chatting.chat"
	serverInfo = "chatting.serverInfo"
)
func NewRepository(cfg *config.Config) (*Repository, error) {
	r := &Repository{cfg: cfg}

	var err error
	
	if r.db, err = sql.Open(cfg.DB.Database, cfg.DB.URL); err != nil  {
		return nil, err
	} else {
		return r,nil
	}
}

func (s *Repository) InsertChatting(user, message, roomName string) error {
	log.Println("insert chatting using wss","from", user, "message", message, "roomName", roomName)
	_, err := s.db.Exec("INSERT INTO chatting.chat(room, name, message)  VALUES( ?, ?, ?)",roomName,user, message)
	
	return err
}


//방관리 , 채팅 데이터 관리
func (s *Repository) Room(name string)(*schema.Room, error) {
	d := new(schema.Room)

	qs := query([]string{"SELECT * FROM", room ,"WHERE name= ?"})
	err := s.db.QueryRow(qs, name).Scan(
		&d.ID,
		&d.Name,
		&d.CreateAt,
		&d.UpdatedAt,
	)	

	return d, err

}
//챗팅 정보
func (s *Repository) GetChatList(roomName string) ([]*schema.Chat, error) {
	
	log.Println("roomName:::",roomName)

	qs := query([]string{"SELECT * FROM ", chat, "WHERE room = ? ORDER BY `sendDtm` DESC LIMIT 10" })
	if cursor, err := s.db.Query(qs, roomName); err != nil {
		return nil, err
	} else {
		defer cursor.Close()

		var result []*schema.Chat

		for cursor.Next() {
			d :=new(schema.Chat)
			if err = cursor.Scan(&d.ID,&d.Room, &d.Name,&d.Message,&d.SendDtm); err != nil {
				return nil, err
			} else {
				result = append(result, d)
			}
		}
		if len(result) == 0  {
			return []*schema.Chat{}, nil
		} else {
			return result, nil
		}
		
	}
}


//방 만들기
func (s *Repository) MakeRoom(name string) error {
	_, err := s.db.Exec("INSERT INTO chatting.room(name) values (?)", name)
	
	return err
}
//func [[참조]] [[functionName]]({{파라미터}}) [[return type]]{}
//public 함수면 대문자 시작, private 소문자로 시작
func (s *Repository) RoomList() ([]*schema.Room, error) {
	qs := query([]string{"SELECT * FROM", room })
	
	if cursor, err := s.db.Query(qs); err != nil {
		return nil, err
	} else {
		defer cursor.Close()

		var result []*schema.Room

		for cursor.Next() {
			d :=new(schema.Room)
			if err = cursor.Scan(&d.ID,&d.Name,&d.CreateAt,&d.UpdatedAt,); err != nil {
				return nil, err
			} else {
				result = append(result, d)
			}
		}
		if len(result) == 0  {
			return []*schema.Room{}, nil
		} else {
			return result, nil
		}
		
	}

}
func query(qs []string) string {
	return strings.Join(qs, " ") +";"
}