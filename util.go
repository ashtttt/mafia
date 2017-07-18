package main

import (
	"encoding/base32"
	"errors"
	"os"

	"github.com/ashtttt/mafia/topt"
)

func getSecret(profile string) (string, error) {
	secret := os.Getenv("TOPT_SECRET")
	if len(profile) > 0 {
		val, err := viewDB(profile)
		if err != nil {
			return "", err
		}
		if len(val) > 0 {
			secret = val
		} else {
			return "", errors.New("Profile name NOT Found! Please check")
		}
	} else if len(secret) <= 0 {
		return "", errors.New("TOPT_SECRET environment variable is not set. Please set it or pass a previously configured profile")
	}
	return secret, nil
}

func getOpt(profile string) (string, error) {
	secret, err := getSecret(profile)
	if err != nil {
		return "", err
	}
	sec, _ := base32.StdEncoding.DecodeString(secret)
	opt := topt.NewTOPT()
	return opt.TokenCode(sec), nil
}

func home() string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = os.Getenv("USERPROFILE")
	}
	return homeDir
}
