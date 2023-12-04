package data

import (
	"fmt"
)

type dummyData struct {
	project   string
	folders   []string
	files     []string
	tests     []string
	testCases []string
}

/*
cmd/main_test.go
internal/service/accounts/service_test.go
internal/service/model/service_test.go
internal/store/model_test.go
*/

// to best represent this data I am thinking a tree structure
type Menu struct {
	Name    string
	SubMenu []*Menu
}

func NewMenuItem(name string) *Menu {
	return &Menu{
		Name:    name,
		SubMenu: []*Menu{},
	}
}

func (m *Menu) AddSubMenu(menu *Menu) {
	m.SubMenu = append(m.SubMenu, menu)
}

func PrintMenu(menu *Menu, level int) {
	prefix := ""
	for i := 0; i < level; i++ {
		prefix += "|--"
	}

	prefix += " "

	fmt.Println(prefix + menu.Name)
	for _, subMenu := range menu.SubMenu {
		PrintMenu(subMenu, level+1)
	}
}

func LoadDummyData() (*Menu, string) {
	menu := NewMenuItem("Tests")

	cmd := NewMenuItem("cmd/")
	cmd.AddSubMenu(NewMenuItem("main_test.go"))

	internal := NewMenuItem("internal/service")

	accounts := NewMenuItem("accounts/")
	accounts.AddSubMenu(NewMenuItem("service_test.go"))
	internal.AddSubMenu(accounts)

	models := NewMenuItem("models/")
	models.AddSubMenu(NewMenuItem("service_test.go"))
	internal.AddSubMenu(models)

	store := NewMenuItem("internal/store")
	store.AddSubMenu(NewMenuItem("model_test.go"))
	store.AddSubMenu(NewMenuItem("account_test.go"))
	internal.AddSubMenu(store)

	menu.AddSubMenu(cmd)
	menu.AddSubMenu(internal)

	return menu, "Retinalytics"
}
