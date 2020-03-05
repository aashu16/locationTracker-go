package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-redis/redis"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	jsoniter "github.com/json-iterator/go"

	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
	uuid "github.com/satori/go.uuid"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var redisClient *redis.Client

type Env struct {
	db *sqlx.DB
}

type Coords struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type Club struct {
	Id       int       `db:"club_id" json:"club_id"`
	ClubUid  uuid.UUID `db:"club_uid"`
	ClubName string    `db:"club_name"`
}

type Clubs []Club

type Events struct {
	Id        uuid.UUID       `db:"event_uid" json:"event_id"`
	EventDesc *sql.NullString `db:"event_description" json:"event_desc,omitempty"`
	// EventWebsite *sql.NullString     `db:"event_website" json:"event_website,omitempty"`
	// EventFB      *sql.NullString     `db:"event_fb" json:"event_db,omitempty"`
	// EventTwitter *sql.NullString     `db:"event_twitter" json:"event_twitter,omitempty"`
	EventLinks *types.NullJSONText `db:"event_links" json:"event_links,omitempty"`
}

type Field struct {
	Name  string      `json:"field_name"`
	Value interface{} `json:"field_value"`
}

type Forms struct {
	EventUID uuid.UUID `json:"event_id"`
	Fields   []Field   `json:"fields"`
}

type Response struct {
	Payload interface{} `json:"data"`
}

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	db, err := sqlx.Connect("pgx", "user=postgres password=postgres dbname=entrysport sslmode=disable")
	if err != nil {
		log.Fatalln("Cannot connect to database...")
	}
	var env = &Env{db: db}

	fmt.Println("No error, mate")

	var clubs Clubs
	if err = db.Select(&clubs, "SELECT * FROM bookings.clubs WHERE club_id < 6883"); err != nil {
		fmt.Println("Cannot fetch records from database...")
	}
	fmt.Println(clubs)

	r := env.CreateRouter()

	_ = r.Run(":8080")
}
