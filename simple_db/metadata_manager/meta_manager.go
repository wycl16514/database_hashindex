package metadata_manager

import (
	"query"
	rm "record_manager"
	"tx"
)

type MetaDataManager struct {
	tblMgr   *TableManager
	viewMgr  *ViewManager
	statMgr  *StatManager
	constant *query.Constant
	//索引管理器以后再处理
	idxMgr *IndexManager
}

func NewMetaDataManager(isNew bool, tx *tx.Transation) *MetaDataManager {
	metaDataMgr := &MetaDataManager{
		tblMgr:   NewTableManager(isNew, tx),
		constant: nil,
	}

	metaDataMgr.viewMgr = NewViewManager(isNew, metaDataMgr.tblMgr, tx)
	metaDataMgr.statMgr = NewStatManager(metaDataMgr.tblMgr, tx)
	metaDataMgr.idxMgr = NewIndexManager(isNew, metaDataMgr.tblMgr,
		metaDataMgr.statMgr, tx)

	return metaDataMgr
}

func (m *MetaDataManager) CreateTable(tblName string, sch *rm.Schema, tx *tx.Transation) {
	m.tblMgr.CreateTable(tblName, sch, tx)
}

func (m *MetaDataManager) CreateView(viewName string, viewDef string, tx *tx.Transation) {
	m.viewMgr.CreateView(viewName, viewDef, tx)
}

func (m *MetaDataManager) GetLayout(tblName string, tx *tx.Transation) *rm.Layout {
	return m.tblMgr.GetLayout(tblName, tx)
}

func (m *MetaDataManager) GetViewDef(viewName string, tx *tx.Transation) string {
	return m.viewMgr.GetViewDef(viewName, tx)
}

func (m *MetaDataManager) GetStatInfo(tblName string, layout *rm.Layout, tx *tx.Transation) *StatInfo {
	return m.statMgr.GetStatInfo(tblName, layout, tx)
}

func (m *MetaDataManager) CreateIndex(idxName string, tblName string, fldName string, tx *tx.Transation) {
	m.idxMgr.CreateIndex(idxName, tblName, fldName, tx)
}

func (m *MetaDataManager) GetIndexInfo(tblName string, tx *tx.Transation) map[string]*IndexInfo {
	return m.idxMgr.GetIndexInfo(tblName, tx)
}
