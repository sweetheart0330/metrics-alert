package app

import "flag"

type options struct {
	Address        string
	ReportInterval int
	PollInterval   int
}

func getAgentOptions() (op options) {
	flag.StringVar(&op.Address, "a", "localhost:8080", "address and port to send requests")
	flag.IntVar(&op.ReportInterval, "r", 10, "interval between sending requests")
	flag.IntVar(&op.PollInterval, "p", 2, "interval between collecting metrics")
	flag.Parse()

	return op
}

func getServerAddress() (fl string) {
	flag.StringVar(&fl, "a", "localhost:8080", "address and port to run server")
	flag.Parse()

	return fl
}
