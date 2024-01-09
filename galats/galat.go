package galats

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"sync"
)

func Error400(w http.ResponseWriter) {
	NormalErrors()
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, "{\"message\":\"Bad Request\"}")
}

func Error500(w http.ResponseWriter) {
	NormalErrors()
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "{\"message\":\"Internal Server Error\"}")
}

type PesimisticLocking struct {
	sync.Mutex
	Tx *sql.Tx
	Db *sql.DB
}

func (p *PesimisticLocking) Done(w http.ResponseWriter) {
	defer p.Unlock()
	if r := recover(); r != nil {
		log.Println("Error catched", r)
		log.Println("Stack trace", string(debug.Stack()))
		if p.Tx != nil {
			if err := p.Tx.Rollback(); err != nil {
				panic(err)
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "{\"message\":\"Internal Server Error\"}")
	} else {
		if p.Tx != nil {
			if err := p.Tx.Commit(); err != nil {
				panic(err)
			}
		}
	}
	if p.Db != nil {
		if err := p.Db.Close(); err != nil {
			panic(err)
		}
	}
}

func Db500(w http.ResponseWriter, db *sql.DB) {
	if r := recover(); r != nil {
		log.Println("Error catched", r)
		log.Println("Stack trace", string(debug.Stack()))
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "{\"message\":\"Internal Server Error\"}")
	}
	if db != nil {
		if err := db.Close(); err != nil {
			panic(err)
		}
	}
}

func DbClose(db *sql.DB) {
	NormalErrors()
	if db != nil {
		if err := db.Close(); err != nil {
			panic(err)
		}
	} else {
		panic(fmt.Errorf("Failed to open connection DB"))
	}
}

func NormalErrors() {
	if r := recover(); r != nil {
		log.Println("Error catched", r)
		log.Println("Stack trace", string(debug.Stack()))
	}
}
