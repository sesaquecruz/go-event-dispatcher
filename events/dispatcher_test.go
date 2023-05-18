package events

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EventDispatcher_NewEventDispatcher(t *testing.T) {
	dispatcher := NewEventDispatcher()
	require.NotNil(t, dispatcher)
}

func Test_EventDispatcher_Register(t *testing.T) {
	// mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event1 := NewMockEvent(ctrl)
	event2 := NewMockEvent(ctrl)
	handler1 := NewMockHandler(ctrl)
	handler2 := NewMockHandler(ctrl)

	event1.EXPECT().Name().Return("event1").AnyTimes()
	event2.EXPECT().Name().Return("event2").AnyTimes()

	// register events and handlers
	dispatcher := NewEventDispatcher()

	err := dispatcher.Register(event1, handler1)
	assert.Nil(t, err)
	err = dispatcher.Register(event1, handler2)
	assert.Nil(t, err)

	err = dispatcher.Register(event2, handler2)
	assert.Nil(t, err)
	err = dispatcher.Register(event2, handler1)
	assert.Nil(t, err)

	err = dispatcher.Register(event1, handler1)
	assert.ErrorIs(t, err, ErrorHandlerAlreadyRegistered)
	err = dispatcher.Register(event1, handler2)
	assert.ErrorIs(t, err, ErrorHandlerAlreadyRegistered)

	err = dispatcher.Register(event2, handler2)
	assert.ErrorIs(t, err, ErrorHandlerAlreadyRegistered)
	err = dispatcher.Register(event2, handler1)
	assert.ErrorIs(t, err, ErrorHandlerAlreadyRegistered)

	// verify registers
	handlers, ok := dispatcher.handlers["event1"]
	assert.True(t, ok)
	assert.Equal(t, 2, len(handlers))
	assert.Same(t, handler1, handlers[0])
	assert.Same(t, handler2, handlers[1])

	handlers, ok = dispatcher.handlers["event2"]
	assert.True(t, ok)
	assert.Equal(t, 2, len(handlers))
	assert.Same(t, handler1, handlers[1])
	assert.Same(t, handler2, handlers[0])
}

func Test_EventDispatcher_Remove(t *testing.T) {
	// mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event1 := NewMockEvent(ctrl)
	event2 := NewMockEvent(ctrl)
	event3 := NewMockEvent(ctrl)
	handler1 := NewMockHandler(ctrl)
	handler2 := NewMockHandler(ctrl)

	event1.EXPECT().Name().Return("event1").AnyTimes()
	event2.EXPECT().Name().Return("event2").AnyTimes()
	event3.EXPECT().Name().Return("event3").AnyTimes()

	// register events and handlers
	dispatcher := NewEventDispatcher()

	dispatcher.Register(event1, handler1)
	dispatcher.Register(event1, handler2)

	dispatcher.Register(event2, handler1)
	dispatcher.Register(event2, handler2)

	// remove handlers
	err := dispatcher.Remove(event1, handler1)
	assert.Nil(t, err)
	err = dispatcher.Remove(event1, handler1)
	assert.ErrorIs(t, err, ErrorHandlerNotRegistered)

	err = dispatcher.Remove(event2, handler2)
	assert.Nil(t, err)
	err = dispatcher.Remove(event2, handler2)
	assert.ErrorIs(t, err, ErrorHandlerNotRegistered)

	err = dispatcher.Remove(event3, handler1)
	assert.ErrorIs(t, err, ErrorEventNotRegistered)

	// verify registers
	handlers, ok := dispatcher.handlers["event1"]
	assert.True(t, ok)
	assert.Equal(t, 1, len(handlers))
	assert.Same(t, handler2, handlers[0])

	handlers, ok = dispatcher.handlers["event2"]
	assert.True(t, ok)
	assert.Equal(t, 1, len(handlers))
	assert.Same(t, handler1, handlers[0])
}

func Test_EventDispatcher_Has(t *testing.T) {
	// mocks
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event1 := NewMockEvent(ctrl)
	event2 := NewMockEvent(ctrl)
	handler1 := NewMockHandler(ctrl)
	handler2 := NewMockHandler(ctrl)

	event1.EXPECT().Name().Return("event1").AnyTimes()
	event2.EXPECT().Name().Return("event2").AnyTimes()

	// register and verify events and handlers
	dispatcher := NewEventDispatcher()

	assert.False(t, dispatcher.Has(event1, handler1))
	assert.False(t, dispatcher.Has(event1, handler2))

	assert.False(t, dispatcher.Has(event2, handler1))
	assert.False(t, dispatcher.Has(event2, handler2))

	dispatcher.Register(event1, handler1)
	dispatcher.Register(event2, handler2)

	assert.True(t, dispatcher.Has(event1, handler1))
	assert.False(t, dispatcher.Has(event1, handler2))

	assert.False(t, dispatcher.Has(event2, handler1))
	assert.True(t, dispatcher.Has(event2, handler2))
}

func Test_EventDispatcher_Dispatch(t *testing.T) {
	// mocks
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event1 := NewMockEvent(ctrl)
	event2 := NewMockEvent(ctrl)
	handler1 := NewMockHandler(ctrl)
	handler2 := NewMockHandler(ctrl)
	handler3 := NewMockHandler(ctrl)
	handler4 := NewMockHandler(ctrl)

	event1.EXPECT().Name().Return("event1").AnyTimes()
	event2.EXPECT().Name().Return("event2").AnyTimes()

	handler1.EXPECT().Handle(ctx, event1).Return(nil).Times(1)
	handler2.EXPECT().Handle(ctx, event1).Return(nil).Times(1)

	handler3.EXPECT().Handle(ctx, event2).Return(nil).Times(1)
	handler4.EXPECT().Handle(ctx, event2).Return(errors.New("fail to run handler 4")).Times(1)

	// register events and handlers
	dispatcher := NewEventDispatcher()

	dispatcher.Register(event1, handler1)
	dispatcher.Register(event1, handler2)

	dispatcher.Register(event2, handler3)
	dispatcher.Register(event2, handler4)

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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	event1 := NewMockEvent(ctrl)
	event2 := NewMockEvent(ctrl)
	handler1 := NewMockHandler(ctrl)
	handler2 := NewMockHandler(ctrl)

	event1.EXPECT().Name().Return("event1").AnyTimes()
	event2.EXPECT().Name().Return("event2").AnyTimes()

	// register events and handlers
	dispatcher := NewEventDispatcher()

	dispatcher.Register(event1, handler1)
	dispatcher.Register(event2, handler2)

	// clean and verify
	assert.Equal(t, 2, len(dispatcher.handlers))

	handlers, ok := dispatcher.handlers["event1"]
	assert.True(t, ok)
	assert.Equal(t, 1, len(handlers))

	handlers, ok = dispatcher.handlers["event2"]
	assert.True(t, ok)
	assert.Equal(t, 1, len(handlers))

	dispatcher.Clear()

	assert.Equal(t, 0, len(dispatcher.handlers))

	_, ok = dispatcher.handlers["event1"]
	assert.False(t, ok)

	_, ok = dispatcher.handlers["event2"]
	assert.False(t, ok)
}
