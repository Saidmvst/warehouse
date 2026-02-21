package main

import (
	"encoding/json"
	"fmt"
)

type Item struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Warehouse struct {
	Items map[int]Item
	L     Logger
}

type Logger interface {
	Log(message string)
}

type ConsoleLogger struct{}
type SimpleLogger struct{}

func (c ConsoleLogger) Log(message string) {
	fmt.Println(message)
}

func (s *SimpleLogger) Log(message string) {
	pref := "[LOG]: "
	fmt.Println(pref + message)
}

func (w *Warehouse) AddItem(item Item) {
	w.Items[item.ID] = item
	w.L.Log("Добавлен товар: " + item.Name)
}

func (w *Warehouse) GetTotalValue() float64 {
	var total float64
	for _, item := range w.Items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}

func (w *Warehouse) FindLowStock(limit int) []Item {
	var lowStockItems []Item
	for _, item := range w.Items {
		if item.Quantity < limit {
			lowStockItems = append(lowStockItems, item)
		}
	}
	return lowStockItems
}

func (w *Warehouse) UpdateQuantity(id int, newQty int) error {
	val, ok := w.Items[id]
	if ok {
		val.Quantity = newQty
		w.Items[id] = val
		fmt.Printf("Новое количество ID: %d %d\n", id, val.Quantity)
		return nil
	}
	return fmt.Errorf("Ошибка: товар с ID: %d не найден", id)
}

func (w *Warehouse) RemoveItem(id int) error {
	delete(w.Items, id)
	return nil
}

func main() {
	myWarehouse := Warehouse{
		Items: make(map[int]Item),
		L:     ConsoleLogger{},
	}

	savedMarsh := `{"1": {"id": 1, "name": "Twix", "price": 40, "quantity": 50}}`
	ware := make(map[int]Item)
	if err := json.Unmarshal([]byte(savedMarsh), &ware); err != nil {
		fmt.Println(err)
	}

	lays := Item{ID: 1, Name: "Lays", Price: 50.0, Quantity: 10}
	milka := Item{ID: 2, Name: "Milka", Price: 120.0, Quantity: 24}
	stom := Item{ID: 3, Name: "Стом", Price: 60.0, Quantity: 12}

	myWarehouse.AddItem(lays)
	myWarehouse.AddItem(milka)
	myWarehouse.AddItem(stom)

	fmt.Printf("Общая стоимость склада: %.2f\n", myWarehouse.GetTotalValue())

	lower := myWarehouse.FindLowStock(5)
	for _, item := range lower {
		fmt.Println(item.Name)
	}

	if err := myWarehouse.UpdateQuantity(3, 24); err != nil {
		fmt.Println("Ошибка!", err)
		return
	}

	if err := myWarehouse.RemoveItem(2); err != nil {
		fmt.Println("Ошибка!", err)
		return
	}

	fmt.Println(myWarehouse.Items)

	jsonData, err := json.MarshalIndent(myWarehouse, "", "    ")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonData))

	fmt.Println("Загружено из JSON:", ware)
	fmt.Println("Программа складского учета завершена")
}
