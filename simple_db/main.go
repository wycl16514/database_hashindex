package main

import (
	bmg "buffer_manager"
	fm "file_manager"
	"fmt"
	lm "log_manager"
	metadata_manager "metadata_management"
	"parser"
	"planner"
	"query"
	"tx"
)

func PrintStudentTable(tx *tx.Transation, mdm *metadata_manager.MetaDataManager) {
	queryStr := "select name, majorid, gradyear from STUDENT"
	p := parser.NewSQLParser(queryStr)
	queryData := p.Query()
	test_planner := planner.CreateBasicQueryPlanner(mdm)
	test_plan := test_planner.CreatePlan(queryData, tx)
	test_interface := (test_plan.Open())
	test_scan, _ := test_interface.(query.Scan)
	for test_scan.Next() {
		fmt.Printf("name: %s, majorid: %d, gradyear: %d\n",
			test_scan.GetString("name"), test_scan.GetInt("majorid"),
			test_scan.GetInt("gradyear"))
	}
}

func CreateInsertUpdateByUpdatePlanner() {
	file_manager, _ := fm.NewFileManager("student", 2048)
	log_manager, _ := lm.NewLogManager(file_manager, "logfile.log")
	buffer_manager := bmg.NewBufferManager(file_manager, log_manager, 3)
	tx := tx.NewTransation(file_manager, log_manager, buffer_manager)
	mdm := metadata_manager.NewMetaDataManager(file_manager.IsNew(), tx)

	updatePlanner := planner.NewBasicUpdatePlanner(mdm)
	createTableSql := "create table STUDENT (name varchar(16), majorid int, gradyear int)"
	p := parser.NewSQLParser(createTableSql)
	tableData := p.UpdateCmd().(*parser.CreateTableData)
	updatePlanner.ExecuteCreateTable(tableData, tx)

	insertSQL := "insert into STUDENT (name, majorid, gradyear) values(\"tylor\", 30, 2020)"
	p = parser.NewSQLParser(insertSQL)
	insertData := p.UpdateCmd().(*parser.InsertData)
	updatePlanner.ExecuteInsert(insertData, tx)
	insertSQL = "insert into STUDENT (name, majorid, gradyear) values(\"tom\", 35, 2023)"
	p = parser.NewSQLParser(insertSQL)
	insertData = p.UpdateCmd().(*parser.InsertData)
	updatePlanner.ExecuteInsert(insertData, tx)

	fmt.Println("table after insert:")
	PrintStudentTable(tx, mdm)

	updateSQL := "update STUDENT set majorid=20 where majorid=30 and gradyear=2020"
	p = parser.NewSQLParser(updateSQL)
	updateData := p.UpdateCmd().(*parser.ModifyData)
	updatePlanner.ExecuteModify(updateData, tx)

	fmt.Println("table after update:")
	PrintStudentTable(tx, mdm)

	deleteSQL := "delete from STUDENT where majorid=35"
	p = parser.NewSQLParser(deleteSQL)
	deleteData := p.UpdateCmd().(*parser.DeleteData)
	updatePlanner.ExecuteDelete(deleteData, tx)

	fmt.Println("table after delete")
	PrintStudentTable(tx, mdm)
}

func TestIndex() {
	file_manager, _ := fm.NewFileManager("student", 4096)
	log_manager, _ := lm.NewLogManager(file_manager, "logfile.log")
	buffer_manager := bmg.NewBufferManager(file_manager, log_manager, 3)
	tx := tx.NewTransation(file_manager, log_manager, buffer_manager)
	fmt.Printf("file manager is new: %v\n", file_manager.IsNew())
	mdm := metadata_manager.NewMetaDataManager(file_manager.IsNew(), tx)

	//创建 student 表，并插入一些记录
	updatePlanner := planner.NewBasicUpdatePlanner(mdm)
	createTableSql := "create table STUDENT (name varchar(16), majorid int, gradyear int)"
	p := parser.NewSQLParser(createTableSql)
	tableData := p.UpdateCmd().(*parser.CreateTableData)
	updatePlanner.ExecuteCreateTable(tableData, tx)

	insertSQL := "insert into STUDENT (name, majorid, gradyear) values(\"tylor\", 30, 2020)"
	p = parser.NewSQLParser(insertSQL)
	insertData := p.UpdateCmd().(*parser.InsertData)
	updatePlanner.ExecuteInsert(insertData, tx)
	insertSQL = "insert into STUDENT (name, majorid, gradyear) values(\"tom\", 35, 2023)"
	p = parser.NewSQLParser(insertSQL)
	insertData = p.UpdateCmd().(*parser.InsertData)
	updatePlanner.ExecuteInsert(insertData, tx)

	fmt.Println("table after insert:")
	PrintStudentTable(tx, mdm)
	//在 student 表的 majorid 字段建立索引
	mdm.CreateIndex("majoridIdx", "STUDENT", "majorid", tx)
	//查询建立在 student 表上的索引并根据索引输出对应的记录信息
	studetPlan := planner.NewTablePlan(tx, "STUDENT", mdm)
	updateScan := studetPlan.Open().(*query.TableScan)
	//先获取每个字段对应的索引对象,这里我们只有 majorid 建立了索引对象
	indexes := mdm.GetIndexInfo("STUDENT", tx)
	//获取 majorid 对应的索引对象
	majoridIdxInfo := indexes["majorid"]
	//将改rid 加入到索引表
	majorIdx := majoridIdxInfo.Open()
	updateScan.BeforeFirst()
	for updateScan.Next() {
		//返回当前记录的 rid
		dataRID := updateScan.GetRid()
		dataVal := updateScan.GetVal("majorid")
		majorIdx.Insert(dataVal, dataRID)
	}

	//通过索引表获得给定字段内容的记录
	majorid := 35
	majorIdx.BeforeFirst(query.NewConstantWithInt(&majorid))
	for majorIdx.Next() {
		datarid := majorIdx.GetDataRID()
		updateScan.MoveToRid(datarid)
		fmt.Printf("student name :%s, id: %d\n", updateScan.GetScan().GetString("name"),
			updateScan.GetScan().GetInt("majorid"))
	}

	majorIdx.Close()
	updateScan.GetScan().Close()
	tx.Commit()

}

func main() {
	TestIndex()
}
