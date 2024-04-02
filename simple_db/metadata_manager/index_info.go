package metadata_manager

import (
	rm "record_manager"
	"tx"
)

type IndexInfo struct {
	idxName   string
	fldName   string
	tblSchema *rm.Schema
	tx        *tx.Transation
	idxLayout *rm.Layout
	si        *StatInfo
}

func NewIndexInfo(idxName string, fldName string, tblSchema *rm.Schema,
	tx *tx.Transation, si *StatInfo) *IndexInfo {
	idxInfo := &IndexInfo{
		idxName:   idxName,
		fldName:   fldName,
		tx:        tx,
		tblSchema: tblSchema,
		si:        si,
		idxLayout: nil,
	}

	idxInfo.idxLayout = idxInfo.CreateIdxLayout()

	return idxInfo
}

func (i *IndexInfo) Open() Index {
	//在这里 构建不同的哈希算法对象 s
	return NewHashIndex(i.tx, i.idxName, i.idxLayout)
}

func (i *IndexInfo) BlocksAccessed() int {
	rpb := int(i.tx.BlockSize()) / i.idxLayout.SlotSize()
	numBlocks := i.si.RecordsOutput() / rpb
	return HashIndexSearchCost(numBlocks, rpb)
}

func (i *IndexInfo) RecordsOutput() int {
	return i.si.RecordsOutput() / i.si.DistinctValues(i.fldName)
}

func (i *IndexInfo) DistinctValues(fldName string) int {
	if i.fldName == fldName {
		return 1
	}

	return i.si.DistinctValues(fldName)
}

func (i *IndexInfo) CreateIdxLayout() *rm.Layout {
	sch := rm.NewSchema()
	sch.AddIntField("block")
	sch.AddIntField("id")
	if i.tblSchema.Type(i.fldName) == rm.INTEGER {
		sch.AddIntField("dataval")
	} else {
		fldLen := i.tblSchema.Length(i.fldName)
		sch.AddStringField("dataval", fldLen)
	}

	return rm.NewLayoutWithSchema(sch)
}
