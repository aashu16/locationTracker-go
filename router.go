package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx/types"
	uuid "github.com/satori/go.uuid"
)

func (env *Env) CreateRouter() *gin.Engine {
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

	// r.GET("/form", env.getForm)

	return r
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

type Uid struct {
	Uid string `form:"uid"`
}

// func (e *Env) getForm(c *gin.Context) {
func (e *Env) getForm() *Response {
	type Formtype struct {
		Uid      uuid.UUID      `db:"event_uid"`
		FormJSON types.JSONText `db:"fields"`
	}
	var forms []Formtype

	if err := e.db.Select(&forms, "SELECT event_uid, event_links -> 'fields' AS fields FROM event_test"); err != nil {
		fmt.Println("Cannot fetch form from database", err)
	}
	fmt.Printf("%v", forms)

	var Form = make([]Forms, len(forms))

	for i := range forms {
		Form[i].EventUID = forms[i].Uid
		if err := json.Unmarshal([]byte(forms[i].FormJSON), &Form[i].Fields); err != nil {
			fmt.Println("Cannot unmarshal JSON forms...", err)
		}
	}

	// if err := json.Unmarshal([]byte(fmt.Sprintf("%v", formJSON)), &form); err != nil {
	// 	fmt.Println("Cannot unmarshal JSON...", err)
	// }

	// for _, val := range formJSON {
	// 	if err := json.Unmarshal([]byte(fmt.Sprintf("%v", val)), &single); err != nil {
	// 		fmt.Println("Cannot unmarshal JSON...", err)
	// 	}
	// 	form = append(form, single)
	// }

	return &Response{Payload: Form}
	// c.JSON(http.StatusOK, &Response{Payload: Form})
}
