package reptile

import (
	"fmt"
	"os"
)

// create folder for download path
func CreateFolder(folderpath string) {

	fmt.Println("check folder path...", folderpath)

	_dir := folderpath
	exist, err := PathExists(_dir)
	if err != nil {
		fmt.Printf("get dir error![%v]\n", err)
		return
	}

	if exist {
		fmt.Printf("folder is exist:", _dir)
	} else {
		fmt.Printf("create folder:", _dir)
		// create folder
		err := os.Mkdir(_dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
		} else {
			fmt.Printf("mkdir success!\n")
		}
	}
}

// check path exist
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
