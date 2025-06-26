package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type FlightNumber struct {
	name string
	id   int
}

type Date struct {
	day, month, year int
}

type Time struct {
	minutes, hours int
}

type Record struct {
	flight_number  FlightNumber
	date           Date
	departure_time Time
	arrival_time   Time
	free_count     int
	status         string
	ticket_cost    int
}

func (record Record) String() string {
	var to_return string
	to_return += RJust(record.flight_number.name+ZeroLJusr(strconv.Itoa(record.flight_number.id), 3), 13, false)
	to_return += " | "
	to_return += RJust(ZeroLJusr(strconv.Itoa(record.date.day), 2)+"."+ZeroLJusr(strconv.Itoa(record.date.month), 2)+"."+strconv.Itoa(record.date.year), 11, false)
	to_return += " | "
	delta_minutes, delta_hours := eval_time_delta(record.departure_time, record.arrival_time)
	to_return += RJust(ZeroLJusr(strconv.Itoa(delta_hours), 2)+":"+ZeroLJusr(strconv.Itoa(delta_minutes), 2), 15, false)
	to_return += " | "
	to_return += RJust(strconv.Itoa(record.free_count), 22, false)
	to_return += " | "
	to_return += RJust(eval_day_of_week(record.date), 11, false)
	to_return += " | "
	to_return += RJust(record.status, 13, true)
	to_return += " | "
	to_return += RJust(strconv.Itoa(record.ticket_cost*record.free_count), 11, false)
	to_return += "\n"
	return to_return
}

func ZeroLJusr(data string, _len int) string {
	var new_data string
	if len(data) == _len {
		return data
	} else {
		start_len := len(data)
		for i := 0; i < _len-start_len; i++ {
			new_data += "0"
		}
		new_data += data
		return new_data
	}
}

func RJust(data string, _len int, is_russion bool) string {
	start_len := len(data)
	if is_russion {
		start_len /= 2
	}
	if _len <= start_len {
		return data
	}
	for i := 0; i < _len-start_len; i++ {
		data += " "
	}
	return data
}

func FindElem(data []Record, elem Record) int {
	for i := range data {
		if data[i] == elem {
			return i
		}
	}
	return -1
}

func IsBigger(first_record, second_record Record) bool {
	if (first_record.free_count * first_record.ticket_cost) >=
		(second_record.free_count * second_record.ticket_cost) {
		// fmt.Println("Here 1")
		return true
	} else if (first_record.free_count * first_record.ticket_cost) <
		(second_record.free_count * second_record.ticket_cost) {
		// fmt.Println("Here 2")
		return false
	} else {
		if first_record.flight_number.name > second_record.flight_number.name {
			// fmt.Println("Here 3")
			return true
		} else if first_record.flight_number.name <= second_record.flight_number.name &&
			first_record.flight_number.id > second_record.flight_number.id {
			// fmt.Println("Here 4")
			return true
		} else if first_record.flight_number.name < second_record.flight_number.name &&
			first_record.flight_number.id < second_record.flight_number.id {
			// fmt.Println("Here 5")
			return false
		} else {
			if status_converter(first_record.status) > status_converter(second_record.status) {
				// fmt.Println("Here 6")
				return true
			} else if status_converter(first_record.status) < status_converter(second_record.status) {
				// fmt.Println("Here 7")
				return false
			}
		}
	}
	return false
}

func NewIsBigger(first_record, second_record Record) bool {
	if (first_record.free_count * first_record.ticket_cost) >
		(second_record.free_count * second_record.ticket_cost) {
		// fmt.Println("Here 1")
		return true
	} else if (first_record.free_count * first_record.ticket_cost) <
		(second_record.free_count * second_record.ticket_cost) {
		// fmt.Println("Here 2")
		return false
	} else {
		if first_record.flight_number.name > second_record.flight_number.name {
			// fmt.Println("Here 3")
			return true
		} else if first_record.flight_number.name < second_record.flight_number.name {
			return false
		} else {
			if first_record.flight_number.id > second_record.flight_number.id {
				return true
			} else if first_record.flight_number.id < second_record.flight_number.id {
				return false
			} else {
				if status_converter(first_record.status) > status_converter(second_record.status) {
					return true
				}
			}
		}
	}

	return false
}

func check_split(data, delimeter string) []string {
	result := []string{}
	var to_add string
	for i := 0; i < len(data); i++ {
		if string(data[i]) != delimeter {
			to_add += string(data[i])
		} else {
			if to_add != "" {
				result = append(result, to_add)
				to_add = ""
			}
		}
	}
	if to_add != "" {
		result = append(result, to_add)
	}
	return result
}

func split(data, delimeter string) []string {
	result := make([]string, strings.Count(data, delimeter)+1)
	count := strings.Count(data, delimeter)
	for i := 0; i < count; i++ {
		if strings.Index(data, delimeter) > 0 {
			result[i] = data[:strings.Index(data, delimeter)]
		} else {
			result[i] = ""
		}
		data = data[strings.Index(data, delimeter)+len(delimeter):]
	}
	result[count] = data
	return result
}

func status_converter(data string) int {
	if data == "по расписанию" {
		return 1
	} else if data == "задержка" {
		return 2
	} else if data == "отменён" {
		return 3
	} else {
		return 5
	}
}

func month_converter(month string) int {
	switch strings.ToLower(month) {
	case "jan":
		return 1
	case "feb":
		return 2
	case "mar":
		return 3
	case "apr":
		return 4
	case "may":
		return 5
	case "jun":
		return 6
	case "jul":
		return 7
	case "aug":
		return 8
	case "sep":
		return 9
	case "oct":
		return 10
	case "nov":
		return 11
	case "dec":
		return 12
	default:
		return 0
	}
}

func year_converter(year int) int {
	if year >= 0 && year <= 24 {
		return year + 2000
	} else {
		return year
	}
}

func eval_date(date string) (int, int, int) {
	data := split(date, ".")
	day, _ := strconv.Atoi(data[0])
	month := month_converter(data[1])
	year, _ := strconv.Atoi(data[2])
	year = year_converter(year)
	return day, month, year
}

func eval_time(data string) (int, int) {
	hours, _ := strconv.Atoi(data[:2])
	minutes, _ := strconv.Atoi(data[3:5])
	if strings.Contains("P", data) || strings.Contains("p", data) {
		hours += 12
	}
	return minutes, hours
}

func eval_time_delta(departure_time Time, arrival_time Time) (int, int) {
	var hours, minutes, time_delta int
	if arrival_time.hours < departure_time.hours {
		time_delta = (arrival_time.hours+12)*60 + arrival_time.minutes - departure_time.hours*60 - departure_time.minutes
	} else {
		time_delta = arrival_time.hours*60 + arrival_time.minutes - departure_time.hours*60 - departure_time.minutes
	}
	hours = time_delta / 60
	minutes = time_delta % 60
	return minutes, hours
}

func eval_day_of_week(date Date) string {
	if date.year > 2000 {
		date.year -= 2000
	}
	_sum := int64(0)
	_sum += int64(date.year/4) * int64(366)
	_sum += int64(date.year-date.year/4) * int64(365)
	for i := 1; i < date.month; i++ {
		if i == 2 && date.year%4 == 0 {
			_sum += 29
		} else if i == 2 {
			_sum += 28
		} else if i == 8 || i%2 == 1 {
			_sum += 31
		} else {
			_sum += 30
		}
	}
	_sum += int64(date.day)
	if _sum%7 == 1 {
		return "Sat"
	} else if _sum%7 == 2 {
		return "Sun"
	} else if _sum%7 == 3 {
		return "Mon"
	} else if _sum%7 == 4 {
		return "Tue"
	} else if _sum%7 == 5 {
		return "Wed"
	} else if _sum%7 == 6 {
		return "Thu"
	} else {
		return "Fri"
	}
}

func eval_cost(data string) int {
	if strings.Contains(data, "\r") == true {
		data = data[:len(data)-1]
	}
	cost, _ := strconv.Atoi(data)
	return cost
}

func check_access(record Record, data []Record) bool {
	for i := 0; i < len(data); i++ {
		if record.flight_number == data[i].flight_number {
			if record.date == data[i].date && record.departure_time.hours <= data[i].departure_time.hours {
				return false
			}
		}
	}
	return true
}

func check_in_output_data(record Record, output_data []Record) bool {
	for i := 0; i < len(output_data); i++ {
		if record == output_data[i] {
			return true
		}
	}
	return false
}

func check_status(data string) string {
	data = strings.ToLower(data)
	for i := 0; i < len(data); i++ {
		if 1103 <= int(data[i]) && int(data[i]) <= 1072 && int(data[i]) != 32 {
			return "incorect"
		}
	}
	if data != "отменён" && data != "задержка" && data != "по расписанию" {
		fmt.Println("Wrong status")
		return "abnormal"
	} else {
		return "output"
	}
}

func check_count_places(data string) string {
	count, err := strconv.Atoi(data)
	if err != nil {
		fmt.Println("Wrong type of free places count")
		return "incorect"
	} else {
		if count >= 0 && count <= 900 {
			return "output"
		} else {
			return "abnormal"
		}
	}
}

func check_time(time string) string { // remake with documentation - done
	if len(time) < 8 {
		fmt.Println("Not time")
		return "incorect"
	}
	hours, err_h := strconv.Atoi(time[:2])
	minutes, err_m := strconv.Atoi(time[3:5])
	if err_h != nil || err_m != nil || (strings.Contains(time, "p") == false &&
		strings.Contains(time, "P") == false && strings.Contains(time, "a") == false &&
		strings.Contains(time, "A") == false) {
		fmt.Println("Wrong type of time ")
		return "incorect"
	}
	if hours < 0 || hours > 24 {
		fmt.Println("Wrong hours")
		return "abnormal"
	}
	if minutes < 0 || minutes > 60 {
		fmt.Println("Wrong minutes")
		return "abnormal"
	}
	return "output"
}

func check_date(date string) string { // remake with documentation - done
	data := split(date, ".")
	if len(data) < 3 {
		fmt.Println("Now date")
		return "incorect"
	}
	day, err_d := strconv.Atoi(data[0])
	month := month_converter(data[1])
	year, err_y := strconv.Atoi(data[2])
	year = year_converter(year)
	if err_d != nil {
		fmt.Println("Error in format of day")
		return "incorect"
	} else if month == 0 {
		fmt.Println("Error in format of month")
		return "incorect"
	} else if err_y != nil {
		fmt.Println("Error in format of year")
		return "incorect"
	}
	if year >= 2000 && year <= 2024 {
		if year%4 == 0 {
			if month == 2 {
				if day >= 1 && day <= 29 {
					return "output"
				} else {
					fmt.Println("Wrong day")
					return "abnormal"
				}
			} else if month%2 == 1 || month == 8 {
				if day >= 1 && day <= 31 {
					return "output"
				} else {
					fmt.Println("Wrong day")
					return "abnormal"
				}
			} else {
				if day >= 1 && day <= 30 {
					return "output"
				} else {
					fmt.Println("Wrong day")
					return "abnormal"
				}
			}
		} else {
			if month == 2 {
				if day >= 1 && day <= 28 {
					return "output"
				} else {
					fmt.Println("Wrong day")
					return "abnormal"
				}
			} else if month%2 == 1 || month == 8 {
				if day >= 1 && day <= 31 {
					return "output"
				} else {
					fmt.Println("Wrong day")
					return "abnormal"
				}
			} else {
				if day >= 1 && day <= 30 {
					return "output"
				} else {
					fmt.Println("Wrong day")
					return "abnormal"
				}
			}
		}
	} else {
		fmt.Println("Wrong year")
		return "abnormal"
	}
}

func check_flight_number(to_analyze string) string {
	if len(to_analyze) < 2 {
		fmt.Println("Not flight number")
		return "incorect"
	}
	_, err := strconv.Atoi(to_analyze[2:])
	if err != nil {
		return "incorect"
	}
	if len(to_analyze) > 6 {
		return "abnormal"
	}
	if 90 >= to_analyze[0] && to_analyze[0] >= 65 &&
		90 >= to_analyze[1] && to_analyze[1] >= 65 && err == nil {
		return "output"
	}
	return "incorect"
}

func check_cost(to_analyze string) string {
	if strings.Contains(to_analyze, "\r") == true {
		to_analyze = to_analyze[:len(to_analyze)-1]
	}
	_, err := strconv.Atoi(to_analyze)
	if err != nil {
		return "incorect"
	}
	return "output"
}

func checker(line string) (string, Record) {
	var to_return Record
	data := check_split(line, ",")
	if len(data) < 7 { //Rework checking for lefted - done
		return "lefted", to_return
	} else if len(data) == 7 {
		data = split(line, ",")
		flight_number_status := check_flight_number(data[0])
		date_status := check_date(data[1])
		start_time_status := check_time(data[2])
		end_time_status := check_time(data[3])
		free_count_status := check_count_places(data[4])
		status_status := check_status(data[5])
		cost_status := check_cost(data[6])
		if flight_number_status == "incorect" ||
			date_status == "incorect" ||
			start_time_status == "incorect" ||
			end_time_status == "incorect" ||
			free_count_status == "incorect" ||
			status_status == "incorect" ||
			cost_status == "incorect" {
			fmt.Println(flight_number_status)
			fmt.Println(date_status)
			fmt.Println(start_time_status)
			fmt.Println(end_time_status)
			fmt.Println(free_count_status)
			fmt.Println(status_status)
			fmt.Println(cost_status)
			return "incorect", to_return
		} else if flight_number_status == "abnormal" ||
			date_status == "abnormal" ||
			start_time_status == "abnormal" ||
			end_time_status == "abnormal" ||
			free_count_status == "abnormal" ||
			status_status == "abnormal" ||
			cost_status == "abnormal" {
			return "abnormal", to_return
		} else if flight_number_status == "output" ||
			date_status == "output" ||
			start_time_status == "output" ||
			end_time_status == "output" ||
			free_count_status == "output" ||
			status_status == "output" ||
			cost_status == "output" {
			var flight_number FlightNumber
			var date Date
			var start_time Time
			var arrival_time Time
			flight_number.id, _ = strconv.Atoi(data[0][2:])
			flight_number.name = data[0][:2]
			to_return.flight_number = flight_number
			date.day, date.month, date.year = eval_date(data[1])
			to_return.date = date
			start_time.minutes, start_time.hours = eval_time(data[2])
			arrival_time.minutes, arrival_time.hours = eval_time(data[3])
			to_return.departure_time = start_time
			to_return.arrival_time = arrival_time
			to_return.free_count, _ = strconv.Atoi(data[4])
			to_return.status = data[5]
			to_return.ticket_cost = eval_cost(data[6])
			return "output", to_return
		}
	}
	return "imposible", to_return
}

func write_to_file_record(record Record, file os.File) {
	file.WriteString(record.String())
}

func write_to_file_line(line string, file os.File) {
	line += "\n"
	file.WriteString(line)
}

func qsort(data []Record, left, right int) []Record { //Need to test more IsBigger (how I see it works correctly => need to check qsort  and test more IsBigger)
	if right-left <= 1 {
		return data
	}
	base_left, base_right := left, right
	new_start_ind := ((right - left) / 2) + left
	start := data[new_start_ind]
	// prev_left, prev_right := left, right
	for left < right {
		for NewIsBigger(start, data[left]) {
			left++
		}
		for NewIsBigger(data[right], start) {
			right--
		}
		data[left], data[right] = data[right], data[left]
	}
	qsort(data, base_left, new_start_ind)
	qsort(data, new_start_ind+1, base_right)
	return data
}

func partition(data []Record, left, right int) int {
	pivot := data[right] // Choosing the last element as pivot
	i := left - 1

	for j := left; j < right; j++ {
		if NewIsBigger(pivot, data[j]) {
			i++
			data[i], data[j] = data[j], data[i] // Swap
		}
	}
	data[i+1], data[right] = data[right], data[i+1] // Place pivot in correct position
	return i + 1
}

func new_qsort(data []Record, left, right int) {
	if left < right {
		pi := partition(data, left, right)
		new_qsort(data, left, pi-1)
		new_qsort(data, pi+1, right)
	}
}

func find_lefted(line string) int {
	data := check_split(line, ",")
	for i := 0; i < len(data); i++ {
		switch i {
		case 0:
			if check_flight_number(data[i]) == "incorect" {
				return 0
			}
		case 1:
			if check_date(data[i]) == "incorect" {
				return 1
			}
		case 2:
			if check_time(data[i]) == "incorect" {
				return 2
			}
		case 3:
			if check_time(data[i]) == "incorect" {
				return 3
			}
		case 4:
			if check_count_places(data[i]) == "incorect" {
				return 4
			}
		case 5:
			if check_status(data[i]) == "incorect" {
				return 5
			}
		case 6:
			if check_cost(data[i]) == "incorect" {
				return 6
			}
		}
	}
	return -1
}

func main() {
	var file_path string
	imposible_file, err_imp := os.Create("imposible.txt")
	output_file, err_out := os.Create("output.txt")
	duplicate_file, err_dup := os.Create("duplicate.txt")
	incorect_file, err_inc := os.Create("incorect.txt")
	lefted_file, err_left := os.Create("lefted.txt")
	abnormal_file, err_abn := os.Create("abnormal.txt")
	manual_file, err_man := os.Create("manual.txt")
	output_data := []Record{}
	if err_imp != nil {
		panic(err_imp)
	}
	if err_out != nil {
		panic(err_out)
	}
	if err_dup != nil {
		panic(err_dup)
	}
	if err_inc != nil {
		panic(err_inc)
	}
	if err_left != nil {
		panic(err_left)
	}
	if err_abn != nil {
		panic(err_abn)
	}
	if err_man != nil {
		panic(err_man)
	}
	// fmt.Println("Enter file name with datatype")
	// fmt.Scan(&file_path)
	file_path = "test.txt"
	file, f_err := os.Open(file_path)
	if f_err != nil {
		panic(f_err)
	}
	b_text, b_err := io.ReadAll(file)
	if b_err != nil {
		panic(b_err)
	}
	text := string(b_text)
	data := split(text, "\n")
	for i := 0; i < len(data); i++ {
		check_status, record := checker(data[i])
		fmt.Println(check_status, data[i])
		if check_status == "lefted" { //rework lefted preprocessing
			index := find_lefted(data[i])
			write_to_file_line(string(index)+" "+data[i], *lefted_file)
		} else if check_status == "incorect" {
			write_to_file_line(data[i], *incorect_file)
		} else if check_status == "abnormal" {
			write_to_file_line(data[i], *abnormal_file)
		} else if check_status == "output" {
			if check_in_output_data(record, output_data) {
				write_to_file_record(record, *duplicate_file)
			} else if check_access(record, output_data) == false {
				write_to_file_record(record, *manual_file)
			} else {
				output_data = append(output_data, record)
			}
		} else {
			write_to_file_line(data[i], *imposible_file)
		}
	}
	fmt.Println(output_data[3])
	fmt.Println(output_data[4])
	fmt.Println(output_data[5])
	new_qsort(output_data, 0, len(output_data)-1)
	write_to_file_line("Flight Number | Flight date | Flight duration | Count of unused places | Day of week | Flight status | Lost profit", *output_file)
	for i := range output_data {
		write_to_file_record(output_data[i], *output_file)
	}
	output_file.Close()
	duplicate_file.Close()
	incorect_file.Close()
	lefted_file.Close()
	abnormal_file.Close()
	manual_file.Close()
	imposible_file.Close()
	fmt.Println("Done")
}
