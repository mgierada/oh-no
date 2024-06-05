package utils

type Table struct {
	Counter               string
	HistoricalCounter     string
	OhnoCounter           string
	HistoricalOhnoCounter string
}

func getTable() Table {
	return Table{
		Counter:               "counter",
		HistoricalCounter:     "historical_counter",
		OhnoCounter:           "ohno_counter",
		HistoricalOhnoCounter: "historical_ohno_counter",
	}
}

var TableInstance = getTable()
