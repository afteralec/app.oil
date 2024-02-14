package interfaces

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2/middleware/session"
	redis "github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"petrichormud.com/app/internal/config"
	"petrichormud.com/app/internal/query"
)

func SetupShared() Shared {
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

	ib := InterfacesBuilder().Database(db).Redis(r).Sessions(s)

	if os.Getenv("DISABLE_SENDING_STONE") != "true" {
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

type Shared struct {
	Database   *sql.DB
	Redis      *redis.Client
	Queries    *query.Queries
	Sessions   *session.Store
	ClientConn *grpc.ClientConn
}

type sharedInterfacesBuilder struct {
	Interfaces Shared
}

func InterfacesBuilder() *sharedInterfacesBuilder {
	return new(sharedInterfacesBuilder)
}

func (b *sharedInterfacesBuilder) Database(db *sql.DB) *sharedInterfacesBuilder {
	b.Interfaces.Database = db
	b.Interfaces.Queries = query.New(db)
	return b
}

func (b *sharedInterfacesBuilder) ClientConn(conn *grpc.ClientConn) *sharedInterfacesBuilder {
	b.Interfaces.ClientConn = conn
	return b
}

func (b *sharedInterfacesBuilder) Redis(r *redis.Client) *sharedInterfacesBuilder {
	b.Interfaces.Redis = r
	return b
}

func (b *sharedInterfacesBuilder) Sessions(s *session.Store) *sharedInterfacesBuilder {
	b.Interfaces.Sessions = s
	return b
}

func (b *sharedInterfacesBuilder) Build() Shared {
	return b.Interfaces
}

func (i *Shared) Ping() {
	if err := i.Database.Ping(); err != nil {
		log.Fatal(err)
	}
	if err := i.Redis.Ping(context.Background()).Err(); err != nil {
		log.Fatal(err)
	}
}

func (i *Shared) Close() {
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
