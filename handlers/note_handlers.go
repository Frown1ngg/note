package handlers

import (
	"fmt"
	"net/http"
	"notes_project/database"
	"notes_project/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func GetNoteHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "GetNoteHandler")
}

func GetNotesHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "GetNotesHandler")
}

func DeleteNoteHandler(ctx *gin.Context) {
	//ctx.JSON(http.StatusOK, "DeleteNoteHandler")
	id := ctx.Param("id")

	collection := database.MongoClient.Database("admin").Collection(fmt.Sprintf("notes/%d", 1))

	filter := bson.M{"id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})

	}
	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusOK, "Заметка не найдена")
	} else {
		ctx.JSON(http.StatusNotFound, "Заметка успешно удалена")
	}
}

func UpdateNoteHandler(ctx *gin.Context) {
	//ctx.JSON(http.StatusOK, "UpdateNoteHandler")
	authorId := 1
	id := ctx.Param("id")

	var note models.Note
	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверные данные"})
		return
	}

	collection := database.MongoClient.Database("admin").Collection(fmt.Sprintf("notes/%d", authorId))
	updateFields := bson.M{}
	if note.Name != nil {
		updateFields["name"] = note.Name
	}
	if note.Content != nil {
		updateFields["content"] = note.Content
	}

	update := bson.M{"$set": updateFields}

	filter := bson.M{"id": id}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusOK, "Заметка не найдена")
	} else {
		ctx.JSON(http.StatusOK, "Заметка успешно обновлена")
	}
}

func CreateNoteHandler(ctx *gin.Context) {
	// ctx.JSON(http.StatusOK, "CreateNoteHandler")
	var note models.Note

	if err := ctx.ShouldBindJSON(&note); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}
	// Получить уникальный id
	note.Id = uuid.New().String()
	// Тестовый ID автора
	note.AuthorID = 1
	// Получаем коллекцию "notes"
	Collection := database.MongoClient.Database("admin").Collection(fmt.Sprintf("notes/%d", note.AuthorID))
	// Вставляем заметку в коллекцию
	_, errInsert := Collection.InsertOne(ctx, note)
	if errInsert != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": errInsert.Error()})
	}
	// Если ошибок нет, то возвращаем заметку и статус 200
	ctx.JSON(http.StatusOK, gin.H{
		"note":    note,
		"message": "Заметка успешно создана"})
}
