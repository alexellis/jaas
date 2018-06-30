package server

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/alexellis/jaas/pkg/proxy"
	"github.com/alexellis/jaas/pkg/task"
	jtypes "github.com/alexellis/jaas/pkg/types"
	"github.com/openfaas/faas/gateway/handlers"
	ftypes "github.com/openfaas/faas/gateway/types"
)

func NewJaaSServer(port int, timeout time.Duration, jaasConfig jtypes.JaaSConfig) JaaSServer {
	return JaaSServer{
		Port:    port,
		Timeout: timeout,
		Config:  jaasConfig,
	}
}

type JaaSServer struct {
	Port    int
	Timeout time.Duration
	Server  *http.Server
	Config  jtypes.JaaSConfig
}

type HTTPServer interface {
	Start(stopCh chan interface{}) error
}

func (j *JaaSServer) Start(stopCh chan interface{}) error {

	j.Server = &http.Server{
		Addr:           fmt.Sprintf(":%d", j.Port),
		ReadTimeout:    j.Timeout,
		WriteTimeout:   j.Timeout,
		MaxHeaderBytes: 1 << 20, // Max header of 1MB
	}

	var creds *ftypes.BasicAuthCredentials

	for _, entry := range j.Config.Auths {

		if strings.Contains(entry.Address, "127.0.0.1") {
			passwordPlain, _ := hex.DecodeString(entry.Password)
			creds = &ftypes.BasicAuthCredentials{
				User:     entry.Username,
				Password: string(passwordPlain),
			}
			break
		}
	}
	if creds == nil {
		return fmt.Errorf("No credentials found for 127.0.0.1")
	}

	http.HandleFunc("/ping", handlers.DecorateWithBasicAuth(handler(), creds))
	http.HandleFunc("/run", handlers.DecorateWithBasicAuth(runHandler(), creds))
	http.HandleFunc("/list", handlers.DecorateWithBasicAuth(listHandler(), creds))

	if err := j.Server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Error ListenAndServe: %v", err)

		close(stopCh)

	}

	return nil
}

func handler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Ping")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	}
}

func runHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Run")

		if r.Body != nil {
			defer r.Body.Close()

			payloadBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				sendErr(w, err)
				return
			}

			req := proxy.RunRequest{}

			unmarshalErr := json.Unmarshal(payloadBytes, &req)
			if unmarshalErr != nil {
				sendErr(w, unmarshalErr)
				return
			}
			if req.Job == nil {
				sendErr(w, fmt.Errorf("job was nil in request"))
				return
			}

			fmt.Printf("Run %s\n", req.Job.Image)

			createStatus, err := task.Create(
				req.Job.Image,
				req.Job.Command,
			)

			if err != nil {
				sendErr(w, err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok", "job": "` + req.Job.Image + `", "ID": "` + createStatus.ID + `"}`))
		}

	}
}

func sendErr(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(`{"status": "error", "msg": "` + err.Error() + `"}`))
}

func listHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("List")

		items, err := task.List()
		if err != nil {
			sendErr(w, err)
			return
		}

		bytesOut, marshalErr := json.Marshal(items)
		if marshalErr != nil {
			sendErr(w, marshalErr)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(bytesOut)
	}
}
