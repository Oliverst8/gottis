package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

func write_and_open(note_name string) {
	println("Vimming")
	cmd := exec.Command("vim", note_name)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func getTree() *tview.TreeView {
	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(file.IsDir())
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, rootDir)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			add(node, path)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	return tree
}

func Draw1() {

	app := tview.NewApplication()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		k := event.Key()
		if k == tcell.KeyEnter {
			write_and_open("test name")
		}
		return event
	})

	box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	if err := app.SetRoot(box, true).SetFocus(box).Run(); err != nil {
		panic(err)
	}
}

func Draw() {
	app := tview.NewApplication()

	home := tview.NewBox().SetBorder(true).SetTitle("GOttis - Your Kattis CLI")
	files := tview.NewBox().SetBorder(true).SetTitle("Files")
	log := tview.NewBox().SetBorder(true).SetTitle("Log")
	help := tview.NewBox().SetBorder(true).SetTitle("Help")

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(files, 0, 1, false).
			AddItem(log, 0, 1, false).
			AddItem(help, 0, 1, false), 0, 1, false)

	home.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		flex.SetRect(x+2, y+2, width-5, height-5) // Adjust inner layout size
		flex.Draw(screen)
		return x, y, width, height
	})

	tree := getTree()

	files.SetDrawFunc(func(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
		tree.SetRect(x+2, y+2, width-5, height-5)
		tree.Draw(screen)
		return x, y, width, height
	})

	err := app.SetRoot(home, true).SetFocus(tree).Run()
	if err != nil {
		panic(err)
	}
}

// Show a navigable tree view of the current directory.
func Draw2() {
	rootDir := "."
	root := tview.NewTreeNode(rootDir).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}
		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name())).
				SetSelectable(file.IsDir())
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	}

	// Add the current directory to the root node.
	add(root, rootDir)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show files in this directory.
			path := reference.(string)
			add(node, path)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	if err := tview.NewApplication().SetRoot(tree, true).Run(); err != nil {
		panic(err)
	}
}
