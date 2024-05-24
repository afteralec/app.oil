package service

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2/middleware/session"
	html "github.com/gofiber/template/html/v2"
	redis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"petrichormud.com/app/internal/config"
	"petrichormud.com/app/internal/query"
	"petrichormud.com/app/web"
)

func NewInterfaces() Interfaces {
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("SET GLOBAL local_infile=true;")
	if err != nil {
		log.Fatal(err)
	}

	opts := config.Redis()
	r := redis.NewClient(&opts)

	s := session.New(config.Session())

	t := web.ViewsEngine()

	ib := InterfacesBuilder().Database(db).Redis(r).Sessions(s).Templates(t)

	if os.Getenv("DISABLE_SENDING_STONE") != "true" {
		// TODO: Migrate this to grpc.NewClient
		conn, err := grpc.Dial(os.Getenv("SENDING_STONE_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}
		ib.ClientConn(conn)
	}

	i := ib.Build()
	i.Ping()

	return i
}

type Interfaces struct {
	Database   *sql.DB
	Redis      *redis.Client
	Queries    *query.Queries
	Sessions   *session.Store
	ClientConn *grpc.ClientConn
	Templates  *html.Engine
}

type interfacesBuilder struct {
	Interfaces Interfaces
}

func InterfacesBuilder() *interfacesBuilder {
	return new(interfacesBuilder)
}

func (b *interfacesBuilder) Database(db *sql.DB) *interfacesBuilder {
	b.Interfaces.Database = db
	b.Interfaces.Queries = query.New(db)
	return b
}

func (b *interfacesBuilder) ClientConn(conn *grpc.ClientConn) *interfacesBuilder {
	b.Interfaces.ClientConn = conn
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

func (b *interfacesBuilder) Templates(e *html.Engine) *interfacesBuilder {
	b.Interfaces.Templates = e
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
	if i.Database != nil {
		i.Database.Close()
	}
	if i.Redis != nil {
		i.Redis.Close()
	}
	if i.ClientConn != nil {
		i.ClientConn.Close()
	}
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
