package model

/**
 * user: ZY
 * Date: 2020/2/16 20:59
 */

type Reminder struct{
	ID uint `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
}

func SelectAll()(reminders []Reminder,err error){
	if err:=DB.Table("cyxbsmobile_transaction").Select("id, title, content").Limit(1000).Find(&reminders).Error;err!=nil{
		return nil,err
	}
	return
}

