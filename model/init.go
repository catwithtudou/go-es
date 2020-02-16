package model

/**
 * user: ZY
 * Date: 2020/2/16 20:58
 */

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var DB *gorm.DB

func ModelInit(){
	var err error
	DB,err = gorm.Open("mysql","root:@tcp(127.0.0.1:3306)/magipoke?parseTime=true&charset=utf8&loc=Local")
	if err != nil {
		log.Println("init the mysql failed")
		return
	}

	DB.DB().SetMaxOpenConns(10)
	DB.DB().SetMaxIdleConns(100)
}

func ModelClose(){
	if DB!=nil{
		err:=DB.Close()
		if err!=nil{
			log.Println("close the mysql failed")
			return
		}
	}
}