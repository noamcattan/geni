package main

import (
	"context"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"github.com/noamcattan/geni/app"
	"github.com/noamcattan/geni/ent"
	_ "github.com/noamcattan/geni/ent/runtime"
	"os"

	"log"
)

func main() {
	client, err := ent.Open("sqlite3", "file:/tmp/geni.sql?cache=shared&_fk=1")
	if err != nil {
		log.Panicf("failed opening connection to sqlite: %v", err)
	}
	defer client.Close()
	if err := client.Schema.Create(context.Background()); err != nil {
		log.Panicf("failed creating schema resources: %v", err)
	}

	server := newServer(client)

	botToken, ok := os.LookupEnv("BOT_TOKEN")
	if !ok {
		log.Fatal("BOT_TOKEN env is missing")
	}

	bot, err := tg.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	router := app.NewApp(bot, client)
	ctx, cancel := context.WithCancel(context.Background())
	go router.Run(ctx)
	defer cancel()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	if err := server.Run(fmt.Sprintf(":%s", port)); err != nil {
		log.Printf(err.Error())
	}

	//http.HandleFunc("/", indexHandler)
	//port := os.Getenv("PORT")
	//if port == "" {
	//	port = "8080"
	//	log.Printf("Defaulting to port %s", port)
	//}
	//
	//log.Printf("Listening on port %s", port)
	//log.Printf("Open http://localhost:%s in the browser", port)
	//log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

//func indexHandler(w http.ResponseWriter, r *http.Request) {
//	if r.URL.Path != "/" {
//		http.NotFound(w, r)
//		return
//	}
//	_, err := fmt.Fprint(w, "Hello, World!")
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//}
