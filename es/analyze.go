package es

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"go-es/model"
	"log"
	"strconv"
	"sync"
)

/**
 * user: ZY
 * Date: 2020/2/16 21:19
 */

// WriteData 将model中数据写入es
// 考虑到model中数据过多,使用单线
// 程速度较慢,所以采用分桶的方式通
// 过调度协程来插入
// 后面发现库中有批量写入的接口故还
// 是采用库中的接口
func WriteData(reminders []model.Reminder) {
	////分桶数量设置为100
	//bucketNum:=100
	//
	////将model数据切割后插入桶中
	//bucketData:=make([][]model.Reminder,bucketNum+1)
	//dataNum := len(reminders)
	//bucketLen := dataNum/bucketNum
	//lastLoc,j := 0, 0
	//for i:=bucketLen;i<=dataNum && j<bucketNum;i+=bucketLen{
	//	bucketData[j]=reminders[lastLoc:i]
	//	j++
	//	lastLoc = i
	//}
	//if lastLoc+bucketLen > dataNum{
	//	bucketData[j] = reminders[lastLoc:]
	//}
	//
	//
	////生成互斥锁
	//var waitGroup sync.WaitGroup
	//
	////开启协程往es中写入数据
	//for i:=0;i<bucketNum+1;i++{
	//	waitGroup.Add(1)
	//	go WriteDocument(bucketData[i],&waitGroup)
	//}
	//
	//waitGroup.Wait()
	//
	//return

	//初始化批量操作接口
	bulkRequest := Client.Bulk()

	for _, v := range reminders {
		doc := elastic.NewBulkIndexRequest().Index(indexName).Id(strconv.Itoa(int(v.ID))).Doc(v)
		bulkRequest = bulkRequest.Add(doc)
	}

	//执行操作
	response, err := bulkRequest.Do(Ctx)
	if err != nil {
		log.Println(err)
		return
	}

	//打印失败次数
	failed := response.Failed()
	l := len(failed)
	if l > 0 {
		fmt.Printf("Error(%d)", l, response.Errors)
	}




}

//ModelData 向es中插入文档
func WriteDocument(reminders []model.Reminder, wait *sync.WaitGroup) (err error) {
	Ctx:=context.Background()
	defer wait.Done()
	for _, v := range reminders {
		_, err = Client.Index().
			Index(indexName).
			Id(strconv.Itoa(int(v.ID))).
			BodyJson(v).
			Do(Ctx)
		if err != nil {
			log.Println(err)
		}
	}
	return
}

type AnalyzeItem struct {
	DocCountError int      `json:"doc_count_error_upper_bound"`
	SumOtherDoc   int      `json:"sum_other_doc_count"`
	Buckets       []Bucket `json:"buckets"`
}
type Bucket struct {
	Key      string `json:"key"`
	DocCount int    `json:"doc_count"`
}

//AnalyzeData 分析文档中title的热词
func AnalyzeData() (err error) {

	termQuery := elastic.NewTermsAggregation().Size(100).Field("title").Include("[\u4E00-\u9FA5][\u4E00-\u9FA5]")
	result, err := Client.Search().Index(indexName).Aggregation("messages", termQuery).Do(Ctx)
	if err != nil {
		log.Println(err)
		return
	}

	data := &AnalyzeItem{}
	messages, _ := result.Aggregations["messages"].MarshalJSON()
	err = json.Unmarshal(messages, &data)
	if err != nil {
		log.Println(err)
		return
	}

	for _, v := range data.Buckets {
		fmt.Println("Key:" + v.Key)
		fmt.Println("Count:" + string(v.DocCount))
	}

	return
}
