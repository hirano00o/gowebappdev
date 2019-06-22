package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

// ErrNoAvatarURL is error
var ErrNoAvatarURL = errors.New("chat: failed to get avatar's url")

// AuthAvatar is ...
type AuthAvatar struct{}

// UseAuthAvatar is ...
var UseAuthAvatar AuthAvatar

// GravatarAvatar is ...
type GravatarAvatar struct{}

// UseGravatar is ...
var UseGravatar GravatarAvatar

// FileSystemAvatar is ...
type FileSystemAvatar struct{}

// UseFileSystemAvatar is ...
var UseFileSystemAvatar FileSystemAvatar

// Avatar represents interface profile image
type Avatar interface {
	GetAvatarURL(u ChatUser) (string, error)
}

// TryAvatars is all avatars
type TryAvatars []Avatar

// GetAvatarURL return avatar url
func (a AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

// GetAvatarURL return avatar url
func (g GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

// GetAvatarURL return avatar path
func (f FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}
	return "", ErrNoAvatarURL
}

// GetAvatarURL try all avatar
func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}
