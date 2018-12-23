package base

import (
	"context"
	"errors"
	"sync"

	"github.com/go-kit/kit/endpoint"
)

// Service is an interface for the profile functions
type Service interface {
	PostProfile(ctx context.Context, p Profile) error
	GetProfile(ctx context.Context, id string) (Profile, error)
	PutProfile(ctx context.Context, id string, p Profile) error
	DeleteProfile(ctx context.Context, id string) error
	FetchData() error
}

// Profile represents a single user profile.
type Profile struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Addresses []Address `json:"addresses,omitempty"`
}

// Address is a field of a user profile.
type Address struct {
	ID       string `json:"id"`
	Location string `json:"location,omitempty"`
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type inmemService struct {
	mtx             sync.RWMutex
	m               map[string]Profile
	pokemonEndpoint endpoint.Endpoint
}

func NewInmemService(pokemonEndpoint endpoint.Endpoint) Service {
	return &inmemService{
		m:               map[string]Profile{},
		pokemonEndpoint: pokemonEndpoint,
	}
}

func (s *inmemService) PostProfile(ctx context.Context, p Profile) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[p.ID]; ok {
		return ErrAlreadyExists
	}
	s.m[p.ID] = p
	return nil
}

func (s *inmemService) GetProfile(ctx context.Context, id string) (Profile, error) {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	p, ok := s.m[id]
	if !ok {
		return Profile{}, ErrNotFound
	}
	return p, nil
}

func (s *inmemService) PutProfile(ctx context.Context, id string, p Profile) error {
	if id != p.ID {
		return ErrInconsistentIDs
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[id] = p
	return nil
}

func (s *inmemService) DeleteProfile(ctx context.Context, id string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}
	delete(s.m, id)
	return nil
}

func (s *inmemService) FetchData() error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	//make outbound request
	res, err := s.pokemonEndpoint(context.Background(), nil)
	if err != nil {
		return err
	}

	pr, ok := res.(PokemonResponse)
	if !ok {
		return errors.New("Response not of type PokemonResponse")
	}

	for _, v := range pr.Cards {
		println(v.ID)
	}

	//save data to cache

	return nil
}

type PokemonResponse struct {
	Cards []card `json:"cards,omitempty"`
}

type card struct {
	ID string `json:"id,omitempty"`
}
