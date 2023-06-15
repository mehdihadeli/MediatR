package contracts

type INotification interface {
	IsNotification() bool
}

type Notification struct {
}

func NewNotification() INotification {
	return &Notification{}
}

func (r *Notification) IsNotification() bool {
	return true
}
