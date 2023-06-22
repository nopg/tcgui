package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os/exec"
	"strings"

	"github.com/danielgtaylor/huma"
	"github.com/danielgtaylor/huma/cli"
	"github.com/danielgtaylor/huma/responses"
)

type Test struct {
	Body io.Reader
}

func test(ctx huma.Context) {
	ctx.Header().Set("Content-Type", "text/html")
	ctx.Write([]byte("HELLO TEST"))
}

func set_interfaces(ctx huma.Context, latency string, loss string, jitter string, bandwidth string) {
	ctx.Header().Set("Content-Type", "text/plain")
	log.Println(latency, loss, jitter, bandwidth)

	args := fmt.Sprintf("%s %s %s %s", latency, loss, jitter, bandwidth)
	out, err := exec.Command("./static/mycommands.sh", args).Output()
	if err != nil {
		log.Println(err)
	}
	resp := fmt.Sprintf("%s", out)

	ctx.Write([]byte(resp))

	// body, err := ioutil.ReadFile("./static/updated.html")
	// if err != nil {
	// 	log.Fatalf("unable to read file: %v", err)
	// }
	// ctx.Header().Set("Content-Type", "text/html")
	// ctx.Write(body)
}

// Main Page
func main_page(ctx huma.Context) {
	body, err := ioutil.ReadFile("./static/index.html")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	ctx.Header().Set("Content-Type", "text/html")
	ctx.Write(body)
}

func set(ctx huma.Context, data Test) {
	ctx.Header().Set("Content-Type", "text/plain")
	//body := fmt.Sprintf("latency: %d\njitter: %d\nbandwidth: %d\nloss: %d", data.Body.Latency, data.Body.Jitter, data.Body.Bandwidth, data.Body.Loss)
	rawbody := data.Body
	rawparams := new(strings.Builder)
	n, _ := io.Copy(rawparams, rawbody)
	log.Println(rawparams.String())

	myurl := fmt.Sprintf("https://x.com/?%s", rawparams.String())
	params, err := url.Parse(myurl)
	if err != nil {
		log.Fatal(err)
	}
	values := params.Query()
	log.Println(values)
	log.Println(n)

	latency := values["latency"][0]
	bandwidth := values["bandwidth"][0]
	jitter := values["jitter"][0]
	loss := values["loss"][0]

	set_interfaces(ctx, latency, loss, jitter, bandwidth)
}

func main() {
	// Create new router & CLI with defaults
	app := cli.NewRouter("Minimal Example", "1.0.0")

	// Endpointsn
	app.Resource("/").Get("get-root", "Main Page",
		// The only response is HTTP 200
		responses.OK().ContentType("text/plain"),
	).Run(main_page)

	app.Resource("/set").Post("set-values", "Set the interface values",
		responses.NoContent(),
	).Run(set)

	app.Resource("/test").Get("sup test", "Get a short text message",
		// The only response is HTTP 200
		responses.OK().ContentType("text/plain"),
	).Run(test)

	// Start Server
	app.Run()
}

// out, err := exec.Command("date").Output()
// if err != nil {
// 	log.Fatal(err)
// }
// resp := fmt.Sprintf("Hello, world!\nThe date is: %s", out)

// ctx.Write([]byte(resp))
