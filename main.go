package main

import (
	"github.com/nepalsaurav/taskflow/models"
	"github.com/nepalsaurav/taskflow/pkg"
)

func main() {
	// if err := cmd.InitCMD(); err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Running taskflow server.. press Ctrl + c for cancel")

	// sig := make(chan os.Signal, 1)
	// signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// // block forever
	// <-sig

	// fmt.Println("Program shutdown")
	//
	//
	models.Migrate()
	// pkg.GetMailLog()
	//

	pkg.ImapLogin()
}
