package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
)

type TaskId = uuid.UUID
type Progress = float64

type FileUploader struct {
	fileProgresses map[TaskId]Progress
	mutex          sync.Mutex
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
func (fu *FileUploader) ScheduleUploadTo(sr SizeableReader, dest string, onFinish func(TaskId)) TaskId {
	taskId := makeTaskId()

	go func() {
		// log.Println("Before")
		err := writeToFileWithProgress(sr, dest, func(p Progress) {
			fu.mutex.Lock()
			defer fu.mutex.Unlock()
			fu.fileProgresses[taskId] = p
		})
		// log.Println("After")

		if err != nil {
			fu.mutex.Lock()
			defer fu.mutex.Unlock()
			delete(fu.fileProgresses, taskId) // Removes task from the map.
			log.Println("Cannot write to file " + dest + ": " + err.Error())
			return
		}

		onFinish(taskId)
	}()

	// log.Println("Out")

	return taskId
}
