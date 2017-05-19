package main

import "fmt"
import "regexp"
import "strconv"
import "errors"

var resultPattern = `Maximum connect burst length: ([0-9]*)

Total: connections ([0-9]*) requests ([0-9]*) replies ([0-9]*) test-duration ([0-9]*\.?[0-9]*) s

Connection rate: ([0-9]*\.?[0-9]*) conn/s \(([0-9]*\.?[0-9]*) ms/conn, <=([0-9]*) concurrent connections\)
Connection time \[ms\]: min ([0-9]*\.?[0-9]*) avg ([0-9]*\.?[0-9]*) max ([0-9]*\.?[0-9]*) median ([0-9]*\.?[0-9]*) stddev ([0-9]*\.?[0-9]*)
Connection time \[ms\]: connect ([0-9]*\.?[0-9]*)
Connection length \[replies/conn\]: ([0-9]*\.?[0-9]*)

Request rate: ([0-9]*\.?[0-9]*) req/s \(([0-9]*\.?[0-9]*) ms/req\)
Request size \[B\]: ([0-9]*\.?[0-9]*)

Reply rate \[replies/s\]: min ([0-9]*\.?[0-9]*) avg ([0-9]*\.?[0-9]*) max ([0-9]*\.?[0-9]*) stddev ([0-9]*\.?[0-9]*) \(([0-9])* samples\)
Reply time \[ms\]: response ([0-9]*\.?[0-9]*) transfer ([0-9]*\.?[0-9]*)
Reply size \[B\]: header ([0-9]*\.?[0-9]*) content ([0-9]*\.?[0-9]*) footer ([0-9]*\.?[0-9]*) \(total ([0-9]*\.?[0-9]*)\)
Reply status: 1xx=([0-9]*) 2xx=([0-9]*) 3xx=([0-9]*) 4xx=([0-9]*) 5xx=([0-9]*)

CPU time \[s\]: user ([0-9]*\.?[0-9]*) system ([0-9]*\.?[0-9]*) \(user ([0-9]*\.?[0-9]*)\% system ([0-9]*\.?[0-9]*)\% total ([0-9]*\.?[0-9]*)\%\)
Net I/O: ([0-9]*\.?[0-9]*) (.*) \((.*) bps\)

Errors: total ([0-9]*) client-timo ([0-9]*) socket-timo ([0-9]*) connrefused ([0-9]*) connreset ([0-9]*)
Errors: fd-unavail ([0-9]*) addrunavail ([0-9]*) ftab-full ([0-9]*) other ([0-9]*)`

var resultRegexp = regexp.MustCompile(resultPattern)

const NUM_RESULTS = 52

func ParseResultsRaw(str string) []string {
	return resultRegexp.FindStringSubmatch(str)
}

func ParseResults(str string, id string, date int64, args *Args) (*PerfData, error) {
	results := ParseResultsRaw(str)
	data := new(PerfData)

	data.BenchmarkId = id
	data.BenchmarkDate = date
	data.ArgHost = args.Host
	data.ArgPort = args.Port
	data.ArgURL = args.URL
	data.ArgNumConnections = args.NumConnections
	data.ArgConnectionRate = args.ConnectionRate
	data.ArgRequestsPerConnection = args.RequestsPerConnection
	data.ArgDuration = args.Duration

	var conv float64
	var err error
	data.All = make([]float64, 49)
	i := 0

	data.Raw = results[0]
	if conv, err = strconv.ParseFloat(results[1], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 1, err.Error()))
	}
	data.ConnectionBurstLength = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[2], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 2, err.Error()))
	}
	data.TotalConnections = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[3], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 3, err.Error()))
	}
	data.TotalRequests = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[4], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 4, err.Error()))
	}
	data.TotalReplies = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[5], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 5, err.Error()))
	}
	data.TestDuration = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[6], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 6, err.Error()))
	}
	data.ConnectionsPerSecond = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[7], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 7, err.Error()))
	}
	data.MsPerConnection = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[8], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 8, err.Error()))
	}
	data.ConcurrentConnections = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[9], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 9, err.Error()))
	}
	data.ConnectionTimeMin = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[10], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 10, err.Error()))
	}
	data.ConnectionTimeAvg = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[11], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 11, err.Error()))
	}
	data.ConnectionTimeMax = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[12], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 12, err.Error()))
	}
	data.ConnectionTimeMedian = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[13], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 13, err.Error()))
	}
	data.ConnectionTimeStddev = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[14], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 14, err.Error()))
	}
	data.ConnectionTimeConnect = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[15], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 15, err.Error()))
	}
	data.RepliesPerConnection = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[16], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 16, err.Error()))
	}
	data.RequestsPerSecond = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[17], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 17, err.Error()))
	}
	data.MsPerRequest = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[18], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 18, err.Error()))
	}
	data.RequestSize = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[19], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 19, err.Error()))
	}
	data.RepliesPerSecMin = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[20], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 20, err.Error()))
	}
	data.RepliesPerSecAvg = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[21], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 21, err.Error()))
	}
	data.RepliesPerSecMax = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[22], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 22, err.Error()))
	}
	data.RepliesPerSecStddev = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[23], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 23, err.Error()))
	}
	data.RepliesPerSecNumSamples = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[24], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 24, err.Error()))
	}
	data.ReplyTimeResponse = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[25], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 25, err.Error()))
	}
	data.ReplyTimeTransfer = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[26], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 26, err.Error()))
	}
	data.ReplySizeHeader = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[27], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 27, err.Error()))
	}
	data.ReplySizeContent = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[28], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 28, err.Error()))
	}
	data.ReplySizeFooter = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[29], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 29, err.Error()))
	}
	data.ReplySizeTotal = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[30], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 30, err.Error()))
	}
	data.ReplyStatus_1xx = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[31], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 31, err.Error()))
	}
	data.ReplyStatus_2xx = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[32], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 32, err.Error()))
	}
	data.ReplyStatus_3xx = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[33], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 33, err.Error()))
	}
	data.ReplyStatus_4xx = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[34], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 34, err.Error()))
	}
	data.ReplyStatus_5xx = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[35], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 35, err.Error()))
	}
	data.CpuTimeUser = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[36], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 36, err.Error()))
	}
	data.CpuTimeSystem = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[37], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 37, err.Error()))
	}
	data.CpuPercUser = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[38], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 38, err.Error()))
	}
	data.CpuPercSystem = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[39], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 39, err.Error()))
	}
	data.CpuPercTotal = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[40], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 40, err.Error()))
	}
	data.NetIOValue = conv
	data.All[i] = conv; i++
	data.NetIOUnit = results[41]
	data.NetIOBytesPerSecond = results[42]
	if conv, err = strconv.ParseFloat(results[43], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 43, err.Error()))
	}
	data.ErrTotal = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[44], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 44, err.Error()))
	}
	data.ErrClientTimeout = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[45], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 45, err.Error()))
	}
	data.ErrSocketTimeout = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[46], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 46, err.Error()))
	}
	data.ErrConnectionRefused = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[47], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 47, err.Error()))
	}
	data.ErrConnectionReset = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[48], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 48, err.Error()))
	}
	data.ErrFdUnavail = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[49], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 49, err.Error()))
	}
	data.ErrAddRunAvail = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[50], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 50, err.Error()))
	}
	data.ErrFtabFull = conv
	data.All[i] = conv; i++
	if conv, err = strconv.ParseFloat(results[51], 64); err != nil {
		return nil, errors.New(fmt.Sprintf("Error parsing field %d:%s", 51, err.Error()))
	}
	data.ErrOther = conv
	data.All[i] = conv; i++

	return data, nil
}

/* The following Lua script was used to generate the above code:

local fields = {"Raw", "ConnectionBurstLength", "TotalConnections", "TotalRequests", "TotalReplies", "TestDuration", "ConnectionsPerSecond", "MsPerConnection", "ConcurrentConnections", "ConnectionTimeMin", "ConnectionTimeAvg", "ConnectionTimeMax", "ConnectionTimeMedian", "ConnectionTimeStddev", "ConnectionTimeConnect", "RepliesPerConnection", "RequestsPerSecond", "MsPerRequest", "RequestSize", "RepliesPerSecMin", "RepliesPerSecAvg", "RepliesPerSecMax", "RepliesPerSecStddev", "RepliesPerSecNumSamples", "ReplyTimeResponse", "ReplyTimeTransfer", "ReplySizeHeader", "ReplySizeContent", "ReplySizeFooter", "ReplySizeTotal", "ReplyStatus_1xx", "ReplyStatus_2xx", "ReplyStatus_3xx", "ReplyStatus_4xx", "ReplyStatus_5xx",
"CpuTimeUser", "CpuTimeSystem", "CpuPercUser", "CpuPercSystem", "CpuPercTotal", "NetIOValue", "NetIOUnit", "NetIOBytesPerSecond", "ErrTotal", "ErrClientTimeout", "ErrSocketTimeout", "ErrConnectionRefused", "ErrConnectionReset", "ErrFdUnavail", "ErrAddRunAvail", "ErrFtabFull", "ErrOther"}

local strings = {Raw = true, NetIOUnit = true, NetIOBytesPerSecond = true}

for idx, field in ipairs(fields) do
	if strings[field] then
		print(string.format("data.%s = results[%d]", field, idx - 1))
	else
		print(string.format("if conv, err = strconv.ParseFloat(results[%d], 64); err != nil {", idx - 1))
		print(string.format([[	return nil, errors.New(fmt.Sprintf("Error parsing field %%d:%%s", %d, err.Error()))]], idx - 1))
		print(string.format("}"))
		print(string.format("data.%s = conv
		data.All[i] = conv; i++", field))
	end
end

*/
