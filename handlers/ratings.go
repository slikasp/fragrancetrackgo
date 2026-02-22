package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/slikasp/fragrancetrackgo/internal/config"
	"github.com/slikasp/fragrancetrackgo/internal/database"
)

// reuse this for Scores after migrations are remade
// how hard would it be to get suggestions in cli after pressing tab? maybe just go to web app straight away?

func RatingAdd(s *config.State, brand, name, comment string, rating int32) error {
	// TODO: maybe capitalise first letters? need to cover all edge cases (l'homme and so on)
	// set all to lowercase
	brand = strings.ToLower(brand)
	name = strings.ToLower(name)

	exits, _ := s.Users.GetRating(context.Background(), database.GetRatingParams{
		UserID: s.Cfg.CurrentUser,
		Brand:  brand,
		Name:   name,
	})
	if exits.Name == name && exits.Brand == brand {
		return fmt.Errorf("Fragrance %s - %s already exists with ID %d.\n", brand, name, exits.ID)
	}

	newFrag, err := s.Users.AddRating(context.Background(), database.AddRatingParams{
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

func RatingRemove(s *config.State, brand, name string) error {
	removedFrag, err := s.Users.RemoveRating(context.Background(), database.RemoveRatingParams{
		UserID: s.Cfg.CurrentUser,
		Brand:  brand,
		Name:   name,
	})
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not remove fragrance %s - %s from database.\n", brand, name)
	}

	log.Printf("Fragrance removed: %v\n", removedFrag)

	return nil
}

func RatingUpdate(s *config.State, brand, name string, comment sql.NullString, rating sql.NullInt32) error {
	updatedFrag, err := s.Users.UpdateRating(context.Background(), database.UpdateRatingParams{
		UserID:  s.Cfg.CurrentUser,
		Brand:   brand,
		Name:    name,
		Rating:  rating,
		Comment: comment,
	})
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not update fragrance %s 0 %s \n", brand, name)
	}

	log.Printf("Fragrance updated: %v\n", updatedFrag)

	return nil
}

func RatingList(s *config.State) error {
	frags, err := s.Users.GetRatings(context.Background(), s.Cfg.CurrentUser)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not get existing users.\n")
	}

	for _, frag := range frags {
		fmt.Printf("%d. %s - %s\n", frag.ID, frag.Brand, frag.Name)
	}

	return nil
}
