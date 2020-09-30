# DataBackup
A Golang Command Line Application which creates a back up of the given folders and Sql databases.

## TABLE OF CONTENTS
- [PROBLEM STATEMENT](#problem-statement)
- [DOCUMENTATION](#documentation)
  *  [PROJECT STRUCTURE](#project-structure)
  *  [TECH STACK USED](#tech-stack-used)
  *  [JSON FILE FORMAT](#json-file-format)
  *  [RUN PROGRAM](#run-program)

## PROBLEM STATEMENT

To Back up list of given Databases using MySQL dump, and to back up folders to the given destination folder and create a zip file using Golang.

## DOCUMENTATION

### PROJECT STRUCTURE

``` bash
├── DataBackup
│   ├── Util 
│   │	├── CopyFolder.go
│   │	└── ZipFolder.go
│   ├── DataBackUp.go
│   ├── README.md
│   └── data.jsone
└── .gitignore
```
### TECH STACK USED 
* GoLang

### JSON FILE FORMAT

``` python 
{
		"userName": "root",
		"password": "Mydb@123",
		"hostName": "localhost",
		"port": "3306",
		"dumpDir": "D:/vicaraBackup",
		"zipFileName":"D:/vicaraBackup.zip",
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

Where,
* **userName**: UserName for MySQL.
* **password**: Password for MySQL.
* **hostName**: localhost or ip address of the connection.
* **port**: Port at which MySQL server is running at.
* **dumpDir**: Destination Directory/Folder to which the Databases and Folders needs to be transfered at.
* **zipFileName**: Name of the Backup file after compressing it to .zip format (Note: Do mention ".zip" extension to the name).
* **databases**: List of Databases to backup.
* **folders**: List of Folders to backup. 

### RUN PROGRAM
Run the program with the following command
**go run ./DataBackup.go**


