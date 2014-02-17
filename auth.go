package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"os/user"

	"github.com/polluxxx/ghubic"
	"github.com/polluxxx/goauth2"
)

const (
	ClientId     = "api_hubic_1366206728U6faUvDSfE1iFImoFAFUIfDRbJytlaY0"
	ClientSecret = "gXfu3KUIO1K57jUsW7VgKmNEhOWIbFdy7r8Z2xBdZn5K6SMkMmnU4lQUcnRy5E26"
	RedirectUrl  = "https://api.hubic.com"
)

type hubicConfig struct {
	AccessToken  string
	RefreshToken string
}

var Hubic *ghubic.HubicApi

func checkAuth() (account *ghubic.Account, err error) {
	Hubic, err := ghubic.NewHubicApi(ClientId, ClientSecret, RedirectUrl)
	if err != nil {
		return nil, err
	}

	config, err := loadConfig()
	if err != nil {
		uri, err := Hubic.GetAuthUrl("state")

		code := getCode(uri)

		account, err = Hubic.FinalizeAuth(code)
		if err != nil {
			return nil, err
		}

		err = saveConfig(&hubicConfig{
			AccessToken:  account.Token.AccessToken,
			RefreshToken: account.Token.RefreshToken,
		})
		if err != nil {
			return nil, err
		}
	} else {
		account, err = Hubic.GetAccountFromToken(&goauth2.OAuthToken{
			AccessToken:  config.AccessToken,
			RefreshToken: config.RefreshToken,
			Type:         "Bearer",
			Client:       Hubic.Client,
		})
		if err != nil {
			return nil, err
		}

		if account.Token.AccessToken != config.AccessToken {
			err = saveConfig(&hubicConfig{
				AccessToken:  account.Token.AccessToken,
				RefreshToken: account.Token.RefreshToken,
			})
			if err != nil {
				return nil, err
			}
		}

	}

	return account, nil
}

var getCodeTemplate = `Welcome to the hubicli command-tool

Before using any command, you will need to authenticate your account in hubiC.

In order to get access to your hubiC account, please login on this URL, and copy/paste the getted code.

    URL : {{.Url}}

Code ? `

func getCode(uri *url.URL) string {

	type CodeUrl struct {
		Url string
	}

	tmpl(os.Stdout, getCodeTemplate, CodeUrl{fmt.Sprintf("%s", uri)})
	answer := ""
	_, err := fmt.Scanf("%s", &answer)
	if err != nil {
		panic(err)
	}

	return answer
}

func saveConfig(hubicCfg *hubicConfig) error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("impossible to find your home directory: %s", err.Error())
	}
	directory := usr.HomeDir + "/.config/hubicli/"
	if err := os.MkdirAll(directory, os.ModeDir|0700); err != nil && err != os.ErrExist {
		return fmt.Errorf("impossible to create %s: %s", directory, err.Error())
	}

	b, err := json.MarshalIndent(hubicCfg, "", " ")
	fmt.Printf("\n\n%s\n\n", b)
	if err != nil {
		return fmt.Errorf("impossible to write json config %s", err.Error())
	}

	file, err := os.Create(directory + "config.json")
	if err != nil {
		return fmt.Errorf("impossible to write create file %s", err.Error())
	}
	defer file.Close()

	_, err = file.Write(b)
	if err != nil {
		return fmt.Errorf("impossible to write %s: %s", directory+"config.json", err.Error())
	}
	return nil
}

func loadConfig() (*hubicConfig, error) {

	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("impossible to find your home directory: %s", err.Error())
	}
	path := usr.HomeDir + "/.config/hubicli/config.json"

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var hubicCfg hubicConfig
	err = json.Unmarshal(file, &hubicCfg)
	if err != nil {
		return nil, err
	}

	return &hubicCfg, nil
}
