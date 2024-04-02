package metadata_manager

import (
	"query"
	rm "record_manager"
	"tx"
)

type TableManagerInterface interface {
	CreateTable(tblName string, sch *rm.Schema, tx *tx.Transation)
	GetLayout(tblName string, tx *tx.Transation) *rm.Layout
}

type Index interface {
	//指向第一条满足查询条件的记录
	BeforeFirst(searchKey *query.Constant)
	//是否还有其余满足条件的记录
	Next() bool
	GetDataRID() *rm.RID
	Insert(dataval *query.Constant, datarid *rm.RID)
	Delete(dataval *query.Constant, datarid *rm.RID)
	Close()
}
