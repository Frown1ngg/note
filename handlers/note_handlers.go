package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"notes_project/database"
	"notes_project/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func GetNoteHandler(ctx *gin.Context) {
	//ctx.JSON(http.StatusOK, "GetNoteHandler")
	authorId := 1
	id := ctx.Param("id")

	collection := database.MongoClient.Database("admin").Collection(fmt.Sprintf("notes/%d", authorId))

	var note models.Note
	filter := bson.M{"id": id}

	errFind := collection.FindOne(ctx, filter).Decode(&note)
	if errFind != nil {
		ctx.JSON(http.StatusNotFound, "Заметка не найдена")
	}
	ctx.JSON(http.StatusOK, &note)

}

func GetNotesHandler(ctx *gin.Context) {
	// ctx.JSON(http.StatusOK, "GetNotesHandler")
	authorId := 1
	var notes []models.Note

	val, err := database.RedisClient.Get(fmt.Sprintf("notes/%d", authorId)).Result()
	if err == redis.Nil {
		log.Printf("Кеш не найден, загружаем из БД")
		// Получаем коллекцию "notes"
		collection := database.MongoClient.Database("admin").Collection(fmt.Sprintf("notes/%d", authorId))

		// Поиск документов без фильтров для получения всех заметок
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Закрытие курсора, при завершении работы функции
		defer cursor.Close(ctx)
		// Итерация по курсору и декодирование документов в заметки
		for cursor.Next(ctx) {
			var note models.Note
			err := cursor.Decode(&note)
			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			notes = append(notes, note)
		}
		// Проверка на ошибки после итерации
		if err := cursor.Err(); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Проверка на наличие заметок
		if len(notes) == 0 {
			ctx.JSON(http.StatusOK, "Заметок не найдено")
		} else {
			notesJSON, err := json.Marshal(notes)
			if err != nil {
				log.Printf("Ошибка при сериализации заметок: %v", err)
			} else {
				err := database.RedisClient.Set(fmt.Sprintf("notes/%d", authorId), string(notesJSON), 1440*time.Minute).Err()
				if err != nil {
					log.Printf("Ошибка при сохранении в Redis: %v", err)
				}
			}
			ctx.JSON(http.StatusOK, notes)
		}
	} else {
		log.Printf("Кеш найден, загружаем из кеша")
		notes := make([]models.Note, 0)
		json.Unmarshal([]byte(val), &notes)
		ctx.JSON(http.StatusOK, notes)
	}

}
func resetCache(val string) {
	database.RedisClient.Del(val)
}

func recordCasheToRedis(notes []models.Note, authorId int) {
	notesJSON, err := json.Marshal(notes)

	if err != nil {
		log.Printf("Ошибка при сериализации заметок: %v", err)
	} else {
		err := database.RedisClient.Set(fmt.Sprintf("notes/%d", authorId), string(notesJSON), 1440*time.Minute).Err()
		if err != nil {
			log.Printf("Ошибка при сохранении в Redis: %v", err)
		}
	}
}
func getFromCache(val string, ctx *gin.Context) {
	log.Printf("Кеш найден, загружаем из Кеша")
	notes := make([]models.Note, 0)
	json.Unmarshal([]byte(val), &notes)
	ctx.JSON(http.StatusOK, notes)
}
func DeleteNoteHandler(ctx *gin.Context) {
	//ctx.JSON(http.StatusOK, "DeleteNoteHandler")
	var authorID = 1
	id := ctx.Param("id")

	collection := database.MongoClient.Database("admin").Collection(fmt.Sprintf("notes/%d", 1))

	filter := bson.M{"id": id}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})

	}
	if result.DeletedCount == 0 {
		ctx.JSON(http.StatusNotFound, "Заметка не найдена")
	} else {
		resetCache(fmt.Sprintf("notes/%d", authorID))
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
		ctx.JSON(http.StatusNotFound, "Заметка не найдена")
	} else {
		resetCache(fmt.Sprintf("notes/%d", note.AuthorID))
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
	resetCache(fmt.Sprintf("notes/%d", note.AuthorID))
	ctx.JSON(http.StatusOK, gin.H{
		"note":    note,
		"message": "Заметка успешно создана"})
}
