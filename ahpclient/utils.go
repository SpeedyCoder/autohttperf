package main

import "fmt"
import "io"
import "log"
import "strings"
import "reflect"
import "os"

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

func PrintAggregateStats(perfdata []*PerfData, workers int) {
	var connRate float64 = 0
	var reqRate float64 = 0

	var totalReq float64 = 0
	var totalRep float64 = 0
	var errs float64 = 0

	var repliesPerSec float64 = 0
	var concurr float64 = 0
	var connTimeAvg float64 = 0
	var replyTime float64 = 0
	
	for _, data := range perfdata {
		connRate += data.ConnectionsPerSecond
		reqRate += data.RequestsPerSecond
		totalReq += data.TotalRequests
		totalRep += data.TotalReplies
		errs += data.ErrTotal

		connTimeAvg += data.ConnectionTimeAvg
		concurr += data.ConcurrentConnections
		repliesPerSec += data.RepliesPerSecAvg
		replyTime += data.ReplyTimeResponse
	}

	fmt.Printf("Connection rate: %.2f\n", connRate)
	fmt.Printf("Request rate: %.2f\n\n", reqRate)

	fmt.Println("\nTotalRequests:", int64(totalReq))
	fmt.Println("TotalReplies:", int64(totalRep))
	fmt.Println("Errors:", int64(errs))
	fmt.Printf("Success rate: %.2f\n\n", totalRep/totalReq)

	fmt.Printf("ConnectionTimeAvg: %.2f\n", connTimeAvg/float64(workers))
	fmt.Println("ConcurrentConnections:", int(concurr))
	fmt.Printf("RepliesPerSecAvg: %.2f\n", repliesPerSec/float64(workers))
	fmt.Printf("ReplyTimeResponse: %.2f\n\n", replyTime/float64(workers))

	res := fmt.Sprintf("%.2f|%.2f|%v|%v|%v|%.2f|%.2f|%v|%.2f|%.2f\n", 
		connRate, reqRate,
		int64(totalReq), int64(totalRep), int64(errs),
		totalRep/totalReq, 
		connTimeAvg/float64(workers), int(concurr), 
		repliesPerSec/float64(workers),
		replyTime/float64(workers))
	fmt.Println(res)

	// Write the result to the file
	file, err := os.OpenFile("results.txt", os.O_RDWR|os.O_APPEND, 0666);
	if err != nil {
		log.Println(err)
	}
	_, err = file.WriteString(res)
	if err != nil {
		log.Println(err)
	}
	file.Close()
}
