package barang

import (
	"context"
	"fmt"

	"RTLS_API/pkg/models"

	"firebase.google.com/go/v4/db"
)

type Service struct {
	DB  *db.Client
	Ctx context.Context
}

func NewService(ctx context.Context, db *db.Client) *Service {
	return &Service{
		Ctx: ctx,
		DB:  db,
	}
}

func (s *Service) GenerateDeviceID() string {
	meta := s.DB.NewRef("meta")
	reusable := meta.Child("reusable_ids")
	counter := meta.Child("device_counter")

	var reused string

	reusable.Transaction(s.Ctx, func(t db.TransactionNode) (interface{}, error) {
		var queue []string
		_ = t.Unmarshal(&queue)

		if len(queue) == 0 {
			return queue, nil
		}

		reused = queue[0]
		return queue[1:], nil
	})

	if reused != "" {
		return reused
	}

	var current int
	counter.Transaction(s.Ctx, func(t db.TransactionNode) (interface{}, error) {
		var old int
		_ = t.Unmarshal(&old)
		current = old + 1
		return current, nil
	})

	return fmt.Sprintf("BOX-%03d", current)
}

func (s *Service) CreateBarang(input models.InputTransaction) (models.OutputTransaction, error) {
	id := s.GenerateDeviceID()

	ref := s.DB.NewRef("Barang").Child(id)
	if err := ref.Set(s.Ctx, input); err != nil {
		return models.OutputTransaction{}, err
	}

	return models.OutputTransaction{
		DeviceID:         id,
		InputTransaction: input,
	}, nil
}

func (s *Service) GetBarang(deviceID string) ([]models.OutputTransaction, error) {
	ref := s.DB.NewRef("Barang")

	var data map[string]models.InputTransaction
	if err := ref.Get(s.Ctx, &data); err != nil {
		return nil, err
	}

	result := []models.OutputTransaction{}
	for k, v := range data {
		if deviceID == "" || k == deviceID {
			result = append(result, models.OutputTransaction{
				DeviceID:         k,
				InputTransaction: v,
			})
		}
	}

	return result, nil
}

func (s *Service) UpdateBarang(id string, payload models.UpdateTransaction) (map[string]interface{}, error) {
	ref := s.DB.NewRef("Barang").Child(id)

	update := map[string]interface{}{}
	if payload.Name != nil {
		update["name"] = *payload.Name
	}

	if len(update) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	if err := ref.Update(s.Ctx, update); err != nil {
		return nil, err
	}

	return update, nil
}

func (s *Service) DeleteBarang(id string) error {
	ref := s.DB.NewRef("Barang").Child(id)

	var check interface{}
	ref.Get(s.Ctx, &check)
	if check == nil {
		return fmt.Errorf("not found")
	}

	ref.Delete(s.Ctx)

	queue := func(t db.TransactionNode) (interface{}, error) {
		var arr []string
		_ = t.Unmarshal(&arr)
		arr = append(arr, id)
		return arr, nil
	}

	s.DB.NewRef("meta/reusable_ids").Transaction(s.Ctx, queue)
	s.DB.NewRef("meta/reusable_ids_esp").Transaction(s.Ctx, queue)

	return nil
}

func (s *Service) ResetSystem() error {
	barangRef := s.DB.NewRef("Barang")
	metaRef := s.DB.NewRef("meta")

	var barangCheck interface{}
	var metaCheck interface{}

	_ = barangRef.Get(s.Ctx, &barangCheck)
	_ = metaRef.Get(s.Ctx, &metaCheck)

	if barangCheck == nil && metaCheck == nil {
		return fmt.Errorf("system already reset")
	}

	if barangCheck != nil {
		if err := barangRef.Delete(s.Ctx); err != nil {
			return err
		}
	}

	if metaCheck != nil {
		if err := metaRef.Delete(s.Ctx); err != nil {
			return err
		}
	}

	return nil
}
