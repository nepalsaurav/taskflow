package main

import "github.com/nepalsaurav/postfix_admin/models"

func main() {

	// maildir := pkg.Maildir{}

	// watcher, err := fsnotify.NewWatcher()
	// if err != nil {
	// 	log.Printf("error while opening watcher for maildir, err: %v\n", err)
	// }
	// defer watcher.Close()

	// maildir.WatchMailDir(watcher)

	// // block main go routine forever
	// <-make(chan struct{})
	//
	//

	models.Migrate()

}
