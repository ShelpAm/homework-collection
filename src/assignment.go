package main

import (
	"crypto/sha512"
	"errors"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type User struct {
	Username       string
	HashedPassword [64]byte
}

func (u *User) Login(password string) error {
	if hashedPassword := sha512.Sum512([]byte(password)); hashedPassword != u.HashedPassword {
		return errors.New("Password is incorrect")
	}

	return nil
}

type Account struct {
	Name     string
	SchoolId string
}

type Student struct {
	Account      Account
	OnSubmitting sync.Mutex
}

func (s *Student) Submit(a *Assignment, r SizeableReader, filename string, onFinish func(TaskId)) (TaskId, error) {
	// Uploading policy for the same user:
	// 1. Queuing
	s.OnSubmitting.Lock()
	// 2. Aborting
	// if !s.OnSubmitting.TryLock() {
	// 	return TaskId{}, errors.New("This student is on submitting another homework. Please cancel that or wait for it to finish.")
	// }

	taskId, err := a.Receive(s, r, filename, func(id TaskId) {
		s.OnSubmitting.Unlock()
		onFinish(id)
	})
	if err != nil {
		return taskId, err
	}

	return taskId, nil
}

type Assignment struct {
	Name      string
	BeginTime time.Time
	EndTime   time.Time
}

func (a *Assignment) Path() string {
	return filepath.Join("homeworks", a.Name)
}

func (a *Assignment) Receive(s *Student, r SizeableReader, filename string, onFinish func(TaskId)) (TaskId, error) {
	if now := time.Now(); now.Before(a.BeginTime) || now.After(a.EndTime) {
		return TaskId{}, errors.New("Submission too late (作业提交超出时限)")
	}

	baseDir := filepath.Join(dataDir, a.Path(), s.Account.SchoolId+s.Account.Name)

	// Overrides origin file/dir.
	if err := os.RemoveAll(baseDir); err != nil {
		log.Fatalln("Failed to remove directory " + baseDir + ": " + err.Error())
	}

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		log.Fatalln("Failed to make directory " + baseDir + ": " + err.Error())
	}

	savePath := filepath.Join(baseDir, filename)

	id := fileUploader.ScheduleUploadTo(r, savePath, onFinish)
	log.Println("Receiving assignment " + a.Name + " " + filename + " from " + s.Account.SchoolId + " " + s.Account.Name + ", task id: " + id.String())
	return id, nil
	// err = writeToFileWithProgress(f, savePath, func(progress float64) {
	// fmt.Println("Progress: ", progress)
	// })
	// if err != nil {
	// 	return errors.New("Cannot copy file")
	// }

	// file, err := os.Create(savePath)
	// if err != nil {
	// 	return err
	// }
	// defer file.Close()
	//
	// _, err = io.Copy(file, f)
	// if err != nil {
	// 	return errors.New("Cannot copy file")
	// }

	// return nil
}
