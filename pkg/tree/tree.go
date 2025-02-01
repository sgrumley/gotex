package tree

// NOTE: This tree should hold the structure of the data extracted from the project with the levels preserved as the file tree paths
// This should probably replace the project struct
// while implementing this, adding concurrency for branches would be nice

type Content interface {
	GetName() string
	GetValue() string
}

type Tree struct {
	RootNode *Node
}

// TODO: should the tree also contain a map for o(1) lookups?
func (t *Tree) Find(string) {}

// TODO: iterate tree and perform function on each node??
func (t *Tree) Traverse(f func() error) {}

type Node struct {
	Data     Content
	Children []*Node
	Parent   *Node
}

func (n *Node) AddChild(elem *Node) {}
