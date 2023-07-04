package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/danielgtaylor/huma"
	"github.com/danielgtaylor/huma/cli"
	"github.com/danielgtaylor/huma/responses"
)

type FormData struct {
	Body io.Reader
}

type TcParams struct {
	Latency   string
	Loss      string
	Jitter    string
	Bandwidth string
}

func save_current_settings(vals TcParams) {
	file, _ := json.MarshalIndent(vals, "", "  ")
	_ = ioutil.WriteFile("./current_settings.json", file, 0644)
}

func read_current_settings() (TcParams, error) {
	var vals TcParams
	settingsFile, err := os.Open("./current_settings.json")
	if err != nil {
		return vals, err
	}
	byteValue, _ := ioutil.ReadAll(settingsFile)
	json.Unmarshal(byteValue, &vals)
	return vals, err
}

func set_interfaces(ctx huma.Context, vals TcParams) {
	ctx.Header().Set("Content-Type", "text/plain")
	log.Println(vals.Latency, vals.Loss, vals.Jitter, vals.Bandwidth)

	args := fmt.Sprintf("%s %s %s %s", vals.Latency, vals.Loss, vals.Jitter, vals.Bandwidth)
	out, err := exec.Command("./static/tccommands.sh", args).Output()
	if err != nil {
		log.Println(err)
	}
	resp := fmt.Sprintf("%s", out)
	save_current_settings(vals)

	ctx.Write([]byte(resp))
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

func read(ctx huma.Context) {
	ctx.Header().Set("Content-Type", "text/html")
	vals, err := read_current_settings()
	if err != nil {
		log.Fatalf("Error reading current settings: %v", err)
	}
	body := fmt.Sprintf("Latency (ms): &emsp; &emsp; %s<br>Jitter (ms): &emsp; &emsp; &emsp; %s<br>Bandwidth (kbit/s): %s<br>Packet Loss (%%): &emsp; %s<br>", vals.Latency, vals.Jitter, vals.Bandwidth, vals.Loss)
	ctx.Write([]byte(body))
}

func set(ctx huma.Context, data FormData) {
	ctx.Header().Set("Content-Type", "text/plain")
	//body := fmt.Sprintf("latency: %d\njitter: %d\nbandwidth: %d\nloss: %d", data.Body.Latency, data.Body.Jitter, data.Body.Bandwidth, data.Body.Loss)
	rawbody := data.Body
	rawparams := new(strings.Builder)
	_, _ = io.Copy(rawparams, rawbody)
	log.Println(rawparams.String())

	myurl := fmt.Sprintf("https://x.com/?%s", rawparams.String())
	params, err := url.Parse(myurl)
	if err != nil {
		log.Fatal(err)
	}
	values := params.Query()
	var vals TcParams
	vals.Latency = values["latency"][0]
	vals.Loss = values["loss"][0]
	vals.Jitter = values["jitter"][0]
	vals.Bandwidth = values["bandwidth"][0]

	log.Printf("Setting: %+v", vals)
	set_interfaces(ctx, vals)
}

func load_defaults() {
	_, err := read_current_settings()
	if err != nil {
		default_vals := TcParams{"0", "0", "0", "0"}
		save_current_settings(default_vals)
	}
}

func main() {
	// Create current_settings.json
	load_defaults()
	// Create new router & CLI with defaults
	app := cli.NewRouter("Minimal Example", "1.0.0")

	// Endpointsn
	app.Resource("/").Get("get-root", "Main Page",
		// The only response is HTTP 200
		responses.OK().ContentType("text/plain"),
	).Run(main_page)

	app.Resource("/read").Get("read-values", "Get the existing latency/loss/jitter/bandwidth values",
		responses.OK().ContentType("text/plain"),
	).Run(read)

	app.Resource("/set").Post("set-values", "Set the latency/loss/jitte/bandwidth",
		responses.NoContent(),
	).Run(set)

	// app.Resource("/test").Get("sup test", "Get a short text message",
	// 	// The only response is HTTP 200
	// 	responses.OK().ContentType("text/plain"),
	// ).Run(test)

	// Start Server
	app.Run()
}

// out, err := exec.Command("date").Output()
// if err != nil {
// 	log.Fatal(err)
// }
// resp := fmt.Sprintf("Hello, world!\nThe date is: %s", out)

// ctx.Write([]byte(resp))
