package file

import (
        "galaxy_walker/internal/github.com/fsnotify/fsnotify"
        "fmt"
)

func RegisterFileWatcher(filepath string, fn func(filename,op string)) error {
        if !Exist(filepath) {
                return fmt.Errorf("%s not exist",filepath)
        }
        watcher, err := fsnotify.NewWatcher()
        if err != nil {
                fmt.Errorf("%v",err)
                return err
        }
        //defer watcher.Close()

        //done := make(chan bool)
        go func() {
                for {
                        select {
                        case event := <-watcher.Events:
                                fmt.Println("event:", event)
                                fn(event.Name,event.Op.String())
                        case err := <-watcher.Errors:
                                fmt.Println("error:", err)
                        }
                }
        }()

        err = watcher.Add(filepath)
        if err != nil {
                fmt.Errorf("%v",err)
                return err
        }
        //<-done
        return nil
}
