package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/JamesStewy/go-mysqldump"
	_ "github.com/go-sql-driver/mysql"
)

// {
//     "username": "root",
//     "password": "Mydb@123",
//     "hostName":"localhost",
//     "port":"3306",
//     "dumpDir":"",
//     "databases": [
//         {
//             "database" : ""
//         },
//         {
//             "database":""
//         }
//     ],

//     "folders": [
//         {
//             "folder":""
//         },
//         {
//             "folder":""
//         }
//     ]
// }

type Data struct {
	UserName  string     `json:"userName"`
	Password  string     `json:"password"`
	HostName  string     `json:"hostName"`
	Port      string     `json:"port"`
	DumpDir   string     `json:"dumpDir"`
	Databases []Database `json:"databases"`
	Folders   []Folder   `json:"folders"`
}

type Database struct {
	DatabaseName string `json:"database"`
}

type Folder struct {
	FolderName string `json:"folder"`
}

// print details in the Data object
func printDetails(data *Data) {
	fmt.Println("user name : ", data.UserName)
	fmt.Println("passord : ", data.Password)
	fmt.Println("HostName : ", data.HostName)
	fmt.Println("port : ", data.Port)
	fmt.Println("destination dir", data.DumpDir)
	fmt.Println("database Name ", data.Databases[0].DatabaseName)
}

// returns "success" on successfull dumping of database to the the destination directory.
func dumpDatabase(data *Data) (string, error) {
	userName := data.UserName
	password := data.Password
	hostName := data.HostName
	port := data.Port
	dumpDir := data.DumpDir
	numberOfDatabases := len(data.Databases)

	for i := 0; i < numberOfDatabases; i++ {
		dbName := data.Databases[i].DatabaseName
		dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", dbName)

		// name of file with which it will be saved at the the destination dir.
		fmt.Printf("dump dile name %s \n", dumpFilenameFormat)

		// establish connection to the database
		// connection format => username:password@tcp(hostname:port)/databaseName
		connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userName, password, hostName, port, dbName)
		db, err := sql.Open("mysql", connectionString)
		if err != nil {
			fmt.Println("Error opening database: ", err)
			return "error opening database", err
		}

		// Register database with mysqldump
		dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
		if err != nil {
			fmt.Println("Error registering databse:", err)
			return "Error registering databse", err
		}

		// Dump database to file
		resultFilename, err := dumper.Dump()
		if err != nil {
			fmt.Println("Error dumping:", err)
			return "Error dumping", err
		}

		fmt.Printf("File is saved to %s \n", resultFilename)

		dumper.Close()
	}

	return "Success!", nil
}

// parse json object and convert it to #Data object
func getJsonObject(fileName string) *Data {
	jsonFile, err := os.Open(fileName)
	// check if error occured
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("successfully loaded the file")

	// close the json connection after program executation
	defer jsonFile.Close()

	// convert the json file in memory to byte array
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data Data
	json.Unmarshal(byteValue, &data)

	return &data
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func copyFolder(data *Data) (string, error) {
	dumpDir := data.DumpDir
	totalFolder := len(data.Folders)

	if totalFolder == 0 {
		return "No folders to move", nil
	}
	for i := 0; i < totalFolder; i++ {
		source := data.Folders[i].FolderName
		destination := dumpDir + sliceString(source)

		err := CopyDir(source, destination)
		if err != nil {
			return "", err
		}

	}

	return "Folders moved successfully", nil
}

func sliceString(fileName string) string {
	returnString := fileName[2:]
	return returnString
}

func main() {
	var jsonFileName string
	fmt.Println("enter json file name")
	fmt.Scanln(&jsonFileName)
	data := getJsonObject(jsonFileName)

	status, err := dumpDatabase(data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Database dump status: %s\n", status)

	copyStatus, er := copyFolder(data)
	if er != nil {
		fmt.Println(er)
	}
	fmt.Printf("Folders status: %s\n", copyStatus)

}
