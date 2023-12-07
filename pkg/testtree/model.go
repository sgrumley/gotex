package testtree

type MenuItem struct {
	Title       string
	Children    []MenuItem
	Expanded    bool
	IsSubmenu   bool
	ParentIndex int // Index of the parent item, -1 for root items
}

type Model struct {
	Menu                []MenuItem
	CurrentIndex        int  // Index of the currently selected item in the main menu
	CurrentSubmenuIndex int  // Index of the currently selected item in the submenu
	InSubmenu           bool // Flag to indicate if the user is currently in a submenu
}
