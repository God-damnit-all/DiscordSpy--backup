package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var (
	t      = flag.String("t", "", "discord token")
	victim = flag.String("v", "", "victim")
	ctx    context.Context
	pool   *sql.DB // database connection pool
)

func init() {
	flag.Parse()
}

func db() error {
	pool, _ = sql.Open("sqlite3", "discord.db")

	// If the database isn't pinged after 5 seconds, exit with error
	timeoutctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := pool.PingContext(timeoutctx); err != nil {
		return errors.Wrap(err, "Database not found, or unable to connect to databsae")
	}

	// Create table
	stmt, err := pool.Prepare("CREATE TABLE IF NOT EXISTS discord (uname TEXT NOT NULL, uid TEXT NOT NULL, msg TEXT NOT NULL)")
	if err != nil {
		fmt.Println(err)
	}
	stmt.Exec()

	return nil
}

func Discord() error {
	// Call DB
	_ = db()

	d, err := discordgo.New("Bot " + *t)
	if err != nil {
		return errors.Wrap(err, "Couldn't create session")
	}

	d.AddHandler(messageCreate)

	d.Identify.Intents = discordgo.IntentsGuildMessages

	err = d.Open()
	if err != nil {
		errors.Wrap(err, "Couldn't open ws connection")
	}

	fmt.Printf("Ready! To cancel, press CTRL-C")

	// Send an interrupt signal to interrupt with ctrl-c
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	s := <-sc
	fmt.Printf("Got signal: %s", s)

	return nil
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Author.Username == *victim {
		fmt.Printf("\n%s (%s): %s", *victim, m.Author.ID, m.Content)
		stmt, err := pool.Prepare("INSERT INTO discord (uname, uid, msg) VALUES (?, ?, ?)")
		if err != nil {
			panic(err)
		}
		stmt.Exec(*victim, m.Author.ID, m.Content)
	}
}

func main() {
	_ = Discord()
}
