package ws

import "github.com/Slimo300/MicroservicesChatApp/backend/group-service/communication"

type HubInterface interface {
	Run()
	Join(*client)
	Leave(*client)
	Forward(*communication.Message)
}

type MockHub struct {
	actionChan <-chan *communication.Action
}

func (m *MockHub) Run() {
	for {
		<-m.actionChan
	}
}
func (m *MockHub) Join(c *client)                     {}
func (m *MockHub) Leave(c *client)                    {}
func (m *MockHub) Forward(msg *communication.Message) {}

func NewMockHub(actionChan <-chan *communication.Action) *MockHub {
	return &MockHub{
		actionChan: actionChan,
	}
}
