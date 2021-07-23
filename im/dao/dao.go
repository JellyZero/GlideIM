package dao

import "go_im/pkg/db"

func Init() {
	InitUserDao()

	tables := []interface{}{
		&User{},
		&Chat{},
		&ChatMessage{},
	}
	for _, tb := range tables {
		if !db.DB.HasTable(tb) {
			db.DB.CreateTable(&tb)
		}
	}
}