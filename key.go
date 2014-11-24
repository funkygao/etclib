package etclib

func keyPath(path ...string) string {
	key := "/" + DIR_PROJECT + "/" + project
	for _, name := range path {
		key += "/" + name
	}

	return key
}
