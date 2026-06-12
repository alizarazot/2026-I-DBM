package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/alizarazot/2026-i-dbm/internal/constants"
	"github.com/alizarazot/2026-i-dbm/internal/database"
	"github.com/alizarazot/2026-i-dbm/internal/model"

	"charm.land/huh/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
	mongoOptions "go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ContextKey uint8

const (
	ContextKeyMongoClient ContextKey = iota
	ContextKeyDatabaseName
)

func main() {
	ctx := context.Background()

	mongoClient, err := mongo.Connect(mongoOptions.Client().ApplyURI(os.Getenv(constants.ENV_MONGODB_URI)))
	if err != nil {
		panic(fmt.Sprintf("unable to connect to MongoDB: %q", err))
	}
	defer func() {
		// TODO: This never runs on panics or [os.Exit] calls.
		if err := mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	ctx = context.WithValue(ctx, ContextKeyMongoClient, mongoClient)
	ctx = context.WithValue(ctx, ContextKeyDatabaseName, os.Getenv(constants.ENV_MONGODB_DATABASE))

	if os.Getenv(constants.ENV_MONGODB_DATABASE) == "" {
		panic(fmt.Sprintf("the environment variable %q is unset", constants.ENV_MONGODB_DATABASE))
	}

	flag.Parse()

	switch flag.Arg(0) {
	case "":
		flag.Usage()
		fmt.Println("You should specify a subcommand.")
		os.Exit(2)

	case "create-user":
		cmdCreateUser(ctx, flag.Args()[1:])

	default:
		fmt.Println("Unknown subcommand.")
		os.Exit(2)
	}
}

func cmdCreateUser(ctx context.Context, args []string) {
	var user model.User

	var birthdate, password string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Email Address").
				Placeholder("example@domain.com").
				Value(&user.Email),

			huh.NewSelect[model.UserRole]().
				Title("System Role").
				Options(
					huh.NewOption("Manager", model.UserRoleManager),
					huh.NewOption("Teacher", model.UserRoleTeacher),
					huh.NewOption("Student", model.UserRoleStudent),
				).
				Value(&user.Role),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("First Name").
				Value(&user.Info.FirstName),

			huh.NewInput().
				Title("Middle Name").
				Value(&user.Info.MiddleName),

			huh.NewInput().
				Title("First Surname").
				Value(&user.Info.FirstSurname),

			huh.NewInput().
				Title("Second Surname").
				Value(&user.Info.SecondSurname),

			huh.NewSelect[model.UserGenre]().
				Title("Genre").
				Options(
					huh.NewOption("Male", model.UserGenreMale),
					huh.NewOption("Female", model.UserGenreFemale),
					huh.NewOption("Other", model.UserGenreOther),
				).
				Value(&user.Info.Genre),

			huh.NewInput().
				Title("Birthdate (YYYY-MM-DD)").
				Placeholder("1990-01-01").
				Value(&birthdate).
				Validate(func(str string) error {
					_, err := time.Parse("2006-01-02", str)
					if err != nil {
						return fmt.Errorf("invalid date format, use YYYY-MM-DD")
					}
					return nil
				}),
		),

		huh.NewGroup(
			huh.NewInput().
				Title("Which's the user's password?").
				Prompt("* >").
				Value(&password),
		),
	)

	form.WithAccessible(true)
	if err := form.Run(); err != nil {
		panic(err)
	}

	var err error
	user.Info.Birthdate, err = time.Parse("2006-01-02", birthdate)
	if err != nil {
		panic(err)
	}

	userStore := database.NewUserStore(ctx.Value(ContextKeyMongoClient).(*mongo.Client), ctx.Value(ContextKeyDatabaseName).(string), constants.DB_COLLECTION_USER)

	if err := userStore.AddUser(ctx, &user, password); err != nil {
		panic(err)
	}
}
