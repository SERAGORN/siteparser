package main

import (
	"context"
	"fmt"
	"github.com/SERAGORN/siteparser/data/mysql"
	"github.com/SERAGORN/siteparser/data/parser"
	"github.com/SERAGORN/siteparser/handlers"
	"github.com/SERAGORN/siteparser/services/articlesrvc"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	cfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := NewConfig(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	validate := validator.New()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	sysCallContext, cancel := context.WithCancel(context.Background())
	go func() {
		sysCall := <-c
		fmt.Println(sysCall)
		cancel()
	}()

	mysqlHost := os.Getenv("PARSER_MYSQL_HOST")
	if mysqlHost != "" {
		cfg.Server.Db.Host = mysqlHost
	}

	sqlDB, err := mysql.InitMysql(sysCallContext, cfg.Server.Db.Host, cfg.Server.Db.Port, cfg.Server.Db.User, cfg.Server.Db.Password, cfg.Server.Db.Name)
	if err != nil {
		fmt.Println(err)
		return
	}

	articleRepository, err := mysql.NewArticleRepository(sqlDB)
	if err != nil {
		log.Fatal("problem while trying to initialize article repository")
	}

	parserRepository, err := parser.NewParserRepository()
	if err != nil {
		log.Fatal("problem while trying to initialize parser repository")
	}

	articleService, err := articlesrvc.NewArticleService(articleRepository, parserRepository)
	if err != nil {
		log.Fatal("problem while trying to initialize article service")
	}

	router, err := handlers.MakeRoutes(&handlers.RouterDependencies{
		ArticleService: articleService,
		Validate:       validate,
		MySql:          sqlDB,
	})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		fmt.Println("server starts:", cfg.Server.Port)
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen and serve error")
		}
	}()

	<-sysCallContext.Done()

	ctxShutDown, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		log.Fatal("server shutdown failed")
	}
}
