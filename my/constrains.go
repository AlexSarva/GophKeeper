package main

import (
	"AlexSarva/GophKeeper/utils"
	"log"
)

func main() {
	a := "124weqwe"
	e := utils.CheckOnlyDigits(a)
	log.Println(e)
}
