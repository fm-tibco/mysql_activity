
# MySQL 
## MySQL Database Activity 
This activity provides your flogo application execute database queries in MySQL.  Currently only select 
queries are supported.


# Third Party Libraries Used
https://github.com/go-sql-driver/mysql

https://github.com/jmoiron/sqlx

## Installation

```bash
flogo install github.com/fm-tibco/mysql_activity
```

## Metadata
Input and Output:

```json
{
  "input":[
    {
      "name": "dataSourceName",
      "type": "string",
      "required": true
    },
    {
      "name": "query",
      "type": "string",
      "required": true
    },
    {
      "name": "params",
      "type": "object"
    },
    {
      "name": "columnTypes",
      "type": "params"
    }
  ],
  "output": [
    {
      "name": "results",
      "type": "array"
    }
  ]
}
  ```
  
  ### Inputs
| Setting     | Description    |
|:------------|:---------------|
| dataSourceName | The db connection string |  
| query          | The db query statement |
| params         | Optional params for named db query |  
| columnTypes      Optional - can specify the types of the columns |  

##  Examples

### Query
```json
{
  "id": "dbquery",
  "name": "DbQuery",
  "activity": {
    "ref": "github.com/fm-tibco/mysql_activity",
    "input": {
      "dataSourceName": "username:password@tcp(host:port)/dbName",
      "query": "select * from test"
    }
  }
}
```

### Named Query

```json
{
  "id": "named_dbquery",
  "name": "Named DbQuery",
  "activity": {
    "ref": "github.com/fm-tibco/mysql_activity",
    "input": {
      "dataSourceName": "username:password@tcp(host:port)/dbName",
      "query": "select * from test where id > :id",
      "params": {
        "id":1
      }
    }
  }
}
```
### Issues
If type cannot be determined, value will default to string.  In order to fix this you can
specify an optional "columnTypes" to specify the type of the column

```json
{
  "id": "dbquery",
  "name": "DbQuery",
  "activity": {
    "ref": "github.com/fm-tibco/mysql_activity",
    "input": {
      "dataSourceName": "username:password@tcp(host:port)/dbName",
      "query": "select * from test",
      "columnTypes" : {
        "id":"integer"
      }
    }
  }
}
```
