# DataBackup
A golang project which creates a back up of the given  folders and Sql databases .

## PROBLEM STATEMENT

To Back up list of given Databases using MySQL dump, and to back up folders to the given destination folder and create a zip file using Golang.

## DOCUMENTATION

### TECH STACK USED 
* GoLang

### JSON FILE FORMAT

``` python 
{
		"username": "root",
		"password": "Mydb@123",
		"hostName": "localhost",
		"port": "3306",
		"dumpDir": "D:/vicaraBackup",
		"zipFileName":"D:/vicaraBackup.zip,
		"databases": [
			{
				"database" : "DB1"
			},
			{
				"database":"DB2"
			}
		],

		"folders": [
			{
				"folder":"FOLDER1"
			},
			{
				"folder":"FOLDER2"
			}
		]
	}

```

