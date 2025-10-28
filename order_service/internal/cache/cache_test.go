package cache

import (
	"order_service/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache_SetAndGet(t *testing.T) {
	cache := New()
	
	order := &models.Order{
		OrderUID:    "test123",
		TrackNumber: "TRACK123",
		Entry:       "WBIL",
		DateCreated: time.Now(),
	}

	// Тестируем Set
	cache.Set(order)

	// Тестируем Get - должен найти
	result, exists := cache.Get("test123")
	assert.True(t, exists, "Заказ должен существовать в кэше")
	assert.Equal(t, order, result, "Заказы должны совпадать")

	// Тестируем Get - не должен найти
	_, exists = cache.Get("nonexistent")
	assert.False(t, exists, "Несуществующий заказ не должен быть найден")
}

func TestCache_GetAll(t *testing.T) {
	cache := New()

	order1 := &models.Order{OrderUID: "test1", TrackNumber: "TRACK1"}
	order2 := &models.Order{OrderUID: "test2", TrackNumber: "TRACK2"}

	cache.Set(order1)
	cache.Set(order2)

	// Тестируем GetAll
	allOrders := cache.GetAll()
	
	assert.Equal(t, 2, len(allOrders), "Должно быть 2 заказа")
	assert.Equal(t, order1, allOrders["test1"])
	assert.Equal(t, order2, allOrders["test2"])
}

func TestCache_Restore(t *testing.T) {
	cache := New()

	orders := map[string]*models.Order{
		"order1": {OrderUID: "order1", TrackNumber: "TRACK1"},
		"order2": {OrderUID: "order2", TrackNumber: "TRACK2"},
	}

	// Тестируем Restore
	cache.Restore(orders)

	// Проверяем что данные восстановились
	result, exists := cache.Get("order1")
	assert.True(t, exists)
	assert.Equal(t, "TRACK1", result.TrackNumber)

	result, exists = cache.Get("order2")
	assert.True(t, exists)
	assert.Equal(t, "TRACK2", result.TrackNumber)
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := New()
	order := &models.Order{OrderUID: "concurrent", TrackNumber: "CONCURRENT"}

	// Запускаем несколько горутин для тестирования конкурентного доступа
	done := make(chan bool)

	// Горутина для записи
	go func() {
		for i := 0; i < 100; i++ {
			cache.Set(order)
		}
		done <- true
	}()

	// Горутина для чтения
	go func() {
		for i := 0; i < 100; i++ {
			cache.Get("concurrent")
		}
		done <- true
	}()

	// Ожидаем завершения
	<-done
	<-done

	// Проверяем что данные не повредились
	result, exists := cache.Get("concurrent")
	assert.True(t, exists)
	assert.Equal(t, order, result)
}