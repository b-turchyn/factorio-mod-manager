package main

import "bufio"
import "fmt"
import "io"
import "io/ioutil"
import "log"
import "os"
import "path/filepath"
import "strconv"
import "strings"
import "unicode"

func main() {
	// Support other OSs later...for now, Windows only
	modpacks := filepath.Join(os.Getenv("APPDATA"), "Factorio", "modpacks")
	mods := filepath.Join(os.Getenv("APPDATA"), "Factorio", "mods")
	packs := GetPacks(modpacks)
	
	fmt.Printf("%02d: %s\n", 0, "Base game")
	for i, f := range packs {
		fmt.Printf("%02d: %s\n", i + 1, f.Name())
	}
	
	pack := SelectPack(packs)
	var name string
	if pack != nil {
		name = pack.Name()
	} else {
		name = "Base"
	}
	fmt.Printf("Migrating to %s", name)
	
	RemoveDir(mods)
	
	if pack != nil {
		CopyDir(filepath.Join(os.Getenv("APPDATA"), "Factorio", "modpacks", pack.Name()), mods)
		fmt.Println("Done!")
	}
}

func RemoveDir(dir string) {
	os.RemoveAll(dir)
	fmt.Print(".")
}

func SelectPack(packs []os.FileInfo) os.FileInfo {
	var result int = -1
	reader := bufio.NewReader(os.Stdin)
	
	for result <= -1 || result > len(packs) {
		var err error
		fmt.Printf("Select Modpack: (0-%d): ", len(packs))
		input, _ := reader.ReadString('\n')
		input = strings.Map(func(r rune) rune {
			if unicode.IsSpace(r) {
				return -1
			}
			return r
		}, input)
		result, err = strconv.Atoi(input)
		
		if err != nil {
			log.Fatal(err)
		}
	}
	
	if result == 0 {
		return nil
	}
	
	return packs[result - 1]
}

func GetPacks(dir string) []os.FileInfo {
	files, err := ioutil.ReadDir(dir)
	var result[] os.FileInfo
	
	if err != nil {
		log.Fatal(err)
	} else {
		for _, f := range files {
			if f.IsDir() {
				result = append(result, f)
			}
		}
	}
	return result
}

// Sourced from https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}
	
	fmt.Print(".")

	return
}

// Sourced from https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04
// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}