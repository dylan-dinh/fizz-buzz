package main

import (
	"github.com/dylan-dinh/fizz-buzz/redis"
	"github.com/dylan-dinh/fizz-buzz/server"
	"log"
	"os"
	"sync"
)

//
//func handleSignals(fbs *FizzBuzzService) {
//	c := make(chan os.Signal, 1)
//	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGINT)
//
//mainloop:
//	for {
//		switch <-c {
//		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL:
//			fbs.server.Close()
//			break mainloop
//		}
//	}
//}

type FizzBuzzService struct {
	wg        *sync.WaitGroup
	rds       *redis.RedisService
	log       *log.Logger
	AppServer *server.AppServer
}

func (fbs *FizzBuzzService) StartFizzBuzz() error {
	fbs.wg.Add(2)
	go func() {
		defer fbs.wg.Done()
		fbs.rds = redis.StartRedis()
	}()

	go func() {
		defer fbs.wg.Done()
		server.StartServer(server.GetRouter(fbs.rds))
	}()

	fbs.wg.Wait()
	return nil
}

func main() {
	var err error
	var wg sync.WaitGroup

	defer func() {
		if err != nil {
			log.Fatalln(err.Error())
		}
	}()

	fbs := &FizzBuzzService{
		wg: &wg,
		log: log.New(os.Stdout, "[FIZZBUZZ]:	", 2)}

	err = fbs.StartFizzBuzz()
}
