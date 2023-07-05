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

var Version string = "0.0.1"

type FormData struct {
	Body io.Reader
}

type TcParams struct {
	Latency   string
	Loss      string
	Jitter    string
	Bandwidth string
	Eth1      string
	Eth2      string
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
	args := []string{vals.Latency, vals.Loss, vals.Jitter, vals.Bandwidth, vals.Eth1, vals.Eth2}

	// Run TC script (for now)
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
	vals, err := read_current_settings()
	if err != nil {
		log.Fatalf("Error reading current settings: %v", err)
	}
	body := fmt.Sprintf("Bridge Port 1: &emsp; &emsp; %s<br>Bridge Port 2: &emsp; &emsp; %s<br><br>Latency (ms): &emsp; &emsp; %s<br>Jitter (ms): &emsp; &emsp; &emsp; %s<br>Bandwidth (kbit/s): %s<br>Packet Loss (%%): &emsp; %s<br>", vals.Eth1, vals.Eth2, vals.Latency, vals.Jitter, vals.Bandwidth, vals.Loss)

	ctx.Header().Set("Content-Type", "text/html")
	ctx.Write([]byte(body))
}

func set_page(ctx huma.Context, data FormData) {
	// Pull Form Values out into TcParams struct
	rawbody := data.Body
	rawparams := new(strings.Builder)
	_, _ = io.Copy(rawparams, rawbody)
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
	vals.Eth1 = values["eth1"][0]
	vals.Eth2 = values["eth2"][0]

	// Set TC Params
	body := set_tc_params(ctx, vals)
	ctx.Header().Set("Content-Type", "text/html")
	ctx.Write(body)

}

func load_defaults() {
	_, err := read_current_settings()
	if err != nil {
		default_vals := TcParams{"0", "0", "0", "1000000", "eth1", "eth2"}
		save_current_settings(default_vals)
	}
}

func main() {
	// Create current_settings.json if needed
	load_defaults()

	// Create new Huma router & CLI with defaults
	app := cli.NewRouter("RPi WAN Emulation", Version)

	// Endpoints:
	app.Resource("/").Get("get-root", "Main Page",
		responses.OK().ContentType("text/html"),
	).Run(main_page)

	app.Resource("/read").Get("read-values", "Get the existing TC values",
		responses.OK().ContentType("text/html"),
	).Run(read_page)

	app.Resource("/set").Post("set-values", "Set the TC values",
		responses.OK().ContentType("text/html"),
	).Run(set_page)

	app.Run()
}
