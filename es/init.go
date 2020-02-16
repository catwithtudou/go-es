package es

/**
 * user: ZY
 * Date: 2020/2/16 21:06
 */

import (
	"context"
	"github.com/olivere/elastic/v7"
	"log"
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
	client *elastic.Client
	ctx    context.Context
)





//EsInit 初始化ES连接获取客户端
func EsInit() {
	ctx = context.Background()

	client, err := elastic.NewClient(elastic.SetURL(serverUrl))
	if err != nil {
		log.Println(err)
		return
	}

	//检查索引是否存在
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err!=nil{
		log.Println(err)
		return
	}
	//若不存在则创造索引
	if !exists{
		//生成映射,创造索引
		_,err:=client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err!=nil{
			log.Println(err)
			return
		}
	}

}


