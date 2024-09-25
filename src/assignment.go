package main

import (
	"crypto/sha512"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
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

type Student struct {
	Name     string
	SchoolId string
}

func (s *Student) Submit(a *Assignment, f multipart.File, filename string) error {
	err := a.Receive(*s, f, filename)
	if err != nil {
		return err
	}

	return nil
}

type Assignment struct {
	Name      string
	BeginTime time.Time
	EndTime   time.Time
}

func (a *Assignment) Path() string {
	return filepath.Join("homeworks", a.Name)
}

func (a *Assignment) Receive(s Student, f multipart.File, filename string) error {
	if now := time.Now(); now.Before(a.BeginTime) || now.After(a.EndTime) {
		return errors.New("Submission time out of bound (作业提交超出时限)")
	}

	baseDir := filepath.Join(a.Path(), s.SchoolId+s.Name)
	err := os.RemoveAll(baseDir) // Overrides origin file/dir.
	if err != nil {
		return err
	}
	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		return err
	}

	savePath := filepath.Join(baseDir, filename)
	file, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, f)
	if err != nil {
		return errors.New("Cannot copy file")
	}

	return nil
}
