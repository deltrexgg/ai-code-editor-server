package helper

import "os"

const ProjectFolder = "project-files/"


func CreateFolder(folder_name string) error {
	location := ProjectFolder+folder_name
	err := os.MkdirAll(location, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func CreateFile(file_name string) error {
	location := ProjectFolder+file_name
	
	file, err := os.Create(location)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

func DeleteFile(file_name string) error {
	err := os.Remove(ProjectFolder+file_name)
	if err != nil {
		return err
	}

	return nil
}

func GetfilesNfolders(folder_loc string) ([]string, error) {
	var file_names []string

	files, err := os.ReadDir(ProjectFolder+folder_loc)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		file_names = append(file_names, file.Name())
	}

	return file_names, nil
}

func ReadFile(path string) (string, error) {

	loaction := ProjectFolder+path
	data, err := os.ReadFile(loaction)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func OverwriteFile(content string, path string) error {
	
	err := os.WriteFile(ProjectFolder+path, []byte(content),0644)
	if err != nil {
		return err
	}

	return nil
}