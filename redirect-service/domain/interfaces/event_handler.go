package interfaces

type EventHandler interface {
	HandleEvent(msg []byte)
}
