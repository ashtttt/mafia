package main

import (
	"errors"
	"flag"
	"os/signal"
	"strings"
	"syscall"

	"os"

	"path/filepath"

	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
	"github.com/mitchellh/colorstring"
)

type KeyCommand struct {
	serial  string
	profile string
	creds   *sts.Credentials
}

var message = "[default] profile has been updated with temporary AWS session keys"

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
	} else if *profile == "default" {
		return errors.New("Can't use default profile, please use a different name for static keys")
	}
	if len(args) > 0 {
		flags.Usage()
	}
	k.profile = *profile
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
			colorstring.Printf("[green]%s \n", message)
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
	os.Setenv("AWS_PROFILE", k.profile)
	sess := session.Must(session.NewSession())
	svc := sts.New(sess)
	var params *sts.GetCallerIdentityInput
	resp, err := svc.GetCallerIdentity(params)

	if err != nil {
		return err
	}
	k.serial = strings.Replace(*resp.Arn, "user", "mfa", 1)
	os.Unsetenv("AWS_PROFILE")
	return nil
}

func (k *KeyCommand) getSessionToken() error {
	os.Setenv("AWS_PROFILE", k.profile)
	sess := session.Must(session.NewSession())
	svc := sts.New(sess)

	code, err := getOpt(k.profile)
	if err != nil {
		return err
	}

	params := &sts.GetSessionTokenInput{
		SerialNumber: aws.String(k.serial),
		TokenCode:    aws.String(code),
	}
	resp, err := svc.GetSessionToken(params)

	if err != nil {
		return err
	}
	k.creds = resp.Credentials
	os.Unsetenv("AWS_PROFILE")
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
	section, err := config.NewSection("default")
	if err != nil {
		return err
	}

	section.NewKey("aws_access_key_id", *k.creds.AccessKeyId)
	section.NewKey("aws_secret_access_key", *k.creds.SecretAccessKey)
	section.NewKey("aws_session_token", *k.creds.SessionToken)

	err = config.SaveTo(filename)
	if err != nil {
		return err
	}

	return nil

}

func (k *KeyCommand) commandHelp() string {
	var usage = `Usage: mafia keys
options:
  -profile  AWS profile name which has user's static keys.(Required)
`
	return usage
}
