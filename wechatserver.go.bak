package main

import (
	//"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	//"net/http"
	"fmt"
	//	"math/rand"
	//"net/url"
	//	"io"
	"os"
	//"io/ioutil"
	//	"path/filepath"
	//	"strconv"
	"time"

	//"github.com/bitly/go-simplejson"
	_ "github.com/go-sql-driver/mysql"
	//"encoding/json"
	//"./MyDBStructs"
	"github.com/itsjamie/gin-cors"
	//	"github.com/virushuo/Go-Apns"
)

const (
	RootResDir = "./res/"
	WebResDir  = "/res"

	TimeFormat = "2006-01-02 15:04:05"
)

var gDB *gorm.DB
var gRelease bool = true
var gLocal bool = false

// **********************************************************************************************************************
// **********************************************************************************************************************
//
//			Database Structures
//
// **********************************************************************************************************************
// **********************************************************************************************************************
/*
type UserView struct {
	LoginName string `gorm:"column:LoginName"`
	FirstName string `gorm:"column:FirstName"`
	LastName  string `gorm:"column:LastName"`
	Icon      string `gorm:"column:Icon"`
	Score     int    `gorm:"column:Score"`
	//	Authority	int		`gorm:"column:Authority"`
	DemoJamId1   int   `gorm:"column:DemoJamId1"`
	DemoJamId2   int   `gorm:"column:DemoJamId2"`
	VoiceVoteId1 int   `gorm:"column:VoiceVoteId1"`
	VoiceVoteId2 int   `gorm:"column:VoiceVoteId2"`
	EggVoted     bool  `gorm:"column:EggVoted"`
	GreenAmb     bool  `gorm:"column:GreenAmb"`
	SubTime      int64 `gorm:"column:SubTime"`
}
*/
// **********************************************************************************************************************
// **********************************************************************************************************************
//
//			router's selection logic function
//
// **********************************************************************************************************************
// **********************************************************************************************************************
func MainGetRouter(c *gin.Context) {
	MyPrint("http message start!")
	c.Request.ParseForm()
	MyPrint("Request : ", c.Request.Form)
	checkWeChat(c)
	MyPrint("http message finished!")
}

// **********************************************************************************************************************
// **********************************************************************************************************************
//
//			Get Function
//
// **********************************************************************************************************************
// **********************************************************************************************************************
func checkWeChat(c *gin.Context) {
	echostring := c.Query("echostr")
	//js, err := simplejson.NewJson([]byte(`{}`))
	//CheckErr(err)
	//js.Set("echostr", echostring)
	//c.JSON(200, js)
	c.Writer.WriteString(echostring)
}

// **********************************************************************************************************************
// **********************************************************************************************************************
//
//			main function
//
// **********************************************************************************************************************
// **********************************************************************************************************************
func main() {
	argCnt := len(os.Args)

	for i := 1; i < argCnt; i++ {
		if os.Args[i] == "debug" {
			gRelease = false
		} else if os.Args[i] == "local" {
			gLocal = true
		}
	}

	fmt.Println("Release Mode : ", gRelease)

	if gRelease {
		gin.SetMode(gin.ReleaseMode)
	}

	//	gDB = ConnectDB(gRelease)

	//TestFunc()

	MyPrint("start server!")

	//router := gin.Default()

	router := gin.New()

	//authorized := router.Group("/")
	router.Use(cors.Middleware((cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	})))

	router.GET("", MainGetRouter)
	//	router.POST("", RouterPost)

	//	router.Static(WebResDir, RootResDir)
	//	router.StaticFile(WebVersionResDir, VersionResDir)

	//	router.Static(WebDemoJamResDir, DemoJamResDir)
	//	router.Static(WebSapVoiceResDir, SapVoiceResDir)
	//	router.Static(WebEggHikingResDir, EggHikingResDir)
	//	router.Static(WebLuckdrawResDir, LuckdrawResDir)
	//	router.Static(WebGuideResDir, GuideResDir)

	router.Run(":80")

	//gDB.Close()
}

// **********************************************************************************************************************
// **********************************************************************************************************************
//
//			common function
//
// **********************************************************************************************************************
// **********************************************************************************************************************
func MyPrint(a ...interface{}) {
	if gRelease {
		return
	} else {
		fmt.Println(a)
	}
}

func CheckErr(err error) bool {
	if err != nil {
		panic(err)
		return true
	}
	return false
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func CheckDirIsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func ConnectDB(isRelease bool) *gorm.DB {
	MyPrint("start to connecting db!")
	var connectStr string
	if gLocal {
		connectStr = "root@tcp(127.0.0.1:3306)/dreamcastleshanghai?charset=utf8&parseTime=True"
	} else {
		connectStr = "root:1011@/dreamcastleshanghai?charset=utf8&parseTime=True"
	}
	db, err := gorm.Open("mysql", connectStr)
	if CheckErr(err) {
		return nil
	}
	MyPrint("start to connecting db finished!")

	MyPrint("start to init db!")
	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	if isRelease {
		db.LogMode(false)
	} else {
		db.LogMode(true)
	}
	db.SingularTable(true)
	//db.AutoMigrate(&User{}, &Tests{})
	MyPrint("start to init db finished!")

	return &db
}

func readError(errorChan <-chan error) {
	for {
		apnerror := <-errorChan
		fmt.Println(apnerror.Error())
	}
}
