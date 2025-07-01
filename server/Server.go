package server

import (
	"log"
	"notes_project/database"
	"notes_project/envs"
)

func InitServer() {
	errEnvs := envs.LoadEnvs()
	if errEnvs != nil {
		log.Fatal("Ошибка инициализации ENV: ", errEnvs)
	} else {
		log.Println("Инициализация ENV прошла успешно")
	}

	//Инициализация БД
	errDB := database.InitDatabase()
	if errDB != nil {
		log.Fatal("Ошибка подключения к базе данных: ", errDB)
	} else {
		log.Println("Успешное подключение к базе данных")
	}
	errRedis := database.InitRedis()
	if errRedis != nil {
		log.Fatal("Ошибка подключения к Redis: ", errRedis)
	} else {
		log.Println("Успешное подключение к Redis")
	}
}

func StartServer() {
	InitRotes()
}
