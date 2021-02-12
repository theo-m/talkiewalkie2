package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
	mrand "math/rand"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/theo-m/talkiewalkie/models"
	"github.com/theo-m/talkiewalkie/pb"
)

//go:generate sqlboiler psql --add-global-variants
// https://github.com/grpc/grpc-go/issues/3669#issuecomment-648433115

//go:generate protoc -I=/usr/local/include/google/protobuf -I=. --go_out=pb --go-grpc_out=pb --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false tw.proto

type contextKey int

const (
	dbKey contextKey = iota
)

type server struct {
	// pb.UnimplementedTalkieWalkieServer
}

var _ pb.TalkieWalkieServer = server{}

func (s server) Register(ctx context.Context, inp *pb.RegisterInput) (*pb.User, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.User{}, status.Errorf(codes.Internal, "could not access request metadata")
	}
	authTokens := md.Get("authToken")
	if len(authTokens) > 0 {
		return nil, status.Errorf(codes.AlreadyExists, "cannot register, there is already an authtoken")
	}

	db := ctx.Value(dbKey).(*sql.DB)

	u := &models.User{
		Username:    inp.Handle,
		Email:       null.String{String: inp.Email, Valid: true},
		Password:    null.String{String: inp.Password, Valid: true},
		Token:       null.String{String: usernameGen(), Valid: true},
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		ConnectedAt: time.Time{},
	}
	if err := u.Insert(ctx, db, boil.Infer()); err != nil {
		return &pb.User{}, status.Errorf(codes.Aborted, "error in inserting new user to db: %v", err)
	}
	return &pb.User{Uuid: u.UUID, Email: u.Email.String}, nil
}

//
//func (a server) Get(ctx context.Context, book *pb.AddressBook) (*pb.AddressBook, error) {
//	log.Println("hey")
//	md, ok := metadata.FromIncomingContext(ctx)
//	if !ok {
//		log.Println(ok)
//	}
//	log.Println(md.Get("custom-header-1"))
//	return &pb.AddressBook{}, nil
//}

func main() {
	//debug := *flag.Bool("debug", false, "")
	err := godotenv.Load(".secrets/dev.env")
	if err != nil {
		log.Fatal(err)
	}
	postgresMigrations()

	serv := grpc.NewServer(grpc.UnaryInterceptor(
		func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
			log.Printf("serving: %v", info.FullMethod)
			boil.DebugMode = true
			db, err := openDb()
			if err != nil {
				log.Printf("could not open db: %v", err)
				return
			}
			newCtx := context.WithValue(ctx, dbKey, db)
			h, err := handler(newCtx, req)
			return h, err
		}))
	service := server{}
	pb.RegisterTalkieWalkieServer(serv, service)

	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Panicf("could not listen to tcp:9090: %v", err)
	}

	if err = serv.Serve(lis); err != nil {
		log.Panicf("could not start the server: %v", err)
	}
}

func postgresMigrations() {
	m, err := migrate.New(
		"file://migrations",
		"postgres://theo:@localhost:5432/talkiewalkie?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Migrating")

	// XXX: https://github.com/golang-migrate/migrate/issues/179
	if err := m.Up(); err != nil && err.Error() != "no change" {
		log.Fatal(err)
	}
}

func buildContext(debug bool) http.HandlerFunc {
	boil.DebugMode = debug

	return func(w http.ResponseWriter, r *http.Request) {
		c := r.Context()
		db, err := openDb()
		if err != nil {
			http.Error(w, fmt.Sprintf("error opening db connection: %s", err), http.StatusInternalServerError)
			return
		}
		context.WithValue(c, "db", db)
		//
		//usernameC, err := r.Cookie("tw-username")
		//var username string
		//if err != nil {
		//	username = usernameGen()
		//	log.Printf("No cookie for username, generating random username: [%v]", username)
		//} else {
		//	username= usernameC.Value
		//}
		//log.Printf("Processing request for '%v'", username)
		//http.SetCookie(w, &http.Cookie{
		//	Name:       "",
		//	Value:      "",
		//	Path:       "",
		//	Domain:     "",
		//	Expires:    time.Time{},
		//	RawExpires: "",
		//	MaxAge:     0,
		//	Secure:     false,
		//	HttpOnly:   false,
		//	SameSite:   0,
		//	Raw:        "",
		//	Unparsed:   nil,
		//})
		//c.SetCookie(
		//	"tw-username",
		//	username,
		//	3600,
		//	"/",
		//	"localhost",
		//	false,
		//	false,
		//)
		//user, err := models.Users(qm.Where("username=?", username)).One(c, db)
		//if err != nil {
		//	log.Printf("Creating a new user for '%v'", username)
		//	user = &models.User{Username: username}
		//	err = user.Insert(c, db, boil.Infer())
		//	if err != nil {
		//		_ = c.AbortWithError(500, fmt.Errorf("couldn't get user for '%v': '%v'", username, err))
		//		return
		//	}
		//
		//}
		//
		//token, err := c.Cookie("tw-token")
		//if token == "" || err != nil || !user.Token.Valid {
		//	token = genToken()
		//	c.SetCookie(
		//		"tw-token",
		//		token,
		//		3600,
		//		"/",
		//		"localhost",
		//		false,
		//		false,
		//	)
		//	user.Token = null.String{String: token, Valid: true}
		//	_, err = user.Update(c, db, boil.Infer())
		//	if err != nil {
		//		_ = c.AbortWithError(500, fmt.Errorf("failed to update user token '%v': %v", user.Username, err))
		//		return
		//	}
		//}
		//if user.Token.Valid && user.Token.String != token {
		//	_ = c.AbortWithError(501, fmt.Errorf("cookie token: '%v' != db token: '%v' for %v", token, user.Token.String, user.Username))
		//	return
		//}
		//
		//c.Set("user", user)
		//c.Set("db", db)
		defer db.Close()
		//c.Next()
	}
}

func openDb() (*sql.DB, error) {
	db, err := sql.Open("postgres", "dbname=talkiewalkie user=theo sslmode=disable")
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}

func genToken() string {
	b := make([]byte, 128)
	_, _ = rand.Read(b)
	m := base64.StdEncoding.EncodeToString(b)
	return m
}

func loadWords(fn string) []string {
	file, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	words := make([]string, 0)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return words
}

var adverbs = loadWords("resources/adverbs.txt")
var adjectives = loadWords("resources/adjectives.txt")
var animals = loadWords("resources/animals.txt")

func usernameGen() string {
	adv := adverbs[mrand.Intn(len(adverbs))]
	adj := adjectives[mrand.Intn(len(adjectives))]
	ani := animals[mrand.Intn(len(animals))]
	const uname = "%v-%v-%v"
	return fmt.Sprintf(uname, adv, adj, ani)
}
