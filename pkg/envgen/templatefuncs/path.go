package templatefuncs

import "path/filepath"

// GetDirName returns the directory name from a path.
// Example: "/path/to/file.txt" -> "to"
func GetDirName(path string) string {
	return filepath.Base(filepath.Dir(path))
}

// GetFileName returns the file name (with extension) from a path.
// Example: "/path/to/file.txt" -> "file.txt"
func GetFileName(path string) string {
	return filepath.Base(path)
}

// GetFileExt returns the file extension from a path.
// Example: "/path/to/file.txt" -> ".txt"
func GetFileExt(path string) string {
	return filepath.Ext(path)
}

// IsAbsolutePath checks if the path is absolute.
// Example: "/path/to/file.txt" -> true, "file.txt" -> false
func IsAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}

// JoinPaths joins path elements using the operating system's path separator.
// Example: JoinPaths("path", "to", "file.txt") -> "path/to/file.txt"
func JoinPaths(elem ...string) string {
	return filepath.Join(elem...)
}
