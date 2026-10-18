package view

import (
	"strings"

	"localizacoes-api/internal/model"
)

type Address struct {
	Street       string `json:"street"`
	Number       string `json:"number"`
	Neighborhood string `json:"neighborhood"`
	City         string `json:"city"`
	State        string `json:"state"`
	ZipCode      string `json:"zipCode"`
}

type UPAResponse struct {
	ID           string  `json:"id"`
	Lat          float64 `json:"lat"`
	Long         float64 `json:"long"`
	Title        string  `json:"title"`
	Phone        string  `json:"phone"`
	Address      Address `json:"address"`
	OpeningHours string  `json:"openingHours"`
}

type DaySchedule struct {
	Open  *string `json:"open"`
	Close *string `json:"close"`
}

type UBSResponse struct {
	ID           string                  `json:"id"`
	Lat          float64                 `json:"lat"`
	Long         float64                 `json:"long"`
	Title        string                  `json:"title"`
	Phone        string                  `json:"phone"`
	Address      Address                 `json:"address"`
	OpeningHours map[string]*DaySchedule `json:"openingHours"`
}

func UPAs(models []model.UPA) []UPAResponse {
	var responses []UPAResponse
	for _, upa := range models {
		responses = append(responses, UPAResponse{
			ID:    upa.ID,
			Lat:   upa.Lat,
			Long:  upa.Long,
			Title: upa.Title,
			Phone: upa.Phone,
			Address: Address{
				Street:       upa.Street,
				Number:       upa.Number,
				Neighborhood: upa.Neighborhood,
				City:         upa.City,
				State:        strings.TrimSpace(upa.State),
				ZipCode:      upa.ZipCode,
			},
			OpeningHours: "24/7",
		})
	}
	return responses
}

func UBSs(models []model.UBS) []UBSResponse {
	responses := make([]UBSResponse, 0, len(models))
	for _, ubs := range models {
		openingHours := make(map[string]*DaySchedule, len(ubs.OpeningHours))
		for day, schedule := range ubs.OpeningHours {
			if schedule == nil {
				openingHours[day] = nil
				continue
			}
			openingHours[day] = &DaySchedule{
				Open:  formatTime(schedule.Open),
				Close: formatTime(schedule.Close),
			}
		}

		responses = append(responses, UBSResponse{
			ID:    ubs.ID,
			Lat:   ubs.Lat,
			Long:  ubs.Long,
			Title: ubs.Title,
			Phone: ubs.Phone,
			Address: Address{
				Street:       ubs.Street,
				Number:       ubs.Number,
				Neighborhood: ubs.Neighborhood,
				City:         ubs.City,
				State:        strings.TrimSpace(ubs.State),
				ZipCode:      ubs.ZipCode,
			},
			OpeningHours: openingHours,
		})
	}
	return responses
}

func formatTime(value *string) *string {
	if value == nil || len(*value) < 5 {
		return nil
	}
	formatted := (*value)[:5]
	return &formatted
}
