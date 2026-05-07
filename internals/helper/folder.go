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

func DeleteFile(file_name string) error {
	err := os.Remove("project-files/"+file_name)
	if err != nil {
		return err
	}

	return nil
}

func GetfilesNfolders(folder_loc string) ([]string, error) {
	var file_names []string

	files, err := os.ReadDir("project-files/"+folder_loc)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		file_names = append(file_names, file.Name())
	}

	return file_names, nil
}