package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("🚀 Starting stress test...")
	
	// Даем сервису время на запуск если нужно
	time.Sleep(2 * time.Second)

	// Запускаем стресс-тест
	cmd := exec.Command("vegeta", "attack", "-duration=10s", "-rate=100", "-targets=stress_test.txt")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal("Error running stress test:", err)
	}

	// Сохраняем результаты
	err = os.WriteFile("stress_results.bin", output, 0644)
	if err != nil {
		log.Fatal("Error saving results:", err)
	}

	// Генерируем отчет
	cmd = exec.Command("vegeta", "report", "stress_results.bin")
	report, err := cmd.Output()
	if err != nil {
		log.Fatal("Error generating report:", err)
	}

	fmt.Println("📊 Stress Test Results:")
	fmt.Println(string(report))
	
	// Генерируем график (опционально)
	cmd = exec.Command("vegeta", "plot", "stress_results.bin")
	plot, err := cmd.Output()
	if err == nil {
		os.WriteFile("stress_plot.html", plot, 0644)
		fmt.Println("📈 Plot saved as stress_plot.html")
	}

	fmt.Println("✅ Stress test completed!")
}