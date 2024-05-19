package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/rudyeila/go-bring-client/bring"
	"github.com/rudyeila/go-bring-client/bring/model"
)

type IngredientWithAmount struct {
	Id     string   `json:"id"`
	Uuid   string   `json:"uuid"`
	Name   string   `json:"name"`
	Amount *float64 `json:"amount"`
	Unit   string   `json:"unit"`
}

func main() {
	options := &slog.HandlerOptions{Level: slog.LevelInfo}
	log := slog.New(slog.NewJSONHandler(os.Stdout, options))

	cfg, err := bring.NewConfig()
	if err != nil {
		log.Error(err.Error())
	}

	bring := bring.New(cfg, log)
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

	ingredientsWithAmount := make([]IngredientWithAmount, 0)
	err = ReadJSONFromFile("ingredients.json", &ingredientsWithAmount)
	if err != nil {
		log.Error(err.Error())
	}

	err = AddIngredientsToList(bring, ingredientsWithAmount, testList.Uuid)
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

func AddIngredientsToList(bring *bring.Bring, ingredients []IngredientWithAmount, listID string) error {
	for _, ingr := range ingredients {
		spec := ""
		if ingr.Amount != nil {
			spec = strconv.FormatFloat(*ingr.Amount, 'f', -1, 64) + ingr.Unit
		}

		err := bring.AddItem(listID, ingr.Name, spec)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadJSONFromFile(filename string, v interface{}) error {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	// Read the file's content
	err = json.NewDecoder(file).Decode(&v)
	if err != nil {
		return fmt.Errorf("could not read file: %v", err)
	}

	return nil
}
