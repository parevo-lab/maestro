package main

/**

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/parevo-lab/maestro"
)

func main() {
	// Engine oluştur
	wfEngine := maestro.NewEngine()

	// Observer ekle
	wfEngine.AddObserver(func(event maestro.Event) {
		log.Printf("Event: %s, Step: %s, Data: %v\n", event.Type, event.StepID, event.Data)
	})

	// Adımları kaydet
	wfEngine.RegisterStep("validate", func(ctx context.Context, data interface{}) (interface{}, error) {
		log.Println("Validating data...")
		time.Sleep(1 * time.Second)
		return "validated", nil
	})

	wfEngine.RegisterStep("process", func(ctx context.Context, data interface{}) (interface{}, error) {
		log.Println("Processing data...")
		time.Sleep(2 * time.Second)
		return "processed", nil
	})

	wfEngine.RegisterStep("notify", func(ctx context.Context, data interface{}) (interface{}, error) {
		log.Println("Sending notification...")
		time.Sleep(1 * time.Second)
		return "notified", nil
	})

	// İş akışı tanımı oluştur
	workflow := maestro.NewWorkflowDefinition(
		"order-flow",
		"Sipariş İşleme",
		"Yeni siparişleri işleme iş akışı",
	)

	// Adımları tanımla
	validateStep := maestro.NewStepDefinition("validate", "Validasyon", maestro.StepTypeTask).
		WithTimeout(5 * time.Second).
		WithNextSteps("process")

	processStep := maestro.NewStepDefinition("process", "İşleme", maestro.StepTypeTask).
		WithTimeout(10 * time.Second).
		WithNextSteps("notify")

	notifyStep := maestro.NewStepDefinition("notify", "Bildirim", maestro.StepTypeTask).
		WithTimeout(5 * time.Second)

	workflow.AddStep(validateStep)
	workflow.AddStep(processStep)
	workflow.AddStep(notifyStep)

	// Çalışma zamanı oluştur
	runtime := maestro.NewWorkflowRuntime(wfEngine, workflow)

	// İş akışını başlat
	ctx := context.Background()
	if err := runtime.Start(ctx); err != nil {
		log.Fatalf("İş akışı başlatılamadı: %v", err)
	}

	// İş akışı durumunu kontrol et
	state := runtime.GetState()
	fmt.Printf("İş akışı tamamlandı. Durum: %s\n", state.Status)

	for stepID, result := range state.StepResults {
		fmt.Printf("Adım %s sonucu: %v\n", stepID, result)
	}
}


**/
