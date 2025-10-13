package delete

// IUsecase - ...
type IUsecase interface {
	Execute(id string) error
}
