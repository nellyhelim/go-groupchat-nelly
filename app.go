package main

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/lolmourne/go-groupchat/resource"
	"github.com/lolmourne/go-groupchat/usecase/groupchat"
	"github.com/lolmourne/go-groupchat/usecase/userauth"
)

var db *sqlx.DB
var dbResource resource.DBItf
var userAuthUsecase userauth.UsecaseItf
var groupchatUsecase groupchat.UsecaseItf

func main() {
	dbInit, err := sqlx.Connect("postgres", "host=34.101.216.10 user=skilvul password=skilvul123apa dbname=skilvul-groupchat sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	dbRsc := resource.NewDBResource(dbInit)
	dbResource = dbRsc
	db = dbInit

	userAuthUsecase = userauth.NewUsecase(dbRsc)
	groupchatUsecase = groupchat.NewUseCase(dbRsc)

	r := gin.Default()
	r.POST("/register", register)
	r.POST("/login", login)
	r.GET("/profile/:username", getProfile)
	r.PUT("/profile", updateProfile)
	r.PUT("/password", changePassword)
	r.PUT("/room", joinRoom)
	r.PUT("/editroom", EditGroupchat)
	r.POST("/room", createRoom)
	r.Run()
}

func register(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")
	confirmPassword := c.Request.FormValue("confirm_password")

	err := userAuthUsecase.Register(username, password, confirmPassword)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err:     err.Error(),
			Message: "Failed",
		})
		return
	}

	c.JSON(201, StandardAPIResponse{
		Err:     "null",
		Message: "Success create new user",
	})
}

func login(c *gin.Context) {
	username := c.Request.FormValue("username")
	password := c.Request.FormValue("password")

	user, err := userAuthUsecase.Login(username, password)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err:     err.Error(),
			Message: "Failed",
		})
		return
	}

	c.JSON(200, StandardAPIResponse{
		Data: user,
	})
}

func getProfile(c *gin.Context) {
	username := c.Param("username")

	user, err := dbResource.GetUserByUserName(username)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: "Unauthorized",
		})
		return
	}

	resp := User{
		Username:   user.Username,
		ProfilePic: user.ProfilePic,
		CreatedAt:  user.CreatedAt.UnixNano(),
	}

	c.JSON(200, StandardAPIResponse{
		Err:  "null",
		Data: resp,
	})
}

func updateProfile(c *gin.Context) {
	username := c.Request.FormValue("username")
	profilepic := c.Request.FormValue("imageURL")

	err := dbResource.UpdateProfile(username, profilepic)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: err.Error(),
		})
		return
	}

	c.JSON(201, StandardAPIResponse{
		Err:     "null",
		Message: "Success update profile picture",
	})

}

func changePassword(c *gin.Context) {
	username := c.Request.FormValue("username")
	oldpass := c.Request.FormValue("old_password")
	newpass := c.Request.FormValue("new_password")

	user, err := dbResource.GetUserByUserName(username)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err: err.Error(),
		})
		return
	}

	oldpass += user.Salt
	h := sha256.New()
	h.Write([]byte(oldpass))
	hashedOldPassword := fmt.Sprintf("%x", h.Sum(nil))

	if user.Password != hashedOldPassword {
		c.JSON(401, StandardAPIResponse{
			Err: "old password is wrong!",
		})
		return
	}

	//new pass
	salt := RandStringBytes(32)
	newpass += salt

	h = sha256.New()
	h.Write([]byte(newpass))
	hashedNewPass := fmt.Sprintf("%x", h.Sum(nil))

	err2 := dbResource.UpdateUserPassword(username, hashedNewPass)

	if err2 != nil {
		c.JSON(400, StandardAPIResponse{
			Err: err.Error(),
		})
		return
	}

	c.JSON(201, StandardAPIResponse{
		Err:     "null",
		Message: "Success update password",
	})

}

func createRoom(c *gin.Context) {
	name := c.Request.FormValue("name")
	description := c.Request.FormValue("description")
	categoryID := c.Request.FormValue("category_id")
	adminID := c.Request.FormValue("admin_id")

	groupchat, err := groupchatUsecase.CreateGroupchat(name, adminID, description, categoryID)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err:     err.Error(),
			Message: "Failed make a groupchat",
		})
		return
	}

	c.JSON(201, StandardAPIResponse{
		Err:     "null",
		Message: "Success create new room",
		Data:    groupchat,
	})
}

func EditGroupchat(c *gin.Context) {
	name := c.Request.FormValue("name")
	description := c.Request.FormValue("description")
	categoryID := c.Request.FormValue("category_id")

	groupchat, err := groupchatUsecase.EditGroupchat(name, description, categoryID)
	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err:     err.Error(),
			Message: "Failed edit the groupchat",
		})
		return
	}

	c.JSON(201, StandardAPIResponse{
		Err:     "null",
		Message: "Edit room Success",
		Data:    groupchat,
	})
}

func joinRoom(c *gin.Context) {
	roomID := c.Request.FormValue("room_id")
	userID := c.Request.FormValue("user_id")

	err := groupchatUsecase.JoinRoom(roomID, userID)

	if err != nil {
		c.JSON(400, StandardAPIResponse{
			Err:     err.Error(),
			Message: "Failed make a groupchat",
		})
		return
	}

	c.JSON(201, StandardAPIResponse{
		Err:     "null",
		Message: "Success join to room with ID " + roomID,
	})
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type StandardAPIResponse struct {
	Err     string      `json:"err"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type User struct {
	Username   string `json:"username"`
	ProfilePic string `json:"profile_pic"`
	CreatedAt  int64  `json:"created_at"`
}

type UserDB struct {
	UserID     sql.NullInt64  `db:"user_id"`
	UserName   sql.NullString `db:"username"`
	ProfilePic sql.NullString `db:"profile_pic"`
	Salt       sql.NullString `db:"salt"`
	Password   sql.NullString `db:"password"`
	CreatedAt  time.Time      `db:"created_at"`
}

//TODO complete all API request
type RoomDB struct {
	RoomID      sql.NullInt64  `db:room_id`
	Name        sql.NullString `db:name`
	Admin       sql.NullInt64  `db:admin_user_id`
	Description sql.NullString `db:description`
	CategoryID  sql.NullInt64  `db:category_id`
	CreatedAt   time.Time      `db:"created_at"`
}

type Room struct {
	RoomID      int64  `json:"room_id"`
	Name        string `json:"name"`
	Admin       int64  `json:"admin"`
	Description string `json:"description"`
}
