package main

import (
	"encoding/base32"
	"errors"
	"flag"
	"os/signal"
	"strings"
	"syscall"

	"os"

	"path/filepath"

	"time"

	"github.com/ashtttt/mafia/topt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
	"github.com/mitchellh/colorstring"
)

type KeyCommand struct {
	serial string
	creds  *sts.Credentials
}

var message = "Credentials got updated. Valid for 12 hours, will be auto roated!"

func (k *KeyCommand) Run(args []string) error {

	flags := flag.NewFlagSet("keys", flag.ContinueOnError)
	flags.Usage = func() {
		colorstring.Println("[red]" + k.commandHelp())
	}
	profile := flags.String("profile", "", "")

	flags.Parse(args[1:])
	args = flags.Args()

	if len(*profile) <= 0 {
		flags.Usage()
		return errors.New("Profile name is required")
	}
	if len(args) > 0 {
		flags.Usage()
	}
	os.Setenv("AWS_PROFILE", *profile)
	err := k.getSerialNumber()
	if err != nil {
		return err
	}

	sigs := make(chan os.Signal, 1)
	errCh := make(chan error)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			err = k.getSessionToken()
			if err != nil {
				errCh <- err
			}
			err = k.updateCreds()
			if err != nil {
				errCh <- err
			}
			colorstring.Printf("[green] %s \n", message)
			time.Sleep(11 * time.Hour)
		}

	}()

	for {
		select {
		case done := <-errCh:
			return done
		case <-sigs:
			return errors.New("process interrupted. Exiting")
		}

	}
}

func (k *KeyCommand) getSerialNumber() error {
	sess := session.Must(session.NewSession())
	svc := sts.New(sess)
	var params *sts.GetCallerIdentityInput
	resp, err := svc.GetCallerIdentity(params)

	if err != nil {
		return err
	}
	k.serial = strings.Replace(*resp.Arn, "user", "mfa", 1)
	return nil
}

func (k *KeyCommand) getSessionToken() error {
	sess := session.Must(session.NewSession())
	svc := sts.New(sess)

	secret := os.Getenv("TOPT_SECRET")
	if len(secret) <= 0 {
		return errors.New("TOPT_SECRET environment variable is not set. Please do so")
	}
	sec, _ := base32.StdEncoding.DecodeString(secret)
	opt := topt.NewTOPT()
	code := opt.TokenCode(sec)

	params := &sts.GetSessionTokenInput{
		SerialNumber: aws.String(k.serial),
		TokenCode:    aws.String(code),
	}
	resp, err := svc.GetSessionToken(params)

	if err != nil {
		return err
	}
	k.creds = resp.Credentials
	return nil
}

func credfile() string {
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = os.Getenv("USERPROFILE")
	}
	return filepath.Join(homeDir, ".aws", "credentials")
}

func (k *KeyCommand) updateCreds() error {
	filename := credfile()
	config, err := ini.Load(filename)
	if err != nil {
		return err
	}

	config.DeleteSection("default")
	config.NewSection("default")
	config.Section("default").NewKey("aws_access_key_id", *k.creds.AccessKeyId)
	config.Section("default").NewKey("aws_secret_access_key", *k.creds.SecretAccessKey)
	config.Section("default").NewKey("aws_session_token", *k.creds.SessionToken)
	config.SaveTo(filename)

	return nil

}

func (k *KeyCommand) commandHelp() string {
	var usage = `Usage: mafia keys
Options:
  -profile  AWS profile to use to generate temporary credentials.(Required)
`
	return usage
}
