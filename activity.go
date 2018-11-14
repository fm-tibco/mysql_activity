package mysql_activity

import (
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"strings"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// log is the default package logger
var log = logger.GetLogger("activity-mysql")

const (
	ivDataSourceName = "dataSourceName"
	ivQuery          = "query"
	ivParams         = "params"
	ivColumnTypes    = "columnTypes"

	ovResults = "results"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new AppActivity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	dsNameIn := context.GetInput(ivDataSourceName)

	dataSourceName, ok := dsNameIn.(string)
	if !ok || dataSourceName == "" {
		return false, fmt.Errorf("dataSourceName not set")
	}
	log.Debugf("dataSourceName: %s", dataSourceName)

	queryIn := context.GetInput(ivQuery)

	query, ok := queryIn.(string)
	if !ok || query == "" {
		return false, fmt.Errorf("query not set")
	}

	query = strings.TrimSpace(query)
	queryParts := strings.Fields(query)

	if len(queryParts) == 0 {
		return false, fmt.Errorf("invalid query '%s'", query)
	}

	log.Debugf("query: '%s'", query)

	queryType := strings.ToLower(queryParts[0])

	switch queryType {
	case "select":
		paramsIn := context.GetInput(ivParams)
		var params map[string]interface{}

		if paramsIn != nil {
			params, ok = paramsIn.(map[string]interface{})
			if !ok {
				return false, fmt.Errorf("params is not valid: %v", params)
			}
		}

		ctIn := context.GetInput(ivColumnTypes)
		var columnTypes map[string]string
		if ctIn != nil {
			columnTypes, ok = ctIn.(map[string]string)
			if !ok {
				return false, fmt.Errorf("columnTypes is not valid: %v", params)
			}
		}

		results, err := DoSelect(dataSourceName, query, params, columnTypes)

		if err != nil {
			return false, err
		}

		context.SetOutput(ovResults, results)
	default:
		return false, fmt.Errorf("queryType '%s' not supported", queryType)
	}

	return true, nil
}

func DoSelect(dataSourceName, query string, params map[string]interface{},  columnTypes map[string]string) ([]map[string]interface{}, error) {

	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var rows *sqlx.Rows

	if len(params) > 0 {
		rows, err = db.NamedQuery(query, params)
	} else {
		rows, err = db.Queryx(query)
	}

	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		row := make(map[string]interface{})
		err := rows.MapScan(row)
		if err != nil {
			return nil, err
		}

		for k, encoded := range row {
			switch t := encoded.(type) {
			case []byte:

				strVal := string(t)
				dataType := data.TypeString

				if columnTypes != nil {
					if ct, ok := columnTypes[k]; ok {
						if dt, ok := data.ToTypeEnum(ct); ok {
							dataType = dt
						}
					}
				}

				row[k], err = data.CoerceToValue(strVal, dataType)
				if err != nil {
					return nil, err
				}
			}
		}

		results = append(results, row)
	}

	return results, nil
}

