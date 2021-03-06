package main

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"net/http"
	"os"
)

func ServerCommand() cli.Command {
	return cli.Command{
		Name:        "server",
		Usage:       "run web app",
		Description: "Use enviroment variable PORT to set which port to listen to.",
		ArgsUsage:   "",
		Action:      ServerAction,
		Flags:       []cli.Flag{},
	}
}

func ServerAction(c *cli.Context) {
	router := httprouter.New()
	router.POST("/convert", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		file, _, err := r.FormFile("messages")
		defer file.Close()
		WebCheck(w, err)

		b, err := ioutil.ReadAll(file)
		WebCheck(w, err)
		fbData := FromHTML(string(b))

		w.Header().Set("Content-type", "application/json")
		w.Write([]byte(ToJSON(fbData, false)))
	})
	router.ServeFiles("/*filepath", http.Dir("public"))

	addr := GetAddr()

	fmt.Println("Listening on http://localhost" + addr)
	check(http.ListenAndServe(addr, handlers.CompressHandler(router)))
}

func GetAddr() string {
	if port := os.Getenv("PORT"); len(port) > 0 {
		return ":" + port
	} else {
		return ":3000"
	}
}

func WebCheck(w http.ResponseWriter, e error) {
	if e != nil {
		fmt.Fprintf(w, "%v", e)
		panic(e)
	}
}
