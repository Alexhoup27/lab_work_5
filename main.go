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
	to_return += record.flight_number.name + "_" + strconv.Itoa(record.flight_number.id)
	to_return += ","
	to_return += strconv.Itoa(record.date.day) + "." + strconv.Itoa(record.date.month) + "." + string(record.date.year)
	to_return += ","
	delta_minutes, delta_hours := eval_time_delta(record.departure_time, record.arrival_time)
	to_return += strconv.Itoa(delta_hours) + ":" + strconv.Itoa(delta_minutes)
	to_return += ","
	to_return += strconv.Itoa(record.free_count)
	to_return += ","
	to_return += eval_day_of_week(record.date)
	to_return += ","
	to_return += record.status
	to_return += ","
	to_return += strconv.Itoa(record.ticket_cost * record.free_count)
	return to_return
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
	if (first_record.free_count * first_record.ticket_cost) >
		(second_record.free_count * second_record.ticket_cost) {
		return true
	} else if (first_record.free_count * first_record.ticket_cost) <
		(second_record.free_count * second_record.ticket_cost) {
		return false
	} else {
		if first_record.flight_number.name > second_record.flight_number.name {
			return true
		} else if first_record.flight_number.name < second_record.flight_number.name &&
			first_record.flight_number.id > second_record.flight_number.id {
			return true
		} else if first_record.flight_number.name < second_record.flight_number.name &&
			first_record.flight_number.id < second_record.flight_number.id {
			return false
		} else {
			if status_converter(first_record.status) > status_converter(second_record.status) {
				return true
			} else if status_converter(first_record.status) < status_converter(second_record.status) {
				return false
			}
		}
	}
	return false
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
	switch month {
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
		if 90 <= data[i] && data[i] <= 65 {
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
	hours, err_h := strconv.Atoi(time[:2])
	minutest, err_m := strconv.Atoi(time[3:5])
	if err_h != nil || err_m != nil || (strings.Contains("p", time) == false &&
		strings.Contains("P", time) == false && strings.Contains("a", time) == false &&
		strings.Contains("A", time) == false) {
		fmt.Println("Wrong type of time ")
		return "incorect"
	}
	if hours <= 0 || hours >= 24 {
		fmt.Println("Wrong hours")
		return "abnormal"
	}
	if minutest <= 0 || minutest >= 60 {
		fmt.Println("Wrong minutes")
		return "abnormal"
	}
	return "output"
}

func check_date(date string) string { // remake with documentation - done
	data := split(date, ".")
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
		if (year/4)%2 == 1 && year%2 == 0 {
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
	_, err := strconv.Atoi(to_analyze[2:])
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
	_, err := strconv.Atoi(to_analyze)
	if err != nil {
		return "incorect"
	}
	return "output"
}

func checker(line string) (string, Record) {
	var to_return Record
	data := split(line, ",")
	if len(data) < 7 {
		return "lefted", to_return
	} else if len(data) == 7 {
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
			flight_number.name = data[0][:1]
			to_return.flight_number = flight_number
			date.day, date.month, date.year = eval_date(data[1])
			to_return.date = date
			start_time.minutes, start_time.hours = eval_time(data[2])
			arrival_time.minutes, arrival_time.hours = eval_time(data[3])
			to_return.departure_time = start_time
			to_return.arrival_time = arrival_time
			to_return.free_count, _ = strconv.Atoi(data[4])
			to_return.status = data[5]
			to_return.ticket_cost, _ = strconv.Atoi(data[6])
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

func qsort(data []Record, left, right int) []Record {
	if len(data) < 2 {
		return data
	}
	start := data[len(data)/2]
	for left <= right {
		for IsBigger(start, data[left]) == true {
			left++
		}
		for IsBigger(data[right], start) == true {
			right--
		}
		data[left], data[right] = data[right], data[left]
	}
	new_start_ind := FindElem(data, start)
	qsort(data, 0, new_start_ind)
	qsort(data, new_start_ind+1, len(data)-1)
	return data
}

func main() {
	var file_path string
	output_file, err_out := os.Create("output.txt")
	duplicate_file, err_dup := os.Create("duplicate.txt")
	incorect_file, err_inc := os.Create("incorect.txt")
	lefted_file, err_left := os.Create("lefted.txt")
	abnormal_file, err_abn := os.Create("abnormal.txt")
	output_data := []Record{}
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
	fmt.Println("Enter file name with datatype")
	fmt.Scan(&file_path)
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
		if check_status == "lefted" {
			write_to_file_line(data[i], *lefted_file)
		} else if check_status == "incorect" {
			write_to_file_line(data[i], *incorect_file)
		} else if check_status == "abnormal" {
			write_to_file_line(data[i], *abnormal_file)
		} else if check_status == "output" {
			if check_in_output_data(record, output_data) {
				write_to_file_record(record, *duplicate_file)
			} else {
				output_data = append(output_data, record)
			}
		}
	}
	output_data = qsort(output_data, 0, len(output_data)-1)
	for i := range output_data {
		write_to_file_line(data[i], *output_file)
	}
	output_file.Close()
	duplicate_file.Close()
	incorect_file.Close()
	lefted_file.Close()
	abnormal_file.Close()
	fmt.Println("Done")
}
