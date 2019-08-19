package domainCheckDb

import (
	"database/sql"
	domainCheckData "domainCheck/src/data"
	"strconv"
	"time"

	"github.com/coderguang/GameEngine_go/sgtime"

	"github.com/coderguang/GameEngine_go/sgwhois"

	"github.com/coderguang/GameEngine_go/sglog"
	"github.com/coderguang/GameEngine_go/sgmysql"
)

var globalMysqldb *sql.DB
var globalMysqlStmt *sql.Stmt

func InitDbConnection() {

	dbuser, dbpwd, dburl, dbport, dbname := domainCheckData.GetDbConnectionInfo()

	db, err := sgmysql.Open(dbuser, dbpwd, dburl, dbport, dbname, "utf8")
	if err != nil {
		sglog.Fatal("connect to db error")
		return
	}

	globalMysqldb = db

	sqlstr := "replace into domain_info (zone,name,name_length,status,create_dt_str,update_dt_str,expiry_dt_str,create_dt,update_dt,expiry_dt) values(?,?,?,?,?,?,?,?,?,?)"

	stmt, err := globalMysqldb.Prepare(sqlstr)

	if err != nil {
		sglog.Fatal("stmt to db error,err=%e", err)
		return
	}
	globalMysqlStmt = stmt
}

func InsertOrUpdateData(result *sgwhois.Whois) {
	_, err := globalMysqlStmt.Exec(result.Zone, result.Name, len(result.Name), result.IsRegist, result.CreateDtStr, result.UpdateDtStr, result.ExpiryDtStr, result.CreateDt.NormalString(), result.UpdateDt.NormalString(), result.ExpiryDt.NormalString())
	if err != nil {
		sglog.Error("stmt to db error,%e", err)
	}
}

func transFiledToStruct(fieldName string, fieldValue string, result *sgwhois.Whois) {
	switch fieldName {
	case "name":
		result.Name = fieldValue
	case "zone":
		result.Zone = fieldValue
	case "name_length":
		result.Name_length, _ = strconv.Atoi(fieldValue)
	case "status":
		result.IsRegist, _ = strconv.Atoi(fieldValue)
	case "create_dt_str":
		result.CreateDtStr = fieldValue
	case "update_dt_str":
		result.UpdateDtStr = fieldValue
	case "expiry_dt_str":
		result.ExpiryDtStr = fieldValue
	case "create_dt":
		result.CreateDt = sgtime.New()
		result.CreateDt.Parse(fieldValue, sgtime.FORMAT_TIME_NORMAL)
	case "update_dt":
		result.UpdateDt = sgtime.New()
		result.UpdateDt.Parse(fieldValue, sgtime.FORMAT_TIME_NORMAL)
	case "expiry_dt":
		result.ExpiryDt = sgtime.New()
		result.ExpiryDt.Parse(fieldValue, sgtime.FORMAT_TIME_NORMAL)
	default:
		sglog.Error("unkonw column:%s,value:%s", fieldName, fieldValue)
	}
}

func transfromDataToStruct(rows *sql.Rows) *sgwhois.Whois {
	result := new(sgwhois.Whois)
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	_ = rows.Scan(scanArgs...)
	for i, col := range values {
		fieldName := columns[i]
		fieldValue := string(col.([]byte))
		if col != nil {
			transFiledToStruct(fieldName, fieldValue, result)
		}
	}
	result.Domain = result.Name + "." + result.Zone
	return result
}

func GetAllExpireDtData(zone string, name_length int) ([]*sgwhois.Whois, error) {
	sqlStr := "select * from domain_info where "
	time_now := time.Now()
	sqlStr += "(expiry_dt<'" + time_now.String() + "' or status=0) and zone=\"" + zone + "\" and name_length=" + strconv.Itoa(name_length) + " and status!=2"
	rows, rowErr := globalMysqldb.Query(sqlStr)
	if rowErr != nil {
		return nil, rowErr
	}
	defer rows.Close()
	domainList := []*sgwhois.Whois{}
	for rows.Next() {
		result := transfromDataToStruct(rows)
		domainList = append(domainList, result)
	}
	return domainList, nil
}

func GetLasteExpiryDt(zone string, name_length int) (int64, *sgtime.DateTime) {

	now := sgtime.New()
	defaultSleep := 60
	last_sqlStr := "select expiry_dt from domain_info where zone=\"" + zone + "\" and name_length=" + strconv.Itoa(name_length) + " and status!=2 order by expiry_dt limit 1"
	last_rows, errEx := globalMysqldb.Query(last_sqlStr)
	if errEx != nil {
		return int64(defaultSleep), now.Add(defaultSleep)
	}
	defer last_rows.Close()
	for last_rows.Next() {
		result := transfromDataToStruct(last_rows)
		sleepTime := result.ExpiryDt.GetTotalSecond() - now.GetTotalSecond()
		return sleepTime, result.ExpiryDt
	}
	return int64(defaultSleep), now.Add(defaultSleep)
}
