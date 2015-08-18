package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type IncCompany struct {
	City      string  `json:"city"`
	Company   string  `json:"company"`
	Growth    float64 `json:"growth"`
	ID        int     `json:"id"`
	Industry  string  `json:"industry"`
	Metro     string  `json:"metro"`
	Rank      int     `json:"rank"`
	Revenue   int     `json:"revenue"`
	StateL    string  `json:"state_l"`
	StateS    string  `json:"state_s"`
	URL       string  `json:"url"`
	Workers   int     `json:"workers"`
	YrsOnList int     `json:"yrs_on_list"`
}

func downloadList(year int) []IncCompany {
	url := fmt.Sprintf("http://www.inc.com/inc5000list/json/inc5000_%d.json", year)

	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err.Error())
	}

	var list []IncCompany
	json.Unmarshal(body, &list)

	return list
}

func writeListToCSV(year int, list []IncCompany) {
	f, err := os.Create(fmt.Sprintf("./inc5000_%d_output.csv", year))
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	w := csv.NewWriter(f)

	headers := []string{"Rank", "Company", "Growth", "Revenue", "Industry", "State", "Metro", "Employees", "URL", "YrsOnList"}
	w.Write(headers)

	for _, obj := range list {
		var record []string
		record = append(record, fmt.Sprintf("%v", obj.Rank))
		record = append(record, obj.Company)
		record = append(record, fmt.Sprintf("%v%%", obj.Growth))
		record = append(record, fmt.Sprintf("$%v", obj.Revenue))
		record = append(record, obj.Industry)
		record = append(record, obj.StateL)
		record = append(record, obj.Metro)
		record = append(record, fmt.Sprintf("%v", obj.Workers))
		record = append(record, fmt.Sprintf("http://www.inc.com/profile/%s", obj.URL))
		record = append(record, fmt.Sprintf("%v", obj.YrsOnList))
		if obj.Company != "" {
			w.Write(record)
		}
	}
	w.Flush()
}

func main() {
	var year int
	flag.IntVar(&year, "year", 2015, "What year to pull for Inc5000 list")
	flag.Parse()

	fmt.Println("Pulling list for year", year)

	list := downloadList(year)
	writeListToCSV(year, list)
}
