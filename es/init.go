package es

/**
 * user: ZY
 * Date: 2020/2/16 21:06
 */

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go-es/model"
	"log"
	"strconv"
)


const (
	//es中的mapping映射
	mapping = `
{
  "mappings": {
    "properties": {
      "title": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_max_word",
        "fielddata": "true",
        "term_vector": "with_positions_offsets"
      },
      "content": {
        "type": "text",
        "analyzer": "ik_max_word",
        "search_analyzer": "ik_max_word",
        "fielddata": "true",
        "term_vector": "with_positions_offsets"
      }
    }
  }
}
`
	//es的url
	serverUrl = "http://localhost:9200/"
	//索引名称
	indexName = "message"

)

//ES客户端
var (
	Client *elastic.Client
	Ctx    context.Context
)





//EsInit 初始化ES连接获取客户端
func EsInit(reminders []model.Reminder) {
	Ctx = context.Background()

	Client, err := elastic.NewClient(elastic.SetURL(serverUrl))
	if err != nil {
		log.Println(err)
		return
	}

	//检查索引是否存在
	exists, err := Client.IndexExists(indexName).Do(Ctx)
	if err!=nil{
		log.Println(err)
		return
	}
	//若不存在则创造索引
	if !exists{
		//生成映射,创造索引
		_,err:=Client.CreateIndex(indexName).BodyString(mapping).Do(Ctx)
		if err!=nil{
			log.Println(err)
			return
		}
	}


	//插入数据
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




	//分析数据
	termQuery := elastic.NewTermsAggregation().Size(100).Field("title")
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
		fmt.Printf("Count:%d\n",v.DocCount)
	}
}


