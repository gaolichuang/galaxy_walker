package file

import (
    "bufio"
    "galaxy_walker/internal/gcodebase"
    "galaxy_walker/internal/gcodebase/conf"
    "fmt"
    "io/ioutil"
    LOG "galaxy_walker/internal/gcodebase/log"
    "os"
    "io"
    "bytes"
    "path/filepath"
    "strings"
    "crypto/md5"
    "galaxy_walker/internal/gcodebase/hash"
)

var CONF = conf.Conf

const (
    kDataPrefixDir               = "mdata"
    kLineReaderBufSize           = 1024 * 1024
    KSoftDeleteSubFix            = ".del"
    kWriteFileTemporarySubfixLen = 5
)

func ReadFileToString(name string) ([]byte, error) {
    name = GetConfFile(name)
    buf := bytes.NewBuffer(nil)
    f, err := os.Open(name)
    if err != nil {
        return []byte(""), err
    }
    defer f.Close()
    _, err = io.Copy(buf, f)
    if err != nil {
        return []byte(""), err
    }
    return buf.Bytes(), nil
}

func genWriteFileTemporaryName(name string) string {
    dir := filepath.Dir(name)
    base := filepath.Base(name)
    subfix := hash.GenerateKey(kWriteFileTemporarySubfixLen)
    newBase := "." + base + "." + subfix
    return filepath.Join(dir, newBase)
}
func WriteLineToFile(lines []string, name string) error {
    content := strings.Join(lines, "\n")
    tmpName := genWriteFileTemporaryName(name)
    err := ioutil.WriteFile(tmpName, []byte(content), 0644)
    if err != nil {
        return err
    }
    return os.Rename(tmpName, name)
}
func WriteStringToFile(content []byte, name string) error {
    tmpName := genWriteFileTemporaryName(name)
    err := ioutil.WriteFile(tmpName, content, 0644)
    if err != nil {
        return err
    }
    return os.Rename(tmpName, name)
}
func Exist(name string) bool {
    _, err := os.Stat(name)
    return !os.IsNotExist(err)
}
func IsDir(name string) bool {
    if !Exist(name) {
        return false
    }
    f, _ := os.Stat(name)
    return f.IsDir()
}
func IsRegular(name string) bool {
    if !Exist(name) {
        return false
    }
    f, _ := os.Stat(name)
    return f.Mode().IsRegular()
}
func DeleteIfExist(name string) error {
    name = GetConfFile(name)
    return os.Remove(name)
}
func SoftDeleteFile(filename string) error {
    filename = GetConfFile(filename)
    return os.Rename(filename, fmt.Sprintf("%s%s", filename, KSoftDeleteSubFix))
}
func DeletePath(filepath string, subfix string) (error, int) {
    filepath = GetConfFile(filepath)
    return deletePathInternal(filepath, subfix)
}
func MkDirAll(name string) error {
    name = GetConfFile(name)
    return os.MkdirAll(name, os.FileMode(0755))
}
func deletePathInternal(fpath string, subfix string) (error, int) {
    files, err := ioutil.ReadDir(fpath)
    if err != nil {
        return err, 0
    }
    delNum := 0
    for _, f := range files {
        if f.IsDir() {
            nfp := filepath.Join(fpath, f.Name())
            err, num := deletePathInternal(nfp, subfix)
            if err != nil {
                return err, num
            }
            delNum += num
        }
        if f.Mode().IsRegular() && strings.HasSuffix(f.Name(), subfix) {
            filename := fmt.Sprintf("%s/%s", fpath, f.Name())
            err := os.Remove(filename)
            if err != nil {
                return err, delNum
            }
            LOG.VLog(1).DebugTag("FILE", "Delete %s", filename)
            delNum += 1
        }
    }
    return nil, delNum
}

// only file, dir is not a file
func FileExist(name string) bool {
    s, err := os.Stat(name)
    return !os.IsNotExist(err) && (!(s != nil && s.IsDir()))
}
func FileLineReaderSoftly(filename string, comment string, f func(line string)) error {
    filename = GetConfFile(filename)
    // Open the file.
    fp, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer fp.Close()
    // Create a new Scanner for the file.
    scanner := bufio.NewScanner(fp)
    buf := make([]byte, kLineReaderBufSize)
    scanner.Buffer(buf, 100*kLineReaderBufSize)
    // Loop over all lines in the file and print them.
    for scanner.Scan() {
        line := scanner.Text()
        if strings.TrimSpace(line) == "" || strings.HasPrefix(line, comment) {
            continue
        }
        f(line)
    }
    return nil
}
func FileLineReader(filename string, comment string, f func(line string)) {
    filename = GetConfFile(filename)
    base.CHECK(Exist(filename), "File %s Not Exist.", filename)

    // Open the file.
    fp, err := os.Open(filename)
    base.CHECK(err == nil, "Read File Error %v.", err)
    defer fp.Close()
    // Create a new Scanner for the file.
    scanner := bufio.NewScanner(fp)
    buf := make([]byte, kLineReaderBufSize)
    scanner.Buffer(buf, 100*kLineReaderBufSize)
    // Loop over all lines in the file and print them.
    for scanner.Scan() {
        line := scanner.Text()
        if strings.TrimSpace(line) == "" || strings.HasPrefix(line, comment) {
            continue
        }
        f(line)
    }
}

func GetConfFile(s string) string {
    /*
       1. if absolute path, return
       2. if ConPathPrefix/s exist. return
       3. else  return ./s
    */
    if filepath.IsAbs(s) {
        return s
    }
    //cpp := *CONF.ConfPathPrefix
    //realFile := fmt.Sprintf("%s/%s/%s", *CONF.ConfPathPrefix, kDataPrefixDir, s)
    //if strings.HasSuffix(cpp, kDataPrefixDir) || strings.HasSuffix(cpp, kDataPrefixDir+"/") {
    //	realFile = fmt.Sprintf("%s/%s", *CONF.ConfPathPrefix, s)
    //}
    realFile := filepath.Join(*CONF.ConfPathPrefix, s)
    if Exist(realFile) {
        rf, err := filepath.Abs(realFile)
        if err != nil {
            return realFile
        }
        return rf
    } else {
        rf, err := filepath.Abs(s)
        if err != nil {
            return s
        }
        return rf
    }
}
func GetMD5FromFile(filename string) string {
    filename = GetConfFile(filename)
    buf := bytes.NewBuffer(nil)
    f, _ := os.Open(filename)
    defer f.Close()
    io.Copy(buf, f)
    return fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
}
func CalDirectoryMd5Sum(fpath, pathPrefix string, invalidSubfix string) map[string]string {
    // calculate md5 for a directory and save to list
    // return key is filename value is md5
    // format:    key:smart/sh.json value:d41d8cd98f00b204e9800998ecf8427e

    if !Exist(fpath) {
        LOG.Errorf("%s Not Exist", fpath)
        return nil
    }
    files, err := ioutil.ReadDir(fpath)
    if err != nil {
        return nil
    }
    retMd5s := make(map[string]string)
    for _, f := range files {
        absoluteFileName := filepath.Join(fpath, f.Name())
        relativeFileName := filepath.Join(pathPrefix, f.Name())
        if f.IsDir() {
            md5s := CalDirectoryMd5Sum(absoluteFileName, relativeFileName, invalidSubfix)
            if md5s != nil && len(md5s) > 0 {
                for k, v := range md5s {
                    retMd5s[k] = v
                }
            }
        }
        if f.Mode().IsRegular() {
            if strings.HasPrefix(f.Name(), ".") {
                continue
            }
            if !strings.HasSuffix(f.Name(), invalidSubfix) {
                md5sum := GetMD5FromFile(absoluteFileName)
                retMd5s[relativeFileName] = md5sum
            }
        }
    }
    return retMd5s
}

func ListFiles(fpath, validSubfix string) []string {
    if !Exist(fpath) {
        LOG.Errorf("%s Not Exist", fpath)
        return []string{}
    }
    files, err := ioutil.ReadDir(fpath)
    if err != nil {
        return []string{}
    }
    fileNames := make([]string, 0)
    for _, f := range files {
        if f.Mode().IsRegular() {
            if strings.HasSuffix(f.Name(), validSubfix) {
                fileNames = append(fileNames, f.Name())
            }
        }
    }
    return fileNames
}
func ListFilesByPrefix(fpath, prefix string, abs bool) (error, []string) {
    if !Exist(fpath) {
        return fmt.Errorf("%s Not Exist", fpath), nil
    }
    files, err := ioutil.ReadDir(fpath)
    if err != nil {
        return err, nil
    }
    fileNames := make([]string, 0)
    for _, f := range files {
        if f.Mode().IsRegular() {
            if strings.HasPrefix(f.Name(), prefix) {
                fname := f.Name()
                if abs {
                    fname = filepath.Join(fpath, f.Name())
                }
                fileNames = append(fileNames, fname)
            }
        }
    }
    return nil, fileNames
}
func FileLineNumAndByteCounts(filename string) (error, int, int) {
    filename = GetConfFile(filename)
    // Open the file.
    fp, err := os.Open(filename)
    if err != nil {
        return err, 0, 0
    }
    defer fp.Close()
    // Create a new Scanner for the file.
    scanner := bufio.NewScanner(fp)
    buf := make([]byte, kLineReaderBufSize)
    scanner.Buffer(buf, 100*kLineReaderBufSize)
    lineNum := 0
    bytes := 0
    // Loop over all lines in the file and print them.
    for scanner.Scan() {
        line := scanner.Bytes()
        lineNum++
        bytes += len(line)
    }
    return nil, lineNum, bytes
}
func CopyFile(srcFilePath, targetPath string) error {
    var err error
    var targetFilePath string
    if !Exist(srcFilePath) {
        return fmt.Errorf("%s NOT exist", srcFilePath)
    }
    srcInfo, err := os.Stat(srcFilePath)
    if err != nil {
        return err
    }
    if Exist(targetPath) {
        if IsDir(targetPath) {
            targetFilePath = filepath.Join(targetPath, srcInfo.Name())
        } else {
            targetFilePath = targetPath
        }
    } else {
        if !Exist(filepath.Dir(targetPath)) {
            if err = MkDirAll(filepath.Dir(targetPath)); err != nil {
                return err
            }
        }
        targetFilePath = targetPath
    }
    srcFile, err := os.Open(srcFilePath)
    if err != nil {
        return err
    }
    defer srcFile.Close()
    targetFile, err := os.Create(targetFilePath)
    if err != nil {
        return err
    }
    defer targetFile.Close()
    _, err = io.Copy(targetFile, srcFile)
    if err != nil {
        return err
    }
    return nil
}
func CopyDir(sourceDirPath, targetPath string) error {
    var err error
    if !IsDir(sourceDirPath) {
        return fmt.Errorf("%s is NOT Dir", sourceDirPath)
    }
    if Exist(targetPath) && !IsDir(targetPath) {
        return fmt.Errorf("%s Specified Must be a Dir", targetPath)
    }
    srcInfo, err := os.Stat(sourceDirPath)
    if err != nil {
        return err
    }
    if err = MkDirAll(filepath.Join(targetPath, srcInfo.Name())); err != nil {
        return err
    }
    files, err := ioutil.ReadDir(sourceDirPath)
    if err != nil {
        return err
    }
    for _, f := range files {
        if f.IsDir() {
            CopyDir(filepath.Join(sourceDirPath, f.Name()), filepath.Join(targetPath, srcInfo.Name()))
        } else {
            CopyFile(filepath.Join(sourceDirPath, f.Name()), filepath.Join(targetPath, srcInfo.Name()))
        }
    }
    return nil
}
func Copy(source, targetPath string) error {
    var err error
    if !Exist(source) {
        return fmt.Errorf("%s NOT exist", source)
    }
    if IsDir(source) {
        err = CopyDir(source, targetPath)
    } else {
        err = CopyFile(source, targetPath)
    }
    if err != nil {
        return err
    }
    return nil
}
