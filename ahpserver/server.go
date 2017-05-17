package main

import "flag"
import "fmt"
import "os/exec"
import "net/http"
import "io/ioutil"
import "log"
import "net"
import "net/rpc"
import "errors"

type Args struct {
	Host                  string
	Port                  int
	URL                   string
	NumConnections        int
	ConnectionRate        int
	RequestsPerConnection int
	Duration              int
}

type Result struct {
	Stdout     string
	Stderr     string
	ExitStatus int
}

type HTTPerf int

const (
	ERR_EXECNOTFOUND = "Could not find the 'httperf' executable: %s"
	ERR_RUNFAILED    = "Failed to run command: %s"
	ERR_WAIT         = "Failed when waiting on pid %d"
	ERR_NOTEXITED    = "Command did not properly exit: %s"
	ERR_READOUT      = "Could not read stdout: %s"
	ERR_READERR      = "Could not read stderr: %s"
)

func (h *HTTPerf) Benchmark(args *Args, result *Result) error {
	// Try to find the 'httpperf' command, which must exist in the PATH
	// of the current user/environment.

	perfexec, err := exec.LookPath("httperf")
	if err != nil {
		return errors.New(fmt.Sprintf(ERR_EXECNOTFOUND, err.Error()))
	}

	// Build the httperf commandline and build a result to return
	argv := []string{
		"--server", args.Host,
		"--port", fmt.Sprintf("%d", args.Port),
		"--uri", args.URL,
		"--num-conns", fmt.Sprintf("%d", args.NumConnections),
		"--rate", fmt.Sprintf("%d", args.ConnectionRate),
		"--num-calls", fmt.Sprintf("%d", args.RequestsPerConnection),
		"--hog",
	}

	log.Printf("++ [%p] Running benchmark of %s on port %d", args, args.Host, args.Port)
	log.Printf("   [%p] Input arguments: %#v", args, args)
	log.Printf("   [%p] Commandline arguments: %#v", args, argv)

	cmd := exec.Command(perfexec, argv...)
	outpipe, err := cmd.StdoutPipe()
	if err != nil {
		return errors.New(fmt.Sprintf(ERR_READOUT, err.Error()))
	}
	errpipe, err := cmd.StderrPipe()
	if err != nil {
		return errors.New(fmt.Sprintf(ERR_READOUT, err.Error()))
	}

	err = cmd.Start()
	if err != nil {
		return errors.New(fmt.Sprintf(ERR_RUNFAILED, err.Error()))
	}

	log.Printf("   [%p] Process successfully started with PID: %d", args, cmd.Process.Pid)

	output, err := ioutil.ReadAll(outpipe)
	log.Println("Finished reading stdout.")
	if err != nil {
		return errors.New(fmt.Sprintf(ERR_READOUT, err.Error()))
	}
	errout, err := ioutil.ReadAll(errpipe)
	log.Println("Finished reading stderr.")
	if err != nil {
		return errors.New(fmt.Sprintf(ERR_READERR, err.Error()))
	}

	log.Printf("   [%p] Finished reading stdout and stderr", args)

	err = cmd.Wait()
	log.Printf("-- [%p] Command joined and finished", args)

	if err != nil {
		log.Println("Error:", err)
		return errors.New(fmt.Sprintf(ERR_WAIT, cmd.Process.Pid))
	}

	result.Stdout = string(output)
	result.Stderr = string(errout)

	return nil
}

var host *string = flag.String("host", "", "The host on which to bind the server")
var port *int = flag.Int("port", 1717, "The port on which to bind the server")

func main() {
	flag.Parse()

	httperf := new(HTTPerf)
	rpc.Register(httperf)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if e != nil {
		log.Fatalf("listen error:", e)
	}

	log.Printf("Now listening for requests on %s:%d", *host, *port)
	http.Serve(l, nil)
}
