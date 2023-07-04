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

func set_tc_params(ctx huma.Context, vals TcParams) []byte {
	log.Printf("Setting: %+v", vals)
	args := []string{vals.Latency, vals.Loss, vals.Jitter, vals.Bandwidth}

	// Run script (for now)
	output, err := exec.Command("./static/exectc.sh", args...).Output()
	if err != nil {
		log.Println(err)
	}
	save_current_settings(vals)
	return output
}

func main_page(ctx huma.Context) {
	body, err := ioutil.ReadFile("./static/index.html")

	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	ctx.Header().Set("Content-Type", "text/html")
	ctx.Write(body)
}

func read_page(ctx huma.Context) {
	ctx.Header().Set("Content-Type", "text/html")
	vals, err := read_current_settings()
	if err != nil {
		log.Fatalf("Error reading current settings: %v", err)
	}
	body := fmt.Sprintf("Latency (ms): &emsp; &emsp; %s<br>Jitter (ms): &emsp; &emsp; &emsp; %s<br>Bandwidth (kbit/s): %s<br>Packet Loss (%%): &emsp; %s<br>", vals.Latency, vals.Jitter, vals.Bandwidth, vals.Loss)
	ctx.Write([]byte(body))
}

func set_page(ctx huma.Context, data FormData) {
	ctx.Header().Set("Content-Type", "text/html")

	// Pull Form Values out into something better
	rawbody := data.Body
	rawparams := new(strings.Builder)
	_, _ = io.Copy(rawparams, rawbody)
	// log.Println(rawparams.String())
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

	// Set TC Params
	out := set_tc_params(ctx, vals)
	ctx.Write(out)

}

func load_defaults() {
	_, err := read_current_settings()
	if err != nil {
		default_vals := TcParams{"0", "0", "0", "1000000"}
		save_current_settings(default_vals)
	}
}

func main() {
	// Create current_settings.json if needed
	load_defaults()

	// Create new Huma router & CLI with defaults
	app := cli.NewRouter("RPi WAN Emulation", "1.0.0")

	// Endpoints:

	// Main
	app.Resource("/").Get("get-root", "Main Page",
		// The only response is HTTP 200
		responses.OK().ContentType("text/html"),
	).Run(main_page)

	// Read
	app.Resource("/read").Get("read-values", "Get the existing TC values",
		responses.OK().ContentType("text/html"),
	).Run(read_page)

	// Set
	app.Resource("/set").Post("set-values", "Set the TC values",
		responses.OK().ContentType("text/html"),
	).Run(set_page)

	// Start Server
	app.Run()
}
