package templatefuncs

import "path/filepath"

// Example: "/path/to/file.txt" -> "to".
func GetDirName(path string) string {
	return filepath.Base(filepath.Dir(path))
}

// Example: "/path/to/file.txt" -> "file.txt".
func GetFileName(path string) string {
	return filepath.Base(path)
}

// Example: "/path/to/file.txt" -> ".txt".
func GetFileExt(path string) string {
	return filepath.Ext(path)
}

// Example: "/path/to/file.txt" -> true, "file.txt" -> false.
func IsAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}

// Example: JoinPaths("path", "to", "file.txt") -> "path/to/file.txt".
func JoinPaths(elem ...string) string {
	return filepath.Join(elem...)
}
