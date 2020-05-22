package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/s3rj1k/jrpc2/client"
)

/*
curl -X POST https://internalrpc.mirohost.net/v1/  -u user:password  -H "Content-Type:application/json" -d '{
    "jsonrpc":"2.0",
    "method": "getContainersList",
    "params": {
        "accessKey":"rpcAccessKey"
    }
}'
*/

// GetListErrInnerDataObj - JSON-RPC error inner data object, JSON-RPC server specific.
type GetListErrInnerDataObj struct {
	Message []string `json:"message"`
}

// GetListResultObj - JSON-RPC GetContainersList result object.
type GetListResultObj []struct {
	CTID          string `json:"ctid"`
	FQDN          string `json:"fqdn"`
	Node          string `json:"node"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	Backup        string `json:"backup"`
	CreateBackups bool   `json:"createBackups"`
}

// GetListParamsObj - JSON-RPC GetContainersList object.
type GetListParamsObj struct {
	AccessKey string `json:"accessKey"`
}

// GetContainersListMethod - calls remote JSON-RPC server to get containers list.
func GetContainersListMethod(c *client.Config, key string) (GetListResultObj, error) {
	// set JSON-RPC method name
	const methodName = "getContainersList"

	// prepare results object
	var resultObj GetListResultObj

	// JSON-RPC params field
	paramsObj := GetListParamsObj{
		AccessKey: key,
	}

	// convert params object to bytes
	paramsData, err := json.Marshal(paramsObj)
	if err != nil {
		return resultObj, fmt.Errorf("method=%s params error: %s", methodName, err.Error())
	}

	// sent JSON-RPC request and get back response
	resultRawData, err := c.Call(methodName, paramsData)
	if err != nil {
		// error data object
		var errObj GetListErrInnerDataObj

		// decode inner (JSON-RPC) error message
		if errInnerDataObjError := json.Unmarshal(resultRawData, &errObj); errInnerDataObjError == nil {
			return resultObj, fmt.Errorf("%s: %s", err.Error(), strings.Join(errObj.Message, ", "))
		}

		// return plain error
		return resultObj, err
	}

	// decode result
	err = json.Unmarshal(resultRawData, &resultObj)
	if err != nil {
		return resultObj, fmt.Errorf("method=%s result error: %s", methodName, err.Error())
	}

	return resultObj, nil
}

// GetServiceDevicesListMethod - calls remote JSON-RPC server to get service nodes and containers list.
func GetServiceDevicesListMethod(c *client.Config, key string) (GetListResultObj, error) {
	// set JSON-RPC method name
	const methodName = "getServiceDevicesList"

	// prepare results object
	var resultObj GetListResultObj

	// JSON-RPC params field
	paramsObj := GetListParamsObj{
		AccessKey: key,
	}

	// convert params object to bytes
	paramsData, err := json.Marshal(paramsObj)
	if err != nil {
		return resultObj, fmt.Errorf("method=%s params error: %s", methodName, err.Error())
	}

	// sent JSON-RPC request and get back response
	resultRawData, err := c.Call(methodName, paramsData)
	if err != nil {
		// error data object
		var errObj GetListErrInnerDataObj

		// decode inner (JSON-RPC) error message
		if errInnerDataObjError := json.Unmarshal(resultRawData, &errObj); errInnerDataObjError == nil {
			return resultObj, fmt.Errorf("%s: %s", err.Error(), strings.Join(errObj.Message, ", "))
		}

		// return plain error
		return resultObj, err
	}

	// decode result
	err = json.Unmarshal(resultRawData, &resultObj)
	if err != nil {
		return resultObj, fmt.Errorf("method=%s result error: %s", methodName, err.Error())
	}

	return resultObj, nil
}

// CleanExpected deletes old expected records.
func CleanExpected(db *sql.DB) error {
	if _, err := db.Exec(`
		DELETE FROM
			expected
		WHERE
			date = NOW()::date;
		`,
	); err != nil {
		Error.Printf("Database DELETE error: %v\n", err)

		return err
	}

	return nil
}

// InsertExpected inserts expected backups, that did not start today.
func (obj GetListResultObj) InsertExpected(db *sql.DB) error {
	// insert records, that absent in compareTable
	for _, el := range obj {
		if el.CreateBackups {
			if _, err := db.Exec(`
				INSERT INTO
					expected (date, source, destination)
				SELECT
					NOW()::date,
					$1,
					$2
				WHERE NOT EXISTS
					(
						SELECT id FROM
							events
						WHERE
							START::date = NOW()::date
						AND
							source_container = $3
						AND
							destination = $4
					)
				AND NOT EXISTS
					(
						SELECT id FROM
							events
						WHERE
							START::date = NOW()::date
						AND
							source_host = $5
						AND
							destination = $6
					);
				`,
				el.FQDN,
				el.Backup,
				el.FQDN,
				el.Backup,
				el.FQDN,
				el.Backup,
			); err != nil {
				Error.Printf("Database INSERT error: %v\n", err)

				return err
			}
		}
	}

	return nil
}
