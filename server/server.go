package main

import (
	"log"
	"net/http"
	fd "github.com/AleksMVP/sybd/pkg/forum/delivery"
	fr "github.com/AleksMVP/sybd/pkg/forum/repository"
	fu "github.com/AleksMVP/sybd/pkg/forum/usecase"
	pd "github.com/AleksMVP/sybd/pkg/post/delivery"
	pr "github.com/AleksMVP/sybd/pkg/post/repository"
	pu "github.com/AleksMVP/sybd/pkg/post/usecase"
	sd "github.com/AleksMVP/sybd/pkg/service/delivery"
	sr "github.com/AleksMVP/sybd/pkg/service/repository"
	su "github.com/AleksMVP/sybd/pkg/service/usecase"
	td "github.com/AleksMVP/sybd/pkg/thread/delivery"
	tr "github.com/AleksMVP/sybd/pkg/thread/repository"
	tu "github.com/AleksMVP/sybd/pkg/thread/usecase"
	ud "github.com/AleksMVP/sybd/pkg/user/delivery"
	ur "github.com/AleksMVP/sybd/pkg/user/repository"
	uu "github.com/AleksMVP/sybd/pkg/user/usecase"
	vd "github.com/AleksMVP/sybd/pkg/vote/delivery"
	vr "github.com/AleksMVP/sybd/pkg/vote/repository"
	vu "github.com/AleksMVP/sybd/pkg/vote/usecase"
	"github.com/gorilla/mux"
	// "net/http/pprof"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

/*const (
	host     = "localhost"
	port     = 5432
	user     = "bober"
	password = "postgres"
	dbname   = "forum"
)*/

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "postgres"
)

func DBInit() *sql.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	return db
}


func RouterInit(db *sql.DB) *mux.Router {
	forumRepository := fr.NewForumRepository(db)
	forumUseCase := fu.NewForumUseCase(&forumRepository)
	forumDelivery := fd.NewForumDelivery(&forumUseCase)

	threadRepository := tr.NewThreadRepository(db)
	threadUseCase := tu.NewThreadUseCase(&threadRepository, &forumRepository)
	threadDelivery := td.NewThreadDelivery(&threadUseCase)

	postRepository := pr.NewPostRepository(db)
	postUseCase := pu.NewPostUseCase(&postRepository, &threadRepository)
	postDelivery := pd.NewPostDelivery(&postUseCase)

	serviceRepository := sr.NewServiceRepository(db)
	serviceUseCase := su.NewServiceUseCase(&serviceRepository)
	serviceDelivery := sd.NewServiceDelivery(&serviceUseCase)

	userRepository := ur.NewUserRepository(db)
	userUseCase := uu.NewUserUseCase(&userRepository, &forumRepository)
	userDelivery := ud.NewUserDelivery(&userUseCase)

	voteRepository := vr.NewVoteRepository(db)
	voteUseCase := vu.NewVoteUseCase(&voteRepository, &threadRepository)
	voteDelivery := vd.NewVoteDelivery(&voteUseCase)

	r := mux.NewRouter()

	r.HandleFunc("/api/forum/create", forumDelivery.PostForumCreate).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/details", forumDelivery.GetForumDetails).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/create", threadDelivery.PostForumCreateThread).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/threads", threadDelivery.GetForumThreads).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/users", userDelivery.GetForumUsers).Methods("GET")

	r.HandleFunc("/api/user/{slug}/create", userDelivery.PostUserCreate).Methods("POST")
	r.HandleFunc("/api/user/{slug}/profile", userDelivery.GetUserProfile).Methods("GET")
	r.HandleFunc("/api/user/{slug}/profile", userDelivery.PostUserProfile).Methods("POST")

	r.HandleFunc("/api/thread/{slug}/details", threadDelivery.GetThreadDetails).Methods("GET")
	r.HandleFunc("/api/thread/{slug}/details", threadDelivery.PostThreadDetails).Methods("POST")
	r.HandleFunc("/api/thread/{slug}/create", postDelivery.PostCreateNewPosts).Methods("POST")
	r.HandleFunc("/api/thread/{slug}/posts", postDelivery.GetThreadPosts).Methods("GET")
	r.HandleFunc("/api/thread/{slug}/vote", voteDelivery.PostThreadVote).Methods("POST")

	r.HandleFunc("/api/post/{slug}/details", postDelivery.PostPostDetails).Methods("POST")
	r.HandleFunc("/api/post/{slug}/details", postDelivery.GetPostDetails).Methods("GET")

	r.HandleFunc("/api/service/status", serviceDelivery.GetServiceStatus).Methods("GET")
	r.HandleFunc("/api/service/clear", serviceDelivery.PostServiceClear).Methods("POST")

	return r
}


func main() {
	db := DBInit()

	r := RouterInit(db)
	/*r.HandleFunc("/api/forum/create", postForumCreate).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/details", getForumDetails).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/create", postForumCreateThread).Methods("POST")
	r.HandleFunc("/api/forum/{slug}/threads", getForumThreads).Methods("GET")
	r.HandleFunc("/api/forum/{slug}/users", getForumUsers).Methods("GET")

	r.HandleFunc("/api/user/{slug}/create", postUserCreate).Methods("POST")
	r.HandleFunc("/api/user/{slug}/profile", getUserProfile).Methods("GET")
	r.HandleFunc("/api/user/{slug}/profile", postUserProfile).Methods("POST")

	r.HandleFunc("/api/thread/{slug}/details", getThreadDetails).Methods("GET")
	r.HandleFunc("/api/thread/{slug}/details", postThreadDetails).Methods("POST")
	r.HandleFunc("/api/thread/{slug}/create", postCreateNewPosts).Methods("POST")
	r.HandleFunc("/api/thread/{slug}/posts", getThreadPosts).Methods("GET")
	r.HandleFunc("/api/thread/{slug}/vote", postThreadVote).Methods("POST")

	r.HandleFunc("/api/post/{slug}/details", postPostDetails).Methods("POST")
	r.HandleFunc("/api/post/{slug}/details", getPostDetails).Methods("GET")

	r.HandleFunc("/api/service/status", getServiceStatus).Methods("GET")
	r.HandleFunc("/api/service/clear", postServiceClear).Methods("POST")

	r.HandleFunc("/debug/pprof/", pprof.Index)
    r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
    r.HandleFunc("/debug/pprof/profile", pprof.Profile)
    r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	
    // Manually add support for paths linked to by index page at /debug/pprof/
    r.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
    r.Handle("/debug/pprof/heap", pprof.Handler("heap"))
    r.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
    r.Handle("/debug/pprof/block", pprof.Handler("block"))*/

	port := "5000"

	var err error = nil

	log.Println("Launching at HTTP port " + port)
	err = http.ListenAndServe(":"+port, r)

	if err != nil {
		log.Fatal("Unable to launch server: ", err)
	}
}
