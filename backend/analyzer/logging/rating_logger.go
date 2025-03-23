package logging

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type RatingLogger struct {
	mu   sync.Mutex
	file *os.File
}

func NewRatingLogger() (*RatingLogger, error) {
	file, err := os.OpenFile("unknown_ratings.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error creando archivo de log: %v", err)
	}

	return &RatingLogger{
		file: file,
	}, nil
}

func (l *RatingLogger) LogUnknownRating(rating string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	logEntry := fmt.Sprintf("[%s] Rating no encontrado: %q\n",
		time.Now().Format("2006-01-02 15:04:05"),
		rating)

	_, err := l.file.WriteString(logEntry)
	if err != nil {
		fmt.Printf("Error al escribir en el archivo de log: %v\n", err)
	}
}

func (l *RatingLogger) Close() error {
	return l.file.Close()
}
