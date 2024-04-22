package interrupt

import (
	"context"
	"fmt"
	"myapp/pkg/database"
	"myapp/pkg/words"
	"myapp/pkg/xkcd"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/schollz/progressbar/v3"
)

// Resources хранит ссылки на базу данных и клиент API
type Resources struct {
	ComicsDB *database.ComicsDatabase
	Client   *xkcd.ComicClient
}

// CreateDatabase инициализирует базу данных и клиента API, запускает воркеров для обработки комиксов
func CreateDatabase(baseURL string, dbPath string, numGoroutines int) (*database.ComicsDatabase, error) {
	client := xkcd.NewComicClient(baseURL)
	lastComicID := client.FetchLatestComicID()

	db, err := database.Open(dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// обрабатывает сигнал прерываня
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// добавлен прогрессбар,
	bar := progressbar.Default(int64(lastComicID - len(db.Records)))

	wg := &sync.WaitGroup{}
	taskChannel := make(chan int, lastComicID)

	// лог
	fmt.Println("\nComics download started")
	if len(db.Records) == 0 {
		fmt.Println("Creating a new database...")
	} else {
		fmt.Println("Continue loading from last saved comic...")
	}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range taskChannel {
				select {
				case <-ctx.Done(): // Когда контекст отменяется  через Ctrl+C
					return // Горутина завершает работу
				default:
					comic := client.FetchComic(id)
					if comic.ID != 0 {
						comicData := convertComic(comic)
						db.AddComic(id, comicData)
					}
					bar.Add(1)
				}
			}
		}()
	}

	enqueueTasks(ctx, lastComicID, db, taskChannel)

	wg.Wait()
	cancel() // Останавливаем все горутины, если процесс завершен
	bar.Finish()
	fmt.Println("Finish. database.json created.")
	return db, nil
}

// Отправляет задачи в канал, предварительно проверяя наличие комикса в базе данных
func enqueueTasks(ctx context.Context, numTasks int, db *database.ComicsDatabase, taskChannel chan<- int) {
	for id := 1; id <= numTasks; id++ {
		if _, exists := db.Records[id]; !exists {
			select {
			case <-ctx.Done():
				close(taskChannel)
				return
			case taskChannel <- id:
			}
		}
	}
	close(taskChannel)
}

// Конвертирует данные комикса из формата API в формат базы данных
func convertComic(comic xkcd.Comics) database.ComicData {
	text := strings.Fields(comic.Transcript + " " + comic.AltText)
	normalText, err := words.Normalization(text)
	if err != nil {
		fmt.Println("Error normalizing words:", err)
		return database.ComicData{}
	}
	return database.ComicData{
		Img:      comic.ImageURL,
		Keywords: normalText,
	}
}
