package main

import (
	"flag"
	"gismart-rest-api/internal/app/apiserver"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	configPath string
	db_user string
	password string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
}

func main()  {
	flag.Parse()

	//	для подключения базы данных нужно ввести пользователя и пароль (db_user = "user=name", password = "password=1")
	//  после запустить Makefile командой $ make
	//  запустить миграции командой $ migrate -path migrations -database "postgres://user:password@localhost/restapi_dev?sslmode=disable" up
	//  параметры конфигурирования базы данных находятся в файле configs/apiserver.toml
	db_user = "user=nichi"
	password = "password=1"

	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}
	if err := apiserver.Start(config, db_user, password); err != nil{
		log.Fatal(err)
	}
}
