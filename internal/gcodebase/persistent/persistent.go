package persistent
import (
        "gcodebase/conf"
        "gcodebase/file"
        "gcodebase/log"
        "path/filepath"
        "io/ioutil"
        "fmt"
)
type Persistent interface {
        ListCategoryKeys(string) (error,[]string)
        // category / key / value
        SetWithCategory(string,string,string) error
        DelWithCategory(string,string) error
        GetWithCategory(string,string) (error, string)
        ExistWithCategory(string,string) bool
        ListWithCategory(string) []string

        Set(string,string) error
        Get(string) (error,string)
        Del(string) error
        Exist(string) bool
        List() []string
}
func NewPersistent() Persistent {
        return &PersistentByFile{
                PPath:getPersistentFilePath(*conf.Conf.FilePersistentPath),
        }
}
func InitPersistentTag(tag string) error {
        ppath := getPersistentFilePath(*conf.Conf.FilePersistentPath)
        fpath := filepath.Join(ppath,tag)
        err := file.MkDirAll(fpath)
        if err != nil {
                return fmt.Errorf("CreatePath File Err %v %s",err,fpath)
        }
        log.Infof("PersistentByFile Tag:%s use %s",tag, fpath)
        return nil
}

func init() {
        fpath := getPersistentFilePath(*conf.Conf.FilePersistentPath)
        err := file.MkDirAll(fpath)
        if err != nil {
                log.Fatalf("CreatePath File Err %v %s",err,fpath)
        }
        log.Infof("PersistentByFile use %s",fpath)
}
func getPersistentFilePath(name string) string {
        if filepath.IsAbs(name) {
                return name
        }
        return filepath.Join(*conf.Conf.ConfPathPrefix,name)
}
type PersistentByFile struct {
        PPath string
}

func (p *PersistentByFile)ListCategoryKeys(category string) (error,[]string) {
        fpath := filepath.Join(p.PPath,category)
        if !file.Exist(fpath) {
                return fmt.Errorf("%s not exist[%s]",category,fpath),nil
        }
        files, err := ioutil.ReadDir(fpath)
        if err != nil {
                return err,nil
        }
        values := make([]string,0)
        for _, f := range files {
                if f.Mode().IsRegular() {
                        values = append(values,f.Name())
                }
        }
        return nil,values
}
func (p *PersistentByFile)ExistWithCategory(category string, key string) bool {
        return file.Exist(filepath.Join(p.PPath,category,key))
}
func (p *PersistentByFile)ListWithCategory(category string) []string {
        return file.ListFiles(filepath.Join(p.PPath,category),"")
}
func (p *PersistentByFile)DelWithCategory(category string, key string) error {
        filename := filepath.Join(p.PPath,category,key)
        return file.DeleteIfExist(filename)
}
func (p *PersistentByFile)SetWithCategory(category string, key string, value string) error {
        err := file.MkDirAll(filepath.Join(p.PPath,category))
        if err != nil {
                return err
        }
        return file.WriteStringToFile([]byte(value),filepath.Join(p.PPath,category,key))
}
func (p *PersistentByFile)GetWithCategory(category string,key string) (error,string) {
        filename := filepath.Join(p.PPath,category,key)
        if !file.Exist(filename) {
                return fmt.Errorf("%s %s Not Exist[%s]",category,key,filename),""
        }
        b,e := file.ReadFileToString(filename)
        return e,string(b)
}
func (p *PersistentByFile)Set(key string,value string) error {
        return file.WriteStringToFile([]byte(value),filepath.Join(p.PPath,key))
}
func (p *PersistentByFile)Get(key string) (error,string) {
        filename := filepath.Join(p.PPath,key)
        if !file.Exist(filename) {
                return fmt.Errorf("%s Not Exist[%s]",key,filename),""
        }
        b,e := file.ReadFileToString(filename)
        return e,string(b)
}

func (p *PersistentByFile)Del(key string) error {
        filename := filepath.Join(p.PPath,key)
        return file.DeleteIfExist(filename)
}
func (p *PersistentByFile)Exist(key string) bool {
        return file.Exist(filepath.Join(p.PPath,key))
}
func (p *PersistentByFile)List() []string {
        return file.ListFiles(p.PPath,"")
}
