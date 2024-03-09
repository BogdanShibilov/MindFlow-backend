package main

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type migrator struct {
	*migrate.Migrate
}

func new(migrationsPath, connUrl string) (*migrator, error) {
	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("%s?sslmode=disable", connUrl),
	)
	if err != nil {
		return nil, err
	}

	return &migrator{
		Migrate: m,
	}, nil
}

func (m *migrator) up() error {
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("There are no new migrations to apply")
			return nil
		}

		return err
	}

	fmt.Println("Migrated up")
	version, isDirty, err := m.Version()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Version: %v\nIs dirty: %v\n", version, isDirty)
	return nil
}

func (m *migrator) down() error {
	if err := m.Down(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println(("No more down migrations"))
			return nil
		}

		return err
	}

	fmt.Println("Migrated down")
	version, isDirty, err := m.Version()
	if err != nil && errors.Is(err, migrate.ErrNoChange) {
		panic(err)
	}
	fmt.Printf("Version: %v\nIs dirty: %v\n", version, isDirty)
	return nil
}

func main() {
	var connUrl, migrationsPath, way string

	flag.StringVar(&connUrl, "conn-url", "", "path to storage")
	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&way, "way", "up", "way to migrate")
	flag.Parse()

	if connUrl == "" {
		panic("conn-url is required")
	}
	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	m, err := new(migrationsPath, connUrl)
	if err != nil {
		panic(err)
	}

	switch way {
	case "up":
		err = m.up()
		if err != nil {
			panic(err)
		}
	case "down":
		err = m.down()
		if err != nil {
			panic(err)
		}
	}
}
