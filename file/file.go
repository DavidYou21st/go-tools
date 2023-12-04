package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
)

//关于大文件的操作，为了避免一次性将整个文件加载到内存中造成内存溢出，我们需要将大文件切片成多个小的文件片段来操作。

const (
	defaultSpitFileSize = 10 * 1024 * 1024
)

// IsDir 判断指定的路径是否是目录
func IsDir(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

// IsFile 判断指定的路径是否是文件
func IsFile(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// SplitFile 将一个大文件分割成多个小文件
// filePath 大文件的文件路径
// prefix 是输出目标文件名称的前缀
// size 每个小文件的大小
func SplitFile(filePath, prefix string, size int) {
	if !IsFile(filePath) {
		panic("文件错误")
	}
	if size <= 0 {
		size = defaultSpitFileSize
	}
	//获取文件的文件名和扩展名
	filename, ext := Basename(filePath)
	//获取文件名前缀
	if len(prefix) == 0 {
		prefix = filename
	}
	dirname := filepath.Dir(filePath)
	fd, err := os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}
	defer fd.Close()

	buf := make([]byte, size)
	reader := bufio.NewReader(fd)
	index := 0
	for {
		subfilename := path.Join(dirname, "/", prefix, "-", strconv.Itoa(index), "."+ext)
		destFile, err := os.OpenFile(subfilename, syscall.O_CREAT|syscall.O_WRONLY, 0777)
		if err != nil {
			panic(err.Error())
		}

		writter := bufio.NewWriter(destFile)
		n, err := reader.Read(buf)
		if err == io.EOF { //文件已读取完
			break
		} else if err != nil {
			return
		}
		writter.Write(buf[n:])
		writter.Flush()
		index++
		destFile.Close()
	}

}

// Basename 分别返回文件的文件名和扩展名
func Basename(filepath string) (filename, ext string) {
	ext = path.Ext(filepath)
	basename := path.Base(filepath)
	filename = strings.TrimSuffix(basename, ext)
	return
}

// MergeFile 将一个目录下的小文件合并成一个大文件
// dirname 是要合并的小文件所在目录， filename 是输出目标文件名称
func MergeFile(dirname, filename string) {
	chunksPath := path.Join(strings.TrimLeft(dirname, "/"), "/")
	files, err := os.ReadDir(chunksPath)
	if err != nil {
		panic("目录不能正常访问")
	}
	// 排序
	filesSort := make(map[string]string)
	for _, f := range files {
		nameArr := strings.Split(f.Name(), "-")
		filesSort[nameArr[1]] = f.Name()
	}
	filesCount := len(files)
	if filesCount != len(filesSort) {
		panic("文件错误：有重复的文件")
	}

	saveFile := path.Join(chunksPath, filename)
	if exists, _ := PathExists(saveFile); exists {
		os.Remove(saveFile)
	}
	fs, _ := os.OpenFile(saveFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	defer fs.Close()

	var wg sync.WaitGroup
	wg.Add(filesCount)
	for i := 0; i < filesCount; i++ {
		// 要注意按顺序读取，否则文件就会损坏
		fileName := path.Join(chunksPath, "/"+filesSort[strconv.Itoa(i)])
		data, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Println(err)
		}
		fs.Write(data)

		wg.Done()
	}
	wg.Wait()
}

// 获取整体文件夹大小
func getDirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

// GetFileListBySuffix 获取指定后缀的文件
func GetFileListBySuffix(dirname, suffix string) ([]string, error) {
	if !IsDir(dirname) {
		return nil, fmt.Errorf("given path does not exist: %s", dirname)
	} else if IsFile(dirname) {
		return []string{dirname}, nil
	}

	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	num := len(infos)

	files := make([]string, 0, num)
	for _, v := range infos {
		if strings.HasSuffix(v.Name(), suffix) {
			files = append(files, v.Name())
		}
	}

	return files, nil
}

// GetFileListByPrefix 获取指定前缀的文件
func GetFileListByPrefix(dirname, prefix string) ([]string, error) {
	if !IsDir(dirname) {
		return nil, fmt.Errorf("given path does not exist: %s", dirname)
	} else if IsFile(dirname) {
		return []string{dirname}, nil
	}

	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	num := len(infos)

	files := make([]string, 0, num)
	for _, v := range infos {
		if strings.HasPrefix(v.Name(), prefix) {
			files = append(files, v.Name())
		}
	}

	return files, nil
}

// GetAllFile 获取指定目录下的所有文件
func GetAllFile(dirname string) ([]string, error) {
	if !IsDir(dirname) {
		return nil, fmt.Errorf("given path does not exist: %s", dirname)
	}
	dirname = strings.TrimSuffix(dirname, string(os.PathSeparator))
	infos, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	files := make([]string, len(infos))

	for _, v := range infos {
		if v.IsDir() {
			temp, err := GetAllFile(path.Join(dirname, string(os.PathSeparator), v.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, temp...)
		} else {
			files = append(files, v.Name())
		}
	}
	return files, nil
}
