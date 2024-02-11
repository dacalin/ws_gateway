package _connection_id

type ConnectionId struct {
	id string
}

func New(id string) ConnectionId {
	return ConnectionId{id: id}
}

func (obj ConnectionId) Value() string {
	return obj.id
}
