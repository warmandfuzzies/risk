package risk

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

var validStates map[string]struct{} = map[string]struct{}{
	"open":          {},
	"closed":        {},
	"accepted":      {},
	"investigating": {},
}

type Risk struct {
	Id          uuid.UUID `json:"id,omitempty"`
	State       string    `json:"state"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
}

var risks = map[uuid.UUID]Risk{}

func newRiskHandler(w http.ResponseWriter, r *http.Request) {
	var risk Risk
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&risk)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := validStates[risk.State]; !ok {
		http.Error(w, "Invalid state. Must be one of [open, closed, accepted, investigating]", http.StatusBadRequest)
		return
	}

	risk.Id = uuid.New()

	risks[risk.Id] = risk
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(&risk)
	if err != nil {
		http.Error(w, "Unable to encode.", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func listRiskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	riskValues := []Risk{}
	for _, risk := range riskValues {
		riskValues = append(riskValues, risk)
	}

	json.NewEncoder(w).Encode(riskValues)
}

func detailRiskHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "ID provided is not a valid UUID.", http.StatusBadRequest)
		return
	}

	risk := risks[id]

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&risk)

}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()
	r := mux.NewRouter()

	r.HandleFunc("/v1/risks", newRiskHandler).Methods("POST")
	r.HandleFunc("/v1/risks", listRiskHandler).Methods("GET")
	r.HandleFunc("/v1/risks/{id}", detailRiskHandler).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	srv := &http.Server{
		Addr: "localhost:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
