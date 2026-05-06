package helper

import "os"


func CreateFolder(folder_name string) error {
	location := "project-files/"+folder_name
	err := os.MkdirAll(location, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func CreateFile(file_name string) error {
	location := "project-files/"+file_name
	
	file, err := os.Create(location)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}