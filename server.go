package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/ks6088ts/graphql-server-go/graph"
	"github.com/ks6088ts/graphql-server-go/graph/generated"

	_ "github.com/lib/pq"
)

const defaultPort = "8080"

var (
	version  string
	revision string
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, err := sql.Open("postgres", "host=localhost port=5432 user=user password=password dbname=db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{
		Db: db,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("version=%v, revision=%v", version, revision)
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
