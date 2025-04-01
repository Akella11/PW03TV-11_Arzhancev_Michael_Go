package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
)

type Input struct {
	Values []float64 `json:"values"`
}

// Функція для обчислення інтегралу Гауса
func integrateGaussian(a, b float64, n int, power, error float64) float64 {
	step := (b - a) / float64(n)
	sum := 0.0
	for i := 0; i < n; i++ {
		x := a + (float64(i)+0.5)*step
		sum += math.Exp(-math.Pow((x-power)/error, 2)/2) / (error * math.Sqrt(2*math.Pi))
	}

	return sum * step
}

// Функція для обчислення фінансових показників до і після вдосконалення
func calculateTask1(power, errorBefore, errorAfter, price float64) string {
	a := power - errorAfter
	b := power + errorAfter
	n := 1000

	shareWithoutImbalancesBefore := integrateGaussian(a, b, n, power, errorBefore)
	profitBefore := power * 24 * shareWithoutImbalancesBefore * price
	fineBefore := power * 24 * (1 - shareWithoutImbalancesBefore) * price

	shareWithoutImbalancesAfter := integrateGaussian(a, b, n, power, errorAfter)
	profitAfter := power * 24 * shareWithoutImbalancesAfter * price
	fineAfter := power * 24 * (1 - shareWithoutImbalancesAfter) * price

	// Форматування результату у вигляді рядка
	output := fmt.Sprintf(
		"Прибуток до вдосконалення: %.2f тис. грн;  \n"+
			"Штраф до вдосконалення: %.2f тис. грн;  \n"+
			"Виручка до вдосконалення: %.2f тис. грн;  \n"+
			"Прибуток після вдосконалення: %.2f тис. грн;  \n"+
			"Штраф після вдосконалення: %.2f тис. грн;  \n"+
			"Виручка після вдосконалення: %.2f тис. грн;  \n",
		profitBefore,
		fineBefore,
		profitBefore-fineBefore,
		profitAfter,
		fineAfter,
		profitAfter-fineAfter,
	)

	return output
}

// Обробник HTTP-запиту для калькулятора
func calculator1Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var input Input
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if len(input.Values) != 4 {
		http.Error(w, "Invalid number of inputs", http.StatusBadRequest)
		return
	}

	result := calculateTask1(input.Values[0], input.Values[1], input.Values[2], input.Values[3])

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

// Основна функція для запуску HTTP-сервера
func main() {
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/calculator1", calculator1Handler)

	fmt.Println("Server running at http://localhost:8083")
	http.ListenAndServe(":8083", nil)
}
