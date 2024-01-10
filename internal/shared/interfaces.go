package shared

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2/middleware/session"
	redis "github.com/redis/go-redis/v9"
	resend "github.com/resend/resend-go/v2"

	"petrichormud.com/app/internal/configs"
	"petrichormud.com/app/internal/queries"
)

func SetupInterfaces() Interfaces {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("SET GLOBAL local_infile=true;")
	if err != nil {
		log.Fatal(err)
	}

	opts := configs.Redis()
	r := redis.NewClient(&opts)

	s := session.New(configs.Session())

	rc := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	i := InterfacesBuilder().Database(db).Redis(r).Sessions(s).Resend(rc).Build()
	i.Ping()
	return i
}

type Interfaces struct {
	Database *sql.DB
	Redis    *redis.Client
	Queries  *queries.Queries
	Sessions *session.Store
	Resend   *resend.Client
}

type interfacesBuilder struct {
	Interfaces Interfaces
}

func InterfacesBuilder() *interfacesBuilder {
	return new(interfacesBuilder)
}

func (b *interfacesBuilder) Database(db *sql.DB) *interfacesBuilder {
	b.Interfaces.Database = db
	b.Interfaces.Queries = queries.New(db)
	return b
}

func (b *interfacesBuilder) Redis(r *redis.Client) *interfacesBuilder {
	b.Interfaces.Redis = r
	return b
}

func (b *interfacesBuilder) Sessions(s *session.Store) *interfacesBuilder {
	b.Interfaces.Sessions = s
	return b
}

func (b *interfacesBuilder) Resend(client *resend.Client) *interfacesBuilder {
	b.Interfaces.Resend = client
	return b
}

func (b *interfacesBuilder) Build() Interfaces {
	return b.Interfaces
}

func (i *Interfaces) Ping() {
	if err := i.Database.Ping(); err != nil {
		log.Fatal(err)
	}
	if err := i.Redis.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}
}

func (i *Interfaces) Close() {
	i.Database.Close()
	i.Redis.Close()
}

func SetupDB(db *sql.DB) error {
	_, err := db.Exec("SET GLOBAL local_infile=true;")
	if err != nil {
		return err
	}

	return nil
}

func PingDB(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return err
	}
	return nil
}
