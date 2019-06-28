//go:generate goversioninfo -icon=SHA512.ico

package main

import (
	"bufio"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var files []string
var extension string = ".sha512"
var bMode bool = false // if true enable binary mode.
var rMode bool = false // if true no calculate hashes, only checking.

func main() {
	if len(os.Args) == 1 {
		fmt.Println("use -b for enable binary mode, all output will be cast to binary type 'true' or 'false' if Error")
		fmt.Println("use -c for enable only checking mode, no calculate hashes, only checking.")
	}
	GetFileist(os.Args[1:], &files)
	for _, filename := range files {
		Hash(filename)
	}
	if !bMode {
		exit()
	}
}

func Hash(FilePath string) {
	rf, err := os.Open(FilePath)
	if err != nil {
		if bMode {
			fmt.Println("false")
		} else {
			fmt.Println("ERROR: Can't read file: ", FilePath)
		}
		return
	}
	defer rf.Close()
	h := sha512.New()
	_, err = io.Copy(h, rf)
	if err != nil {
		if bMode {
			fmt.Println("false")
		} else {
			fmt.Println("ERROR: Can't generate hash for file ", FilePath)
		}
		return
	}
	hFilePath := FilePath + extension
	if fileExists(hFilePath) && !rMode {
		rf, err := os.Open(hFilePath)
		if err != nil {
			if bMode {
				fmt.Println("false")
			} else {
				fmt.Println("ERROR: Can't read HASH file ", hFilePath)
			}
			return
		}
		defer rf.Close()
		hashdata := make([]byte, 128)
		_, err = rf.Read(hashdata)
		if err != nil {
			if bMode {
				fmt.Println("false")
			} else {
				fmt.Println("ERROR: Can't read HASH file ", hFilePath)
			}
		}
		if strings.ToUpper(hex.EncodeToString(h.Sum(nil))) == string(hashdata) {
			if bMode {
				fmt.Println("true")
			} else {
				fmt.Printf("Hash sum for file: \"%v\" matched\n", FilePath)
			}
		} else {
			if bMode {
				fmt.Println("false")
			} else {
				fmt.Printf("Hash sum for file \"%v\" NOT MATCHED!!!\n", FilePath)
			}
		}
	} else {
		wf, err := os.Create(hFilePath)
		if err != nil {
			if bMode {
				fmt.Println("false")
			} else {
				fmt.Println("ERROR: Can't create hash file: ", hFilePath)
			}
			return
		}
		defer wf.Close()
		wfb := bufio.NewWriter(wf)
		ok, err := wfb.WriteString(strings.ToUpper(hex.EncodeToString(h.Sum(nil))))
		if err != nil {
			if bMode {
				fmt.Println("false")
			} else {
				fmt.Println("ERROR: Can't write hash file: ", hFilePath)
			}
			return
		}
		if ok == 128 {
			if bMode {
				fmt.Println("true")
			} else {
				fmt.Printf("Hash sum for file \"%v\" writed to \"%v\"\n", FilePath, hFilePath)
			}
		} else {
			if bMode {
				fmt.Println("false")
			} else {
				fmt.Println("Hash of ", FilePath, " have incorrect length and may be wrong")
			}
		}
		wfb.Flush()
	}
}

func GetFileist(paths []string, files *[]string) {
	for _, path := range paths {
		if path == "-b" {
			bMode = true
		} else if path == "-r" {
			rMode = true
		} else {
			fd, err := os.Stat(path)
			if err != nil {
				if bMode {
					fmt.Println("false")
				} else {
					fmt.Println("ERROR: path ", path, " not exists")
				}
				return
			}
			if fd.IsDir() {
				err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
					if !info.IsDir() && (filepath.Ext(path) != extension) {
						*files = append(*files, path)
					}
					if err != nil {
						if bMode {
							fmt.Println("false")
						} else {
							fmt.Println("ERROR:", err)
						}
					}
					return nil
				})
				if err != nil {
					if bMode {
						fmt.Println("false")
					} else {
						fmt.Println("ERROR:", err)
					}
				}
			} else if filepath.Ext(path) != extension {
				*files = append(*files, path)
			}
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func exit() {
	fmt.Printf("Press Enter for exit...")
	b := make([]byte, 10)
	if _, err := os.Stdin.Read(b); err != nil {
		log.Fatal(err)
	}
}
