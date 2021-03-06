package main

import "fmt"
import "io"
import "log"
import "strings"
import "reflect"
import "os"
import "math"

// The 'Raw' field is omitted here, since all of the data is already included
var fieldNames = []string{"BenchmarkId", "BenchmarkDate", "ArgHost", "ArgPort", "ArgURL", "ArgNumConnections", "ArgConnectionRate", "ArgRequestsPerConnection", "ArgDuration", "ConnectionBurstLength", "TotalConnections", "TotalRequests", "TotalReplies", "TestDuration", "ConnectionsPerSecond", "MsPerConnection", "ConcurrentConnections", "ConnectionTimeMin", "ConnectionTimeAvg", "ConnectionTimeMax", "ConnectionTimeMedian", "ConnectionTimeStddev", "ConnectionTimeConnect", "RepliesPerConnection", "RequestsPerSecond", "MsPerRequest", "RequestSize", "RepliesPerSecMin", "RepliesPerSecAvg", "RepliesPerSecMax", "RepliesPerSecStddev", "RepliesPerSecNumSamples", "ReplyTimeResponse", "ReplyTimeTransfer", "ReplySizeHeader", "ReplySizeContent", "ReplySizeFooter", "ReplySizeTotal", "ReplyStatus_1xx", "ReplyStatus_2xx", "ReplyStatus_3xx", "ReplyStatus_4xx", "ReplyStatus_5xx", "CpuTimeUser", "CpuTimeSystem", "CpuPercUser", "CpuPercSystem", "CpuPercTotal", "NetIOValue", "NetIOUnit", "NetIOBytesPerSecond", "ErrTotal", "ErrClientTimeout", "ErrSocketTimeout", "ErrConnectionRefused", "ErrConnectionReset", "ErrFdUnavail", "ErrAddRunAvail", "ErrFtabFull", "ErrOther"}

// Write a CSV header to the given writer including each of the field names
// above, and an optional list of additional column names specified. In the
// resulting file, the optional columns are listed first.
func WriteTSVHeader(w io.Writer) {
	numColumns := len(fieldNames)
	columns := make([]string, 0, numColumns)

	columns = append(columns, fieldNames...)

	io.WriteString(w, strings.Join(columns, ","))
	io.WriteString(w, "\n")
}

func WriteTSVParseDataSet(w io.Writer, data []*PerfData) {
	for _, result := range data {
		WriteTSVParseData(w, result)
	}
}

func WriteTSVParseData(w io.Writer, data *PerfData) {
	numColumns := len(fieldNames)
	columns := make([]string, 0, numColumns)

	// Turn the struct into a Type so we can use reflection
	ptr := reflect.ValueOf(data)
	kind := ptr.Kind()
	if kind != reflect.Ptr {
		log.Fatalf("Could not convert results into a pointer value")
		return
	}

	val := reflect.Indirect(ptr)
	kind = val.Kind();
	if kind != reflect.Struct {
		log.Fatalf("Failed when reflecting on struct")
		return
	}

	// Move through every field, fetching the value by name and adding
	// it to the columns slice

	for _, field := range fieldNames {
		column := val.FieldByName(field)
		if !column.IsValid() {
			log.Fatalf("Failed when reflecting field %s", field)
		}

		t := column.Kind()
		switch t {
		case reflect.String:
			columns = append(columns, column.String())
		case reflect.Float64:
			columns = append(columns, fmt.Sprintf("%#v", column.Float()))
		case reflect.Int:
			columns = append(columns, fmt.Sprintf("%#v", column.Int()))
		case reflect.Int64:
			columns = append(columns, fmt.Sprintf("%#v", column.Int()))
		default:
			log.Println("Type: ", t.String())
			log.Fatalf("Got a field we cannot handle: %s", field)
		}
	}

	io.WriteString(w, strings.Join(columns, ","))
	io.WriteString(w, "\n")
}

func SetHasErrors(perfdata []*PerfData, threshold int) bool {
	total := 0
	for _, data := range perfdata {
		total = total + int(data.ErrTotal)
	}

	if total >= threshold {
		return true
	}

	return false
}

func HasClientErrors(perfdata []*PerfData) bool {
	total := 0
	for _, data := range perfdata {
		total += int(data.ErrFdUnavail)
		total += int(data.ErrAddRunAvail)
		total += int(data.ErrFtabFull)
		total += int(data.ErrOther)
	}

	return total > 0
}

var iFieldNames = []string{"ConnectionBurstLength", "TotalConnections", "TotalRequests", "TotalReplies", "TestDuration", "ConnectionsPerSecond", "MsPerConnection", "ConcurrentConnections", "ConnectionTimeMin", "ConnectionTimeAvg", "ConnectionTimeMax", "ConnectionTimeMedian", "ConnectionTimeStddev", "ConnectionTimeConnect", "RepliesPerConnection", "RequestsPerSecond", "MsPerRequest", "RequestSize", "RepliesPerSecMin", "RepliesPerSecAvg", "RepliesPerSecMax", "RepliesPerSecStddev", "RepliesPerSecNumSamples", "ReplyTimeResponse", "ReplyTimeTransfer", "ReplySizeHeader", "ReplySizeContent", "ReplySizeFooter", "ReplySizeTotal", "ReplyStatus_1xx", "ReplyStatus_2xx", "ReplyStatus_3xx", "ReplyStatus_4xx", "ReplyStatus_5xx", "CpuTimeUser", "CpuTimeSystem", "CpuPercUser", "CpuPercSystem", "CpuPercTotal", "NetIOValue", "ErrTotal", "ErrClientTimeout", "ErrSocketTimeout", "ErrConnectionRefused", "ErrConnectionReset", "ErrFdUnavail", "ErrAddRunAvail", "ErrFtabFull", "ErrOther"}
var iType = []string{"max", "sum", "sum", "sum", "max", "sum", "avg", "sum", "min", "avg", "max", "avg", "avg", "avg", "avg", "sum", "avg", "avg", "min", "avg", "max", "avg", "sum", "avg", "avg", "avg", "avg", "avg", "avg", "sum", "sum", "sum", "sum", "sum", "avg", "avg", "avg", "avg", "avg", "avg", "sum", "sum", "sum", "sum", "sum", "sum", "sum", "sum", "sum"}

func PrintAggregateStats(perfdata []*PerfData, workers int) {
	var res []float64 = make([]float64, 49)
	for i, n := range perfdata[0].All {
		res[i] = n
	}
	
	for _, data := range perfdata[1:] {
		for i, n := range data.All {
			switch iType[i] {
				case "sum", "avg": {
					res[i] += n
				}
				case "min": {
					res[i] = math.Min(res[i], n)
				}
				case "max": {
					res[i] = math.Max(res[i], n)
				}
			}
			
		}
	}

	for i, t := range iType {
		if (t == "avg") {
			res[i] = res[i]/float64(workers)
		}
	}

	sres := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(res)), "|"), "[]") + "\n"
	fmt.Println(sres)

	// Write the result to the file
	file, err := os.OpenFile("results.csv", os.O_RDWR|os.O_APPEND, 0666);
	if err != nil {
		log.Println("Writing results error:", err)
	}
	_, err = file.WriteString(sres)
	if err != nil {
		log.Println("Writing results error:", err)
	}
	file.Close()
}
