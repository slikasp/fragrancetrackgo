package handlers

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/slikasp/fragrancetrackgo/internal/config"
	"github.com/slikasp/fragrancetrackgo/internal/database"
)

func FragranceAdd(s *config.State, brand, name string) error {
	// TODO: maybe capitalise first letters? need to cover all edge cases (l'homme and so on)
	// set all to lowercase
	brand = strings.ToLower(brand)
	name = strings.ToLower(name)

	exits, _ := s.Db.GetFragrance(context.Background(), database.GetFragranceParams{
		Brand: brand,
		Name:  name,
	})
	if exits.Name == name && exits.Brand == brand {
		return fmt.Errorf("Fragrance %s - %s already exists with ID %d.\n", brand, name, exits.ID)
	}

	newFrag, err := s.Db.AddFragrance(context.Background(), database.AddFragranceParams{
		Brand: brand,
		Name:  name,
	})
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not add fragrance %s - %s to database.\n", brand, name)
	}

	log.Printf("Fragrance added: %s - %s\n", newFrag.Brand, newFrag.Name)

	return nil
}

func FragranceRemove(s *config.State, id int32) error {
	removedFrag, err := s.Db.RemoveFragrance(context.Background(), id)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not remove fragrance with ID %d from database.\n", id)
	}

	log.Printf("Fragrance removed: %v\n", removedFrag)

	return nil
}

func FragranceUpdate(s *config.State, id int32, brand, name string) error {
	updatedFrag, err := s.Db.UpdateFragrance(context.Background(), database.UpdateFragranceParams{
		ID:    id,
		Brand: brand,
		Name:  name,
	})
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not update fragrance with ID %d\n", id)
	}

	log.Printf("Fragrance updated: %v\n", updatedFrag)

	return nil
}

func FragranceList(s *config.State) error {
	frags, err := s.Db.GetFragrances(context.Background())
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not get existing users.\n")
	}

	for _, frag := range frags {
		fmt.Printf("%d. %s - %s\n", frag.ID, frag.Brand, frag.Name)
	}

	return nil
}
