package etclib

func keyPath(path ...string) string {
	if project == "" {
		panic("empty project name")
	}

	key := projectPath()
	for _, name := range path {
		key += "/" + name
	}

	return key
}

func nodePath(nodeType, addr string) string {
	return keyPath(DIR_NODES, nodeType, addr)
}

func nodeRoot(nodeType string) string {
	return keyPath(DIR_NODES, nodeType)
}

func projectPath() string {
	return "/" + DIR_PROJECT + "/" + project
}
