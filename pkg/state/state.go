package state

import "sync"

// Listener is a function that will be invoked whenever the reactive
// value changes. It receives the new value as its argument.
// Generics are used so the listener can work with any underlying
// type. Listeners should not block for long periods of time.
type Listener[T any] func(value T)

// Reactive manages a single piece of state and notifies all registered
// listeners when the state changes. It is safe for concurrent use.
// A zero value Reactive is not safe for use; always construct via
// NewReactive.
type Reactive[T any] struct {
    mu        sync.Mutex
    value     T
    listeners []Listener[T]
}

// NewReactive returns a new Reactive wrapping the provided initial
// value. The returned Reactive can be used to store and update state.
func NewReactive[T any](initial T) *Reactive[T] {
    return &Reactive[T]{value: initial}
}

// Get returns the current value of r. It acquires a lock to ensure
// consistent reads across concurrent writers.
func (r *Reactive[T]) Get() T {
    r.mu.Lock()
    defer r.mu.Unlock()
    return r.value
}

// Set updates the reactive value to val and notifies all registered
// listeners. Listeners will be invoked in the order they were
// registered. Set copies the current listeners slice under the
// protection of the mutex before unlocking, so listeners may call
// Subscribe safely from within callbacks.
func (r *Reactive[T]) Set(val T) {
    r.mu.Lock()
    r.value = val
    listenersCopy := append([]Listener[T](nil), r.listeners...)
    r.mu.Unlock()
    for _, l := range listenersCopy {
        l(val)
    }
}

// Subscribe registers a listener that will be called whenever the
// reactive value changes. Multiple calls to Subscribe will append
// additional listeners, which will be invoked on every call to Set.
// It is safe to register the same listener multiple times, though it
// will then be called multiple times per Set. Listeners should not
// call Set on the same Reactive, otherwise a deadlock may occur. To
// update state from within a listener use a separate Reactive.
func (r *Reactive[T]) Subscribe(listener Listener[T]) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.listeners = append(r.listeners, listener)
}
