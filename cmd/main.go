package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/rudyeila/go-bring-api/src"
	"github.com/rudyeila/go-bring-api/src/model"
)

func main() {
	options := &slog.HandlerOptions{Level: slog.LevelInfo}
	log := slog.New(slog.NewJSONHandler(os.Stdout, options))

	cfg, err := src.NewConfig()
	if err != nil {
		log.Error(err.Error())
	}

	bring := src.New(cfg, log)
	err = bring.Login()
	if err != nil {
		log.Error(err.Error())
	}

	listsRes, err := bring.GetLists()
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Println(listsRes)
	list := findList(listsRes, "Test")
	if list == nil {
		log.Error("Test list not found...")
	}

	testList, err := bring.GetList(list.ListUuid)
	if err != nil {
		log.Error(err.Error())
	}
	fmt.Println(testList)

	err = bring.AddItem(testList.Uuid, model.ListItem{
		Specification: "TestItem",
		Name:          "250gr",
	})
	if err != nil {
		log.Error(err.Error())
	}
}

func findList(listsRes *model.GetListsResponse, name string) *model.List {
	for _, list := range listsRes.Lists {
		if list.Name == name {
			return &list
		}
	}

	return nil
}
