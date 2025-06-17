package server

import (
	"notes_project/handlers"

	"github.com/gin-gonic/gin"
)

func InitRotes() {
	router := gin.Default()                                // Инициализация роута (по умолчанию)
	router.PUT("/note", handlers.CreateNoteHandler)        // Редактирование заметки
	router.DELETE("/note/:id", handlers.DeleteNoteHandler) // Удаление заметки по id
	router.GET("/note/:id", handlers.GetNoteHandler)       // Получение заметки по id
	router.POST("/note/:id", handlers.UpdateNoteHandler)   // Создание заметки по id
	router.GET("/notes", handlers.GetNotesHandler)         // Получение списка всех заметок
	router.Run(":8080")                                    // Запуск сервера по порту в нашем случаи 8080
}
