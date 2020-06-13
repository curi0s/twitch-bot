package main

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/curiTTV/twirgo"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	command struct {
		Name  string `bson:"name"`
		Value string `bson:"value"`
	}

	taler struct {
		Username string `bson:"username"`
		Amount   int    `bson:"amount"`
	}
)

var (
	commands []*command
	talers   []*taler
)

func findCommand(commandName string) (*command, error) {
	commandName = strings.ToLower(commandName)

	for _, command := range commands {
		if command.Name == commandName {
			return command, nil
		}
	}

	command := &command{
		Name: commandName,
	}
	commands = append(commands, command)

	return command, errors.New("Command not found")
}

func findTalers(username string) *taler {
	username = strings.ToLower(username)

	for _, taler := range talers {
		if taler.Username == username {
			return taler
		}
	}

	taler := &taler{
		Username: username,
	}
	talers = append(talers, taler)

	return taler
}

func normalizeUsername(username string) string {
	if strings.HasPrefix(username, "@") {
		username = strings.TrimLeft(username, "@")
	}
	return strings.ToLower(username)
}

func main() {
	token := os.Getenv("TOKEN")
	mongodbHost := os.Getenv("MONGODB_HOST")
	logLevel := os.Getenv("LOG_LEVEL")

	log := logrus.New()

	switch logLevel {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.Fatal("Invalid log level provided")
	}

	twirgoOptions := twirgo.Options{
		Username: "curi_BOT_",
		Channels: []string{"curi"},
		Log:      log,
		Token:    token,
	}

	t := twirgo.New(twirgoOptions)

	ch, err := t.Connect()
	if err == twirgo.ErrInvalidToken {
		log.Fatal(err)
	}

	client, err := mongo.NewClient(options.Client().
		ApplyURI("mongodb://" + mongodbHost + ":27017"))

	if err != nil {
		log.Fatal(err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.TODO())

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	cCommands := client.Database("bots").Collection("commands")
	cTalers := client.Database("bots").Collection("talers")

	for event := range ch {
		switch ev := event.(type) {
		case twirgo.EventConnected:
			// get all commands from db
			cur, err := cCommands.Find(context.TODO(), bson.M{}, options.Find())
			if err != nil {
				log.Fatal(err)
			}

			err = cur.All(context.TODO(), &commands)
			if err != nil {
				log.Fatal(err)
			}

			// get all taler from db
			cur, err = cTalers.Find(context.TODO(), bson.M{}, options.Find())
			if err != nil {
				log.Fatal(err)
			}

			err = cur.All(context.TODO(), &talers)
			if err != nil {
				log.Fatal(err)
			}

			log.Info("Connected to Twitch chat, loaded all command and taler")

		case twirgo.EventMessageReceived:
			parts := strings.Split(ev.Message.Content, " ")
			cmd := parts[0]

			if !strings.HasPrefix(cmd, "!") {
				continue
			}

			cmd = strings.TrimLeft(cmd, "!")

			switch cmd {
			case "quit":
				if ev.ChannelUser.IsBroadcaster {
					return
				}

			case "editcmd":
				if !ev.ChannelUser.IsBroadcaster && !ev.ChannelUser.IsMod {
					continue
				}

				command, _ := findCommand(parts[1])

				command.Value = strings.Join(parts[2:], " ")

				_, err = cCommands.UpdateOne(context.TODO(), bson.M{
					"name": command.Name,
				}, bson.M{"$set": command}, options.Update().SetUpsert(true))

				if err != nil {
					t.SendWhisper(ev.ChannelUser.User.Username, err.Error())
					log.Fatal(err)
				}

			case "deletecmd":
				if !ev.ChannelUser.IsBroadcaster && !ev.ChannelUser.IsMod {
					continue
				}

				command, _ := findCommand(parts[1])

				_, err := cCommands.DeleteOne(context.TODO(), bson.M{
					"name": command.Name,
				}, options.Delete())

				if err != nil {
					t.SendWhisper(ev.ChannelUser.User.Username, err.Error())
					log.Fatal(err)
				}

			case "addtaler":
				if !ev.ChannelUser.IsBroadcaster {
					continue
				}

				taler := findTalers(normalizeUsername(parts[1]))

				taler.Amount++

				_, err = cTalers.UpdateOne(context.TODO(), bson.M{
					"username": taler.Username,
				}, bson.M{"$set": taler}, options.Update().SetUpsert(true))

				if err != nil {
					t.SendWhisper(ev.ChannelUser.User.Username, err.Error())
					log.Fatal(err)
				}

			case "movetaler":
				if !ev.ChannelUser.IsBroadcaster {
					continue
				}

				fromUser := normalizeUsername(parts[1])
				toUser := normalizeUsername(parts[2])

				taler := findTalers(fromUser)
				taler.Username = toUser

				_, err = cTalers.UpdateOne(context.TODO(), bson.M{
					"username": taler.Username,
				}, bson.M{"$set": taler}, options.Update().SetUpsert(true))

				if err != nil {
					t.SendWhisper(ev.ChannelUser.User.Username, err.Error())
					log.Fatal(err)
				}

				_, err := cTalers.DeleteOne(context.TODO(), bson.M{
					"username": fromUser,
				}, options.Delete())

				if err != nil {
					t.SendWhisper(ev.ChannelUser.User.Username, err.Error())
					log.Fatal(err)
				}

			case "help":
				cmds := []string{
					"!taler", "!followage",
				}
				for _, cmd := range commands {
					cmds = append(cmds, "!"+cmd.Name)
				}

				sort.Strings(cmds)

				t.SendMessage(ev.Channel.Name, strings.Join(cmds, ", "))

			case "taler":
				taler := findTalers(ev.ChannelUser.User.Username)

				t.SendMessage(ev.Channel.Name, "@"+ev.ChannelUser.User.DisplayName+": Du hast "+strconv.Itoa(taler.Amount)+" Taler!")

			case "followage":
				resp, err := http.Get("https://api.crunchprank.net/twitch/followage/" + ev.Channel.Name + "/" + ev.ChannelUser.User.Username + "?precision=4")
				if err != nil {
					t.SendWhisper(ev.Channel.Name, err.Error())
					log.Fatal(err)
					continue
				}

				body, err := ioutil.ReadAll(resp.Body)
				resp.Body.Close()
				if err != nil {
					t.SendWhisper(ev.Channel.Name, err.Error())
					log.Fatal(err)
					continue
				}

				followage := string(body)
				message := "@" + ev.ChannelUser.User.DisplayName + ": Du folgst " + ev.Channel.Name + " schon " + string(body)

				if followage == "Follow not found" {
					message = "@" + ev.ChannelUser.User.DisplayName + ": Du folgst " + ev.Channel.Name + " nicht :("
				}

				t.SendMessage(ev.Channel.Name, message)

			default:
				command, err := findCommand(cmd)
				if err != nil {
					continue
				}

				t.SendMessage(ev.Channel.Name, command.Value)
			}
		}
	}
}
