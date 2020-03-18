package main

import (
	"fmt"
	"github.com/TylerGrey/graphql-chat-sample/resolver"
	"github.com/TylerGrey/graphql-chat-sample/schema"
	"github.com/go-redis/redis"
	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/graph-gophers/graphql-transport-ws/graphqlws"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var httpPort = 8080

func init() {
	port := os.Getenv("HTTP_PORT")
	if port != "" {
		var err error
		httpPort, err = strconv.Atoi(port)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// graphiql handler
	http.HandleFunc("/", http.HandlerFunc(graphiql))

	// init graphQL schema
	s, err := graphql.ParseSchema(schema.GetRootSchema(), resolver.NewResolver(client))
	if err != nil {
		panic(err)
	}

	// graphQL handler
	graphQLHandler := graphqlws.NewHandlerFunc(s, &relay.Handler{Schema: s})
	http.HandleFunc("/graphql", graphQLHandler)

	// start HTTP server
	log.Printf("Listening for requests on http://localhost:%d", httpPort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil); err != nil {
		panic(err)
	}
}

var graphiql = func(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("graphiql").Parse(`
  	<!DOCTYPE html>
	<html>
	
	<head>
	  <meta charset=utf-8/>
	  <meta name="viewport" content="user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, minimal-ui">
	  <title>GraphQL Playground</title>
	  <link rel="stylesheet" href="//cdn.jsdelivr.net/npm/graphql-playground-react/build/static/css/index.css" />
	  <link rel="shortcut icon" href="//cdn.jsdelivr.net/npm/graphql-playground-react/build/favicon.png" />
	  <script src="//cdn.jsdelivr.net/npm/graphql-playground-react/build/static/js/middleware.js"></script>
	</head>
	
	<body>
	  <div id="root">
		<style>
		  body {
			background-color: rgb(23, 42, 58);
			font-family: Open Sans, sans-serif;
			height: 90vh;
		  }
	
		  #root {
			height: 100%;
			width: 100%;
			display: flex;
			align-items: center;
			justify-content: center;
		  }
	
		  .loading {
			font-size: 32px;
			font-weight: 200;
			color: rgba(255, 255, 255, .6);
			margin-left: 20px;
		  }
	
		  img {
			width: 78px;
			height: 78px;
		  }
	
		  .title {
			font-weight: 400;
		  }
		</style>
		<img src='//cdn.jsdelivr.net/npm/graphql-playground-react/build/logo.png' alt=''>
		<div class="loading"> Loading
		  <span class="title">GraphQL Playground</span>
		</div>
	  </div>
	  <script>window.addEventListener('load', function (event) {
		  GraphQLPlayground.init(document.getElementById('root'), {
			endpoint: '/graphql',
			subscriptionEndpoint: '/graphql'
		  })
		})</script>
	</body>
	
	</html>
  `))
	t.Execute(w, httpPort)
}
