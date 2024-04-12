package events

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type TestEvent struct {
	Name    string
	Payload interface{}
}

func (e *TestEvent) GetName() string {
	return e.Name
}

func (e *TestEvent) GetPayload() interface{} {
	return e.Payload
}

func (e *TestEvent) GetDateTime() time.Time {
	return time.Now()
}

type TestEventHandler struct {
	ID int
}

func (h *TestEventHandler) Handle(event IEvent, wg *sync.WaitGroup) {}

type EventDispatcherTestSuite struct {
	suite.Suite
	event1          TestEvent
	event2          TestEvent
	handler         TestEventHandler
	handler2        TestEventHandler
	handler3        TestEventHandler
	eventDispatcher *EventDispatcher
}

func (suite *EventDispatcherTestSuite) SetupTest() {
	suite.eventDispatcher = NewEventDispatcher()
	suite.event1 = TestEvent{Name: "Event 1", Payload: "Payload 1"}
	suite.event2 = TestEvent{Name: "Event 2", Payload: "Payload 2"}
	suite.handler = TestEventHandler{ID: 1}
	suite.handler2 = TestEventHandler{ID: 2}
	suite.handler3 = TestEventHandler{ID: 3}
}

func (suite *EventDispatcherTestSuite) TearDownTest() {
	suite.eventDispatcher = NewEventDispatcher()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Register() {
	suite.Run("Register handles for an event", func() {
		defer suite.TearDownTest()

		err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler)

		suite.Nil(err)
		suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

		err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)

		suite.Nil(err)
		suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

		event1_handlers := suite.eventDispatcher.handlers[suite.event1.GetName()]
		assert.Equal(suite.T(), &suite.handler, event1_handlers[0])
		assert.Equal(suite.T(), &suite.handler2, event1_handlers[1])
	})

	suite.Run("Register same handler for a same event", func() {
		defer suite.TearDownTest()

		err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler)

		suite.Nil(err)
		suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

		err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler)

		suite.NotNil(err)
		suite.Equal(ErrHandlerAlreadyRegistered, err)
		suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))
	})
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Clear() {
	suite.Run("Clear all handlers", func() {
		defer suite.TearDownTest()

		err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler)
		suite.Nil(err)
		suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

		err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)
		suite.Nil(err)
		suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

		err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler3)
		suite.Nil(err)
		suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

		suite.eventDispatcher.Clear()
		suite.Empty(suite.eventDispatcher.handlers)
	})
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Has() {
	suite.Run("When handlers are included", func() {
		defer suite.TearDownTest()

		err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler)
		suite.Nil(err)
		suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

		err = suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler2)
		suite.Nil(err)
		suite.Equal(2, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

		assert.True(suite.T(), suite.eventDispatcher.Has(suite.event1.GetName(), &suite.handler))
		assert.True(suite.T(), suite.eventDispatcher.Has(suite.event1.GetName(), &suite.handler2))
	})

	suite.Run("When handlers are not included", func() {
		defer suite.TearDownTest()

		err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler)
		suite.Nil(err)
		suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

		assert.False(suite.T(), suite.eventDispatcher.Has(suite.event1.GetName(), &suite.handler2))
		assert.False(suite.T(), suite.eventDispatcher.Has(suite.event2.GetName(), &suite.handler2))
	})
}

type MockEventHandler struct {
	mock.Mock
}

func (m *MockEventHandler) Handle(event IEvent, wg *sync.WaitGroup) {
	m.Called(event)
	wg.Done()
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Dispatch() {
	suite.Run("Should execute handlers of the event using multithreads", func() {
		defer suite.TearDownTest()

		handler := &MockEventHandler{}
		handler.On("Handle", &suite.event1)

		handler2 := &MockEventHandler{}
		handler2.On("Handle", &suite.event1)

		handler3 := &MockEventHandler{}
		handler3.On("Handle", &suite.event1)

		suite.eventDispatcher.Register(suite.event1.GetName(), handler)
		suite.eventDispatcher.Register(suite.event1.GetName(), handler2)
		suite.eventDispatcher.Register(suite.event1.GetName(), handler3)

		suite.eventDispatcher.Dispatch(&suite.event1)

		handler.AssertExpectations(suite.T())
		handler.AssertNumberOfCalls(suite.T(), "Handle", 1)

		handler2.AssertExpectations(suite.T())
		handler2.AssertNumberOfCalls(suite.T(), "Handle", 1)

		handler3.AssertExpectations(suite.T())
		handler3.AssertNumberOfCalls(suite.T(), "Handle", 1)
	})
}

func (suite *EventDispatcherTestSuite) TestEventDispatcher_Remove() {
	suite.Run("Remove handler from event dispatcher", func() {
		defer suite.TearDownTest()

		err := suite.eventDispatcher.Register(suite.event1.GetName(), &suite.handler)
		suite.Nil(err)
		suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event1.GetName()]))

		err = suite.eventDispatcher.Register(suite.event2.GetName(), &suite.handler)
		suite.Nil(err)
		suite.Equal(1, len(suite.eventDispatcher.handlers[suite.event2.GetName()]))

		suite.eventDispatcher.Remove(suite.event1.GetName(), &suite.handler)
		suite.Empty(suite.eventDispatcher.handlers[suite.event1.GetName()])

		suite.eventDispatcher.Remove(suite.event2.GetName(), &suite.handler)
		suite.Empty(suite.eventDispatcher.handlers[suite.event2.GetName()])
	})
}

func TestSuite(t *testing.T) {
	testSuite := new(EventDispatcherTestSuite)
	testSuite.SetupTest()

	suite.Run(t, testSuite)
}
