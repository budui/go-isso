package main

import (
	"fmt"
	"github.com/RayHY/go-isso/internal/pkg/conf"
	"log"
)

func main() {
	c, err := conf.Load("./configs/go-isso.toml")
	if err != nil {
		log.Fatalf("[FATA] Load Config Failed %v", err)
	}
	fmt.Println(c.Notify.Email.SMTP.Timeout)
	fmt.Println(c.Admin.Enable)
}
