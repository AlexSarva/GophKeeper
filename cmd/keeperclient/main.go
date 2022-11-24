package main

import (
	"AlexSarva/GophKeeper/models"
	"AlexSarva/GophKeeper/workclient"
	"io"
	"log"
	"os"

	"github.com/google/uuid"
)

func main() {
	cli, cliErr := workclient.WorkClient("http://localhost:8005/api/v1")
	if cliErr != nil {
		log.Fatalln(cliErr)
	}

	regData := &models.UserRegister{
		Username: "alexsarva",
		Email:    "alexsarva@yandex.ru",
		Password: "123456",
	}

	userReg, userRegErr := cli.Register(regData)
	if userRegErr != nil {
		log.Println(userRegErr)
	}
	log.Printf("%+v\n", userReg)

	loginData := &models.UserLogin{
		Email:    "alexsarva@gmail.com",
		Password: "12345",
	}
	userLog, userLogErr := cli.Login(loginData)
	if userLogErr != nil {
		log.Println(userLogErr)
	}
	log.Printf("%+v\n", userLog)

	userMe, userMeErr := cli.Me()
	if userMeErr != nil {
		log.Println(userMeErr)
	}
	log.Printf("%+v\n", userMe)

	info, infoErr := cli.ElementList("notes")
	if infoErr != nil {
		log.Println(infoErr)
	}
	log.Printf("%+v\n", info)

	elemId, _ := uuid.Parse("66cf84da-1ca4-4e15-9cd7-089254564115")
	elem, elemErr := cli.Element("cards", elemId)
	if elemErr != nil {
		log.Println(elemErr)
	}
	log.Printf("%+v\n %T", elem, elem)

	//tempCred := models.NewCred{
	//	Title:  "Тестовый набор 1",
	//	Login:  "alexsarva",
	//	Passwd: "77ofnWFF",
	//	Notes:  "Хрень какая-то",
	//}

	file, fileErr := os.Open("/Users/alexsarva/Documents/Бронь_СВО.xlsx")
	if fileErr != nil {
		log.Println(fileErr)
	}
	bodyBytes, readErr := io.ReadAll(file)
	if readErr != nil {
		log.Println(readErr)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	tempFile := models.NewFile{
		Title: "Тестовый файл",
		File:  bodyBytes,
		Notes: "",
	}

	result, resultErr := cli.AddElement("files", tempFile)
	if resultErr != nil {
		log.Println(resultErr)
	}
	log.Printf("%+v\n %T", result, result)

	e, eErr := cli.Element("notes", elemId)
	if eErr != nil {
		log.Println(eErr)
	}
	log.Printf("%+v\n %T", e, e)

	editNode := &models.NewNote{
		Title: "Третья заметка",
		Note:  "",
	}

	a, aErr := cli.EditElement("notes", editNode, elemId)
	if aErr != nil {
		log.Println(aErr)
	}

	log.Printf("%+v\n %T", a, a)

	d, dErr := cli.Delete("notes", elemId)
	if dErr != nil {
		log.Println(dErr)
	}

	log.Println(d)

}
