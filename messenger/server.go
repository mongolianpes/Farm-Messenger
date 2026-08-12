package messenger

import (
	"database/sql"
	"fmt"
	pb "messenger/proto"
	"os"
)

type MessengerServer struct {
	pb.UnimplementedMessengerServer
	db *sql.DB
}

func NewMessengerServer() *MessengerServer {
	return &MessengerServer{
		db: connectToDB(),
	}
}

func connectToDB() *sql.DB {
	// host := "localhost"
	// port := "5432"
	// user := "postgres"
	// password := "123"
	// dbname := "project_farm"

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	return db
}
