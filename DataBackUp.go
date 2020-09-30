package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"vicara/Util"

	"github.com/JamesStewy/go-mysqldump"
	_ "github.com/go-sql-driver/mysql"
)

/*
	{
		"username": "root",
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
*/

/*
	Define Structure for the given Json file.
*/
type Data struct {
	UserName    string     `json:"userName"`
	Password    string     `json:"password"`
	HostName    string     `json:"hostName"`
	Port        string     `json:"port"`
	DumpDir     string     `json:"dumpDir"`
	ZipFileName string     `json:"zipFileName"`
	Databases   []Database `json:"databases"`
	Folders     []Folder   `json:"folders"`
}

type Database struct {
	DatabaseName string `json:"database"`
}

type Folder struct {
	FolderName string `json:"folder"`
}

/*
	Param *Data: Data object
	returns: Success on performing mysql Dump operation to the given directory.
*/
func dumpDatabase(data *Data) (string, error) {
	userName := data.UserName
	password := data.Password
	hostName := data.HostName
	port := data.Port
	dumpDir := data.DumpDir
	numberOfDatabases := len(data.Databases)

	for i := 0; i < numberOfDatabases; i++ {
		dbName := data.Databases[i].DatabaseName
		dumpFilenameFormat := fmt.Sprintf("%s-backup", dbName)

		// establish connection to the database
		// connection format => username:password@tcp(hostname:port)/databaseName
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, password, hostName, port, dbName)
		db, err := sql.Open("mysql", connectionString)
		if err != nil {
			// fmt.Println("Error opening database: ", err)
			return "error opening database", err
		}

		// Register database with mysqldump
		dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
		if err != nil {
			// fmt.Println("Error registering databse:", err)
			return "Error registering databse", err
		}

		// Dump database to file
		resultFilename, err := dumper.Dump()
		if err != nil {
			// fmt.Println("Error dumping:", err)
			return "Error dumping", err
		}

		fmt.Printf("Database is saved to %s \n", resultFilename)

		dumper.Close()
	}

	return "Success!", nil
}

/*
	Param1: Json file name.
	returns: Data object after parsing the given json.
*/
func getJsonObject(fileName string) *Data {
	jsonFile, err := os.Open(fileName)
	// check if error occured
	if err != nil {
		fmt.Println(err)
	}
	// close the json connection after program executation
	defer jsonFile.Close()

	// convert the json file in memory to byte array
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data Data
	json.Unmarshal(byteValue, &data)

	return &data
}

/*
	transfers the folder from source to destination. Uses CopyDir(source, destination) function from Util package.
	Param: *Data. can get the following information from the parameter
		1) Destination Folder to store the back up at.
		2) List of folders to backup.
*/
func copyFolder(data *Data) (string, error) {
	dumpDir := data.DumpDir
	totalFolder := len(data.Folders)

	if totalFolder == 0 {
		return "No folders to move", nil
	}
	for i := 0; i < totalFolder; i++ {
		source := data.Folders[i].FolderName
		destination := dumpDir + sliceString(source)

		err := Util.CopyDir(source, destination)
		if err != nil {
			return "", err
		}

	}

	return "Folders moved successfully", nil
}

//HELPER FUNCTIONS

/*
	print contents of Data object
*/
func printDetails(data *Data) {
	fmt.Println("user name : ", data.UserName)
	fmt.Println("passord : ", data.Password)
	fmt.Println("HostName : ", data.HostName)
	fmt.Println("port : ", data.Port)
	fmt.Println("destination dir", data.DumpDir)
	fmt.Println("zip file name :", data.ZipFileName)
	fmt.Println("database Name ", data.Databases[0].DatabaseName)
}

/*
	function used in copyFolder function to create destination folder name.
	Param: filename
	returns: String ==> removes the Drive name and semicolon
	example: given filename D:/dirName
	Output: /dirName
*/
func sliceString(fileName string) string {
	returnString := fileName[2:]
	return returnString
}

func main() {

	// get file name as input
	var jsonFileName string
	fmt.Println("enter json file name")
	fmt.Scanln(&jsonFileName)

	// get data object
	data := getJsonObject(jsonFileName)
	status, err := dumpDatabase(data)
	if err != nil {
		fmt.Printf("%s %v \n", status, err)
	}
	fmt.Printf("Database dump status: %s\n", status)

	// copy folders to the destination folder
	copyStatus, er := copyFolder(data)
	if er != nil {
		fmt.Println(er)
	}
	fmt.Printf("Folders status: %s\n", copyStatus)

	// create zip file of the backup folder.
	// targetFile := "D:/dump.zip"
	targetFile := data.ZipFileName
	sourceFile := data.DumpDir
	zipError := Util.Zipit(sourceFile, targetFile)

	if zipError != nil {
		fmt.Println(zipError)
	} else {
		fmt.Printf("Database and Folders zipped at %s", targetFile)
	}

}
