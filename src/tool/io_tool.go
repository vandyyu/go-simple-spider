package tool
import (
	"os"
)
func ExistFile(fpath string) bool{
	fi, err := os.Stat(fpath)
	return err == nil && !fi.IsDir()
}
func ExistDir(fpath string) bool{
	fi, err := os.Stat(fpath)
	return err == nil && fi.IsDir()
}
