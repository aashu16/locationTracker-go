package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/types"
	jsoniter "github.com/json-iterator/go"
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

type Clubs struct {
	Id       int       `db:"club_id" json:"club_id"`
	ClubUid  uuid.UUID `db:"club_uid"`
	ClubName string    `db:"club_name"`
}

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
	Fields []Field `json:"fields"`
}

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=entrysport sslmode=disable")
	if err != nil {
		log.Fatalln("Cannot connect to database...")
	}
	var env = &Env{db: db}

	fmt.Println("No error, mate")

	var clubs []Clubs
	if err = db.Select(&clubs, "SELECT * FROM bookings.clubs WHERE club_id < 6883"); err != nil {
		fmt.Println("Cannot fetch records from database...")
	}
	fmt.Println(clubs)

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.SecureJSON(http.StatusOK, gin.H{
			"message": "Hello, world!",
			"status":  "success",
		})
	})

	r.GET("/getLocation/:code", getLocation)

	r.POST("/setLocation/:code", setLocation)

	r.GET("/clubs", env.getClubs)

	r.POST("/events", env.createEvent)

	r.POST("/form", env.enterForm)

	r.GET("/form", env.getForm)

	_ = r.Run(":8080")
}

func getLocation(c *gin.Context) {
	locationCode := c.Param("code")
	reply, err := redisClient.Get(locationCode).Result()
	if err != nil {
		fmt.Println("The key does not exist")

		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": "failed",
		})
	}

	coords := strings.Split(reply, ",")

	c.SecureJSON(http.StatusOK, gin.H{
		"status": "success",
		"lat":    coords[0],
		"lng":    coords[1],
	})
}

func setLocation(c *gin.Context) {
	locationCode := c.Param("code")

	var coordinates Coords
	if err := c.ShouldBindJSON(&coordinates); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": "failed",
		})
	}

	serialisedCoords := coordinates.Lat + "," + coordinates.Lng
	if err := redisClient.Set(locationCode, serialisedCoords, 0).Err(); err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"status": "failed",
		})
	}

	c.SecureJSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

func (e *Env) getClubs(c *gin.Context) {
	var clubs []Clubs
	if err := e.db.Select(&clubs, "SELECT * FROM bookings.clubs WHERE club_id < 6883"); err != nil {
		fmt.Println("Cannot fetch records from database...")
	}

	c.JSON(http.StatusOK, clubs)
}

func (e *Env) createEvent(c *gin.Context) {
	// var event Events
	var eventLinks types.JSONText

	if err := eventLinks.UnmarshalJSON([]byte(`{"first": "success!","second": "success!!","third": "success!!!"}`)); err != nil {
		fmt.Println("Cannot insert JSON...")
	}
	fmt.Println(eventLinks)
	result, err := e.db.Exec("INSERT INTO event_test (event_links) VALUES ($1)", eventLinks)
	if err != nil {
		fmt.Println("Cannot create a new event...", err)
	}

	fmt.Println(result)
}

func (e *Env) enterForm(c *gin.Context) {
	var form Forms
	// var formJSON types.JSONText
	if err := c.ShouldBindJSON(&form); err != nil {
		fmt.Println("Error parsing form fields...")
	}
	fmt.Println(form)

	preparedJSON, err := json.Marshal(form)
	if err != nil {
		fmt.Println("Cannot marshal JSON...")
	}

	//if err := formJSON.UnmarshalJSON([]byte(preparedJSON)); err != nil {
	//fmt.Println("Cannot get JSON from request...", err)
	//}
	//fmt.Println(formJSON)

	result, err := e.db.Exec("INSERT INTO event_test (event_links) VALUES ($1)", preparedJSON)
	if err != nil {
		fmt.Println("Cannot insert form into database", err)
	}
	fmt.Println(result)
}

func (e *Env) getForm(c *gin.Context) {
	var form Forms
	var formJSON []types.JSONText
	if err := e.db.Select(&formJSON, "SELECT event_links FROM event_test LIMIT 1"); err != nil {
		fmt.Println("Cannot fetch form from database", err)
	}
	fmt.Println(formJSON[0])

	if err := json.Unmarshal([]byte(fmt.Sprintf("%v", formJSON[0])), &form); err != nil {
		fmt.Println("Cannot unmarshal JSON...", err)
	}

	c.JSON(http.StatusOK, form)
}
