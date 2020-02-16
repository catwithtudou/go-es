package main

import (
	"go-es/es"
	"go-es/model"
)

/**
 * user: ZY
 * Date: 2020/2/16 22:09
 */

func main(){
	model.ModelInit()
	defer model.ModelClose()
	es.EsInit()
	modelData,_:=model.SelectAll()
	es.WriteData(modelData)
}