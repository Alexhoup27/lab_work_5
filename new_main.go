package main

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
)

var file_path string
var flight_number_status, date_status, start_time_status, end_time_status, places_status, status_status string

// var flight_number_flag, date_flag, start_time_flag, end_time_flag, places_flag, status_flag bool

// переделать evalы для нужной мне сортировки и переделать саму сортировку

func eval_time_delta(start_time int, end_time int) int {
	if end_time < start_time {
		end_time += 24 * 60
	}
	return end_time - start_time
}

func eval_time(data string) int {
	hours, _ := strconv.Atoi(data[:2])
	minutes, _ := strconv.Atoi(data[3:5])
	if strings.Contains("P", data) || strings.Contains("p", data) {
		hours += 12
	}
	return hours*24 + minutes
}

func eval_number(data []string) int {
	_sum := int32(0)
	for i := 0; i < len(data[0]); i++ {
		_sum += rune(data[0][i])
	}
	return int(_sum)
}

func eval_day_of_week(data string) string {
	new_data := split(data, ".")
	year, _ := strconv.Atoi(new_data[0])
	year = year_converter(year)
	month := month_converter(new_data[1])
	day, _ := strconv.Atoi(new_data[2])
	_sum := int64(0)
	year -= 2000
	_sum += int64(year/4) * int64(366)
	_sum += int64(year-year/4) * int64(365)
	for i := 1; i < month; i++ {
		if i == 2 && year%4 == 0 {
			_sum += 29
		} else if i == 2 {
			_sum += 28
		} else if i == 8 || i%2 == 1 {
			_sum += 31
		} else {
			_sum += 30
		}
	}
	_sum += int64(day)
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

func eval_status(data string) int {
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

func eval_date(data []string) int {
	date := split(data[1], "/")
	hours, _ := strconv.Atoi(date[0])
	minutes, _ := strconv.Atoi(date[1])
	return hours*60 + minutes
}

func eval_lost_income(free_seats_count string, cost string) int {
	count_int, _ := strconv.Atoi(free_seats_count)
	cost_int, _ := strconv.Atoi(cost)
	return cost_int * count_int
}

//create new qsort (done) and check 17 line

func qsort(data [][]string) [][]string {
	if len(data) < 2 {
		return data
	}
	start := len(data) / 2
	for i := range data {
		now_lost, _ := strconv.Atoi(data[i][6])
		start_lost, _ := strconv.Atoi(data[start][6])
		if now_lost > start_lost { // recreate all convditions with documentation - done
			elem := data[i]
			data = slices.Concat(data[:i], data[i+1:start], data[start:])
			data = append(data, elem)
		} else if now_lost < start_lost {
			new_data := [][]string{}
			new_data[0] = data[i]
			for g := 0; g < len(data); g++ {
				if g != i {
					new_data = append(new_data, data[g])
				}
			}
			data = new_data
		} else {
			if eval_number(data[i][0]) > eval_number(data[start][0]) {
				elem := data[i]
				data = slices.Concat(data[:i], data[i+1:start], data[start:])
				data = append(data, elem)
			} else if eval_number(data[i][0]) < eval_number(data[start][0]) {
				new_data := [][]string{}
				new_data[0] = data[i]
				for g := 0; g < len(data); g++ {
					if g != i {
						new_data = append(new_data, data[g])
					}
				}
				data = new_data
			} else {
				if eval_status(data[i][5]) > eval_status(data[start][5]) {
					elem := data[i]
					data = slices.Concat(data[:i], data[i+1:start], data[start:])
					data = append(data, elem)
				} else if eval_status(data[i][5]) < eval_status(data[start][5]) {
					new_data := [][]string{}
					new_data[0] = data[i]
					for g := 0; g < len(data); g++ {
						if g != i {
							new_data = append(new_data, data[g])
						}
					}
					data = new_data
				}
			}
		}
	}
	return data
}

func read_line(data string) []string {
	to_return := []string{}
	to_add := ""
	comma_flag := false
	for i := 0; i < len(data); i++ {
		if data[i] != ',' && comma_flag == false {
			to_add += string(data[i])
		} else if data[i] != ',' && comma_flag == true {
			to_return = append(to_return, to_add)
			to_add = string(data[i])
		} else {
			comma_flag = true
		}
	}
	if comma_flag {
		to_return = append(to_return, "")
	} else {
		to_return = append(to_return, to_add)
	}
	return to_return
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

func check_flight_number(to_analyze string) string {
	_, err := strconv.Atoi(to_analyze[2:])
	if len(to_analyze) > 6 {
		return "abnormal"
	}
	if 90 >= to_analyze[0] && to_analyze[0] >= 65 &&
		90 >= to_analyze[1] && to_analyze[1] >= 65 && err != nil {
		return "output"
	}
	return "incorect"
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

func check_time(time string) string { // remake with documentation - done
	// data := split(time, "/")
	hours, err_h := strconv.Atoi(time[:2])
	minutest, err_m := strconv.Atoi(time[3:5])
	if err_h != nil || err_m != nil || (strings.Contains("p", time) == false &&
		strings.Contains("P", time) == false && strings.Contains("a", time) == false &&
		strings.Contains("A", time) == false) {
		fmt.Println("Wron type of time ")
		return "incorect"
	}
	if hours <= 0 || hours >= 24 {
		fmt.Println("Wrong hours")
		return "abnormal"
	}
	if minutest <= 0 && minutest >= 60 {
		fmt.Println("Wrong minutes")
		return "abnormal"
	}
	return "output"
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

func check_in_writer(writer [][]string, line []string) bool {
	for i := 0; i < len(writer); i++ {
		flag := true
		for g := 0; g < len(writer[i]); g++ {
			if writer[i][g] != line[g] {
				flag = false
			}
		}
		if flag {
			return true
		}
	}
	return false
}

func main() {
	fmt.Println("Enter path to file")
	fmt.Scan(&file_path)
	file, f_err := os.Open(file_path)
	output_writer := [][]string{}
	duplicate_writer := [][]string{}
	abnormal_writer := [][]string{}
	lefted_writer := [][]string{}
	incorect_writer := [][]string{}
	// file, f_err := os.Open("test.txt")
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
		line := read_line(data[i])
		if len(line) < 6 {
			lefted_writer = append(lefted_writer, line)
		} else if len(line) == 6 {
			flight_number_status = check_flight_number(line[0])
			date_status = check_date(line[1])
			start_time_status = check_time(line[2])
			end_time_status = check_time(line[3])
			places_status = check_count_places(line[4])
			status_status = check_status(line[5])
			if flight_number_status == "incorect" ||
				date_status == "incorect" ||
				start_time_status == "incorect" ||
				end_time_status == "incorect" ||
				places_status == "incorect" ||
				status_status == "incorect" {
				incorect_writer = append(incorect_writer, line)
			} else if flight_number_status == "abnormal" ||
				date_status == "abnormal" ||
				start_time_status == "abnormal" ||
				end_time_status == "abnormal" ||
				places_status == "abnormal" ||
				status_status == "abnormal" {
				abnormal_writer = append(abnormal_writer, line)
			} else if flight_number_status == "ouput" &&
				date_status == "ouput" &&
				start_time_status == "ouput" &&
				end_time_status == "ouput" &&
				places_status == "ouput" &&
				status_status == "ouput" {
				if check_in_writer(output_writer, line) {
					duplicate_writer = append(duplicate_writer, line)
				} else {
					output_writer = append(output_writer, line)
				}
			} else {
				incorect_writer = append(incorect_writer, line)
			}
		}
	}
	for i := 0; i < len(output_writer); i++ {
		now_line := []string{}
		now_line[0] = output_writer[i][0]
		now_line[1] = output_writer[i][1]
		// ask about format of output time delta
		now_line[2] = string(eval_time_delta(eval_time(output_writer[i][2]), eval_time(output_writer[i][3])))
		now_line[3] = output_writer[i][4]
		now_line[4] = eval_day_of_week(output_writer[i][1])
		now_line[5] = output_writer[i][5]
		now_line[6] = string(eval_lost_income(output_writer[i][4], output_writer[i][6]))
	}
	// go to sort and after it writing
	output_writer = qsort(output_writer)
	output_file, err_out := os.Create("output.txt")
	duplicate_file, err_dup := os.Create("duplicate.txt")
	incorect_file, err_inc := os.Create("incorect.txt")
	lefted_file, err_left := os.Create("lefted.txt")
	abnormal_file, err_abn := os.Create("abnormal.txt")
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
	// go to file writing
	for i := 0; i < len(output_writer); i++ {
		line := ""
		for g := 0; i < len(output_writer[i])-1; g++ {
			line += output_writer[i][g] + ","
		}
		line += output_writer[i][6] + "\n"
		output_file.WriteString(line)
	}
	for i := 0; i < len(duplicate_writer); i++ {
		line := ""
		for g := 0; i < len(duplicate_writer[i])-1; g++ {
			line += duplicate_writer[i][g] + ","
		}
		line += duplicate_writer[i][6] + "\n"
		duplicate_file.WriteString(line)
	}
	for i := 0; i < len(incorect_writer); i++ {
		line := ""
		for g := 0; i < len(incorect_writer[i])-1; g++ {
			line += incorect_writer[i][g] + ","
		}
		line += incorect_writer[i][6] + "\n"
		incorect_file.WriteString(line)
	}
	for i := 0; i < len(lefted_writer); i++ {
		line := ""
		for g := 0; i < len(lefted_writer[i])-1; g++ {
			line += lefted_writer[i][g] + ","
		}
		line += lefted_writer[i][6] + "\n"
		lefted_file.WriteString(line)
	}
	for i := 0; i < len(abnormal_writer); i++ {
		line := ""
		for g := 0; i < len(abnormal_writer[i])-1; g++ {
			line += abnormal_writer[i][g] + ","
		}
		line += abnormal_writer[i][6] + "\n"
		abnormal_file.WriteString(line)
	}
	fmt.Println("WW")
}
