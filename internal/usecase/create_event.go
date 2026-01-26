package usecase

type CreateEventInputDTO struct {
	Name         string  `json:"name"`
	Location     string  `json:"location"`
	Organization string  `json:"organization"`
	Rating       string  `json:"rating"`
	Date         string  `json:"date"` // Recebemos string (ISO8601) e convertemos
	Capacity     int     `json:"capacity"`
	Price        float64 `json:"price"`
	ImageURL     string  `json:"image_url"`
	PartnerID    int     `json:"partner_id"`
}

type CreateEventOutputDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateEventUseCase struct {
	eventRepo EventRepository
}

func NewCreateEventUseCase(eventRepo EventRepository, txManager TransactionManager) *CreateEventUseCase {
	return &CreateEventUseCase{
		eventRepo: eventRepo,
	}
}
