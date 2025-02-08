package scanner

// NOTE: thinking a tree should be made from the lists if needed
// would live in models? take scanner/tree.go
// func GenerateTree(p *models.Project) (*models.Tree, error) {
// 	tree := models.NewTree(p)
// 	dirNodes := make(map[string]*models.NodeTree)
// 	dirNodes[p.GetName()] = tree.RootNode
// }

// getPathComponents returns the path components relative to the root directory
// func getPathComponents(rootDir, pkgPath string) ([]string, error) {
// 	relPath, err := filepath.Rel(rootDir, pkgPath)
// 	if err != nil {
// 		return nil, fmt.Errorf("error getting relative path: %w", err)
// 	}
// 	return strings.Split(filepath.Clean(relPath), string(filepath.Separator)), nil
// }

// createDirectoryStructure creates intermediate directory nodes and returns the last parent node
// func createDirectoryStructure(rootPath string, components []string, rootNode *models.NodeTree, dirNodes map[string]*models.NodeTree) (*models.NodeTree, int) {
// 	currentPath := rootPath
// 	parentNode := rootNode
// 	maxLevel := 0
//
// 	// Process all components except the last one (which will be the package)
// 	for i, component := range components[:len(components)-1] {
// 		currentPath = filepath.Join(currentPath, component)
// 		level := i + 1
// 		maxLevel = level
//
// 		if node, exists := dirNodes[currentPath]; exists {
// 			parentNode = node
// 			continue
// 		}
//
// 		newDirNode := &models.NodeTree{
// 			Level: level,
// 			Data: models.DirectoryContent{
// 				Name: component,
// 				Path: currentPath,
// 			},
// 			Type:   models.NODE_TYPE_DIRECTORY,
// 			Parent: parentNode,
// 		}
// 		parentNode.Children = append(parentNode.Children, newDirNode)
//
// 		dirNodes[currentPath] = newDirNode
// 		parentNode = newDirNode
// 	}
//
// 	return parentNode, maxLevel
// }
