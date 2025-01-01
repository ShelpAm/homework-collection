package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
)

type TaskId = uuid.UUID
type Progress = float64

type FileUploader struct {
	fileProgresses map[TaskId]Progress
}

func MakeFileUploader() *FileUploader {
	return &FileUploader{
		fileProgresses: make(map[TaskId]Progress),
	}
}

func (fu *FileUploader) GetProgress(id string) (Progress, error) {
	taskId, err := uuid.Parse(id)
	if err != nil {
		return 0, err
	}

	progress, exists := fu.fileProgresses[taskId]
	if !exists {
		return 0, fmt.Errorf("file with id \"%s\" doesn't exist", id)
	}

	return progress, nil
}

func makeTaskId() TaskId {
	return TaskId(uuid.New())
}

// onFinish will be called iff the write process succeeded.
func (fu *FileUploader) ScheduleUploadTo(sr SizeableReader, dest string, onFinish func()) TaskId {
	taskId := makeTaskId()

	go func() {
		log.Println("Before")
		err := writeToFileWithProgress(sr, dest, func(p Progress) {
			fu.fileProgresses[taskId] = p
		})
		log.Println("After")

		if err != nil {
			delete(fu.fileProgresses, taskId) // Removes task from the map.
			log.Println("Error: " + err.Error())
			return
		}

		onFinish()
	}()

	log.Println("Out")

	return taskId
}
