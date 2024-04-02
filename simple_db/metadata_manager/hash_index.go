package metadata_manager

import (
	"fmt"
	"query"
	rm "record_manager"
	"tx"
)

const (
	NUM_BUCKETS = 100
)

type HashIndex struct {
	tx        *tx.Transation
	idxName   string
	layout    *rm.Layout
	searchKey *query.Constant
	ts        *query.TableScan
}

func NewHashIndex(tx *tx.Transation, idxName string, layout *rm.Layout) *HashIndex {
	return &HashIndex{
		tx:      tx,
		idxName: idxName,
		layout:  layout,
		ts:      nil,
	}
}

func (h *HashIndex) BeforeFirst(searchKey *query.Constant) {
	h.Close()
	h.searchKey = searchKey
	bucket := searchKey.HashCode() % NUM_BUCKETS
	//构造索引记录对应的表名称
	tblName := fmt.Sprintf("%s#%d", h.idxName, bucket)
	h.ts = query.NewTableScan(h.tx, tblName, h.layout)
}

func (h *HashIndex) Next() bool {
	for h.ts.Next() {
		if h.ts.GetVal("dataval").Equals(h.searchKey) {
			return true
		}
	}

	return false
}

func (h *HashIndex) GetDataRID() *rm.RID {
	//返回记录所在的区块信息
	blkNum := h.ts.GetInt("block")
	id := h.ts.GetInt("id")
	return rm.NewRID(blkNum, id)
}

func (h *HashIndex) Insert(val *query.Constant, rid *rm.RID) {
	h.BeforeFirst(val)
	h.ts.Insert()
	h.ts.SetInt("block", rid.BlockNumber())
	h.ts.SetInt("id", rid.Slot())
	h.ts.SetVal("dataval", val)
}

func (h *HashIndex) Delete(val *query.Constant, rid *rm.RID) {
	h.BeforeFirst(val)
	for h.Next() {
		if h.GetDataRID().Equals(rid) {
			h.ts.Delete()
			return
		}
	}
}

func (h *HashIndex) Close() {
	if h.ts != nil {
		h.ts.Close()
	}
}

func HashIndexSearchCost(numBlocks int, rpb int) int {
	return numBlocks / NUM_BUCKETS
}
