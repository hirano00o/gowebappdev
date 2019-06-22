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
	GetAvatarURL(c *client) (string, error)
}

// GetAvatarURL return avatar url
func (a AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// GetAvatarURL return avatar url
func (g GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

// GetAvatarURL return avatar path
func (f FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			if files, err := ioutil.ReadDir("avatars"); err == nil {
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					if match, _ := filepath.Match(useridStr+"*", file.Name()); match {
						return "/avatars/" + file.Name(), nil
					}
				}
			}
		}
	}
	return "", ErrNoAvatarURL
}
