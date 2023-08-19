package events

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_EventDispatcher_NewEventDispatcher(t *testing.T) {
	dispatcher := NewEventDispatcher()
	require.NotNil(t, dispatcher)
}

func Test_EventDispatcher_Register(t *testing.T) {
	// mocks
	var event1Name EventName
	var event2Name EventName

	event1Name = "event1"
	event2Name = "event2"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event1 := NewMockEvent(ctrl)
	event2 := NewMockEvent(ctrl)
	handler1 := NewMockHandler(ctrl)
	handler2 := NewMockHandler(ctrl)

	event1.EXPECT().Name().Return(event1Name).AnyTimes()
	event2.EXPECT().Name().Return(event2Name).AnyTimes()

	// register events and handlers
	dispatcher := NewEventDispatcher()

	err := dispatcher.Register(event1Name, handler1)
	assert.Nil(t, err)
	err = dispatcher.Register(event1Name, handler2)
	assert.Nil(t, err)

	err = dispatcher.Register(event2Name, handler2)
	assert.Nil(t, err)
	err = dispatcher.Register(event2Name, handler1)
	assert.Nil(t, err)

	err = dispatcher.Register(event1Name, handler1)
	assert.ErrorIs(t, err, ErrorHandlerAlreadyRegistered)
	err = dispatcher.Register(event1Name, handler2)
	assert.ErrorIs(t, err, ErrorHandlerAlreadyRegistered)

	err = dispatcher.Register(event2Name, handler2)
	assert.ErrorIs(t, err, ErrorHandlerAlreadyRegistered)
	err = dispatcher.Register(event2Name, handler1)
	assert.ErrorIs(t, err, ErrorHandlerAlreadyRegistered)

	// verify registers
	handlers, ok := dispatcher.handlers[event1Name]
	assert.True(t, ok)
	assert.Equal(t, 2, len(handlers))
	assert.Same(t, handler1, handlers[0])
	assert.Same(t, handler2, handlers[1])

	handlers, ok = dispatcher.handlers[event2Name]
	assert.True(t, ok)
	assert.Equal(t, 2, len(handlers))
	assert.Same(t, handler1, handlers[1])
	assert.Same(t, handler2, handlers[0])
}

func Test_EventDispatcher_Remove(t *testing.T) {
	// mocks
	var event1Name EventName
	var event2Name EventName
	var event3Name EventName

	event1Name = "event1"
	event2Name = "event2"
	event3Name = "event3"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event1 := NewMockEvent(ctrl)
	event2 := NewMockEvent(ctrl)
	event3 := NewMockEvent(ctrl)
	handler1 := NewMockHandler(ctrl)
	handler2 := NewMockHandler(ctrl)

	event1.EXPECT().Name().Return(event1Name).AnyTimes()
	event2.EXPECT().Name().Return(event2Name).AnyTimes()
	event3.EXPECT().Name().Return(event3Name).AnyTimes()

	// register events and handlers
	dispatcher := NewEventDispatcher()

	dispatcher.Register(event1Name, handler1)
	dispatcher.Register(event1Name, handler2)

	dispatcher.Register(event2Name, handler1)
	dispatcher.Register(event2Name, handler2)

	// remove handlers
	err := dispatcher.Remove(event1Name, handler1)
	assert.Nil(t, err)
	err = dispatcher.Remove(event1Name, handler1)
	assert.ErrorIs(t, err, ErrorHandlerNotRegistered)

	err = dispatcher.Remove(event2Name, handler2)
	assert.Nil(t, err)
	err = dispatcher.Remove(event2Name, handler2)
	assert.ErrorIs(t, err, ErrorHandlerNotRegistered)

	err = dispatcher.Remove(event3Name, handler1)
	assert.ErrorIs(t, err, ErrorEventNotRegistered)

	// verify registers
	handlers, ok := dispatcher.handlers[event1Name]
	assert.True(t, ok)
	assert.Equal(t, 1, len(handlers))
	assert.Same(t, handler2, handlers[0])

	handlers, ok = dispatcher.handlers[event2Name]
	assert.True(t, ok)
	assert.Equal(t, 1, len(handlers))
	assert.Same(t, handler1, handlers[0])
}

func Test_EventDispatcher_Has(t *testing.T) {
	// mocks
	var event1Name EventName
	var event2Name EventName

	event1Name = "event1"
	event2Name = "event2"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event1 := NewMockEvent(ctrl)
	event2 := NewMockEvent(ctrl)
	handler1 := NewMockHandler(ctrl)
	handler2 := NewMockHandler(ctrl)

	event1.EXPECT().Name().Return(event1Name).AnyTimes()
	event2.EXPECT().Name().Return(event2Name).AnyTimes()

	// register and verify events and handlers
	dispatcher := NewEventDispatcher()

	assert.False(t, dispatcher.Has(event1Name, handler1))
	assert.False(t, dispatcher.Has(event1Name, handler2))

	assert.False(t, dispatcher.Has(event2Name, handler1))
	assert.False(t, dispatcher.Has(event2Name, handler2))

	dispatcher.Register(event1Name, handler1)
	dispatcher.Register(event2Name, handler2)

	assert.True(t, dispatcher.Has(event1Name, handler1))
	assert.False(t, dispatcher.Has(event1Name, handler2))

	assert.False(t, dispatcher.Has(event2Name, handler1))
	assert.True(t, dispatcher.Has(event2Name, handler2))
}

func Test_EventDispatcher_Dispatch(t *testing.T) {
	// mocks
	var event1Name EventName
	var event2Name EventName

	event1Name = "event1"
	event2Name = "event2"

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event1 := NewMockEvent(ctrl)
	event2 := NewMockEvent(ctrl)
	handler1 := NewMockHandler(ctrl)
	handler2 := NewMockHandler(ctrl)
	handler3 := NewMockHandler(ctrl)
	handler4 := NewMockHandler(ctrl)

	event1.EXPECT().Name().Return(event1Name).AnyTimes()
	event2.EXPECT().Name().Return(event2Name).AnyTimes()

	handler1.EXPECT().Handle(ctx, event1).Return(nil).Times(1)
	handler2.EXPECT().Handle(ctx, event1).Return(nil).Times(1)

	handler3.EXPECT().Handle(ctx, event2).Return(nil).Times(1)
	handler4.EXPECT().Handle(ctx, event2).Return(errors.New("fail to run handler 4")).Times(1)

	// register events and handlers
	dispatcher := NewEventDispatcher()

	dispatcher.Register(event1Name, handler1)
	dispatcher.Register(event1Name, handler2)

	dispatcher.Register(event2Name, handler3)
	dispatcher.Register(event2Name, handler4)

	// verify dispatch
	errs := dispatcher.Dispatch(ctx, event1)
	assert.Nil(t, errs)

	errs = dispatcher.Dispatch(ctx, event2)
	assert.NotNil(t, errs)
	assert.Equal(t, 1, len(errs))
	assert.EqualError(t, errs[0], "fail to run handler 4")
}

func Test_EventDispatcher_Clear(t *testing.T) {
	// mocks
	var event1Name EventName
	var event2Name EventName

	event1Name = "event1"
	event2Name = "event2"

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event1 := NewMockEvent(ctrl)
	event2 := NewMockEvent(ctrl)
	handler1 := NewMockHandler(ctrl)
	handler2 := NewMockHandler(ctrl)

	event1.EXPECT().Name().Return(event1Name).AnyTimes()
	event2.EXPECT().Name().Return(event2Name).AnyTimes()

	// register events and handlers
	dispatcher := NewEventDispatcher()

	dispatcher.Register(event1Name, handler1)
	dispatcher.Register(event2Name, handler2)

	// clean and verify
	assert.Equal(t, 2, len(dispatcher.handlers))

	handlers, ok := dispatcher.handlers[event1Name]
	assert.True(t, ok)
	assert.Equal(t, 1, len(handlers))

	handlers, ok = dispatcher.handlers[event2Name]
	assert.True(t, ok)
	assert.Equal(t, 1, len(handlers))

	dispatcher.Clear()

	assert.Equal(t, 0, len(dispatcher.handlers))

	_, ok = dispatcher.handlers[event1Name]
	assert.False(t, ok)

	_, ok = dispatcher.handlers[event2Name]
	assert.False(t, ok)
}
