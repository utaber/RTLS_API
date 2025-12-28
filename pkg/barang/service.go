package barang

import (
	"context"
	"fmt"
	"sort"
	"strconv"

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

	var reusableIDs map[string]string
	reusable.Get(s.Ctx, &reusableIDs)

	if len(reusableIDs) > 0 {
		keys := make([]int, 0)
		for k := range reusableIDs {
			i, _ := strconv.Atoi(k)
			keys = append(keys, i)
		}
		sort.Ints(keys)

		id := reusableIDs[strconv.Itoa(keys[0])]
		reusable.Child(strconv.Itoa(keys[0])).Delete(s.Ctx)
		return id
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

/* ===== UPDATE ===== */

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

/* ===== DELETE ===== */

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
