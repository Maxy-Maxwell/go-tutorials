package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Maxy-Maxwell/go-tutorials/Gator/internal/database"
	rss "github.com/Maxy-Maxwell/go-tutorials/Gator/internal/rss"

	"github.com/google/uuid"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		// Get the user record (currently logged in) so we have the user_id
		username := sql.NullString{
			String: s.cfg.Current_user_name,
			Valid:  true,
		}

		user, err := s.db.GetUser(context.Background(), username)
		if err != nil {
			return fmt.Errorf("Failed to retrieve user with name: %v", username)
		}

		if !user.ID.Valid {
			return fmt.Errorf("Failed to find user: %v.", username)
		}

		return handler(s, cmd, user)
	}

}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("The login handler expects a single argument - username")
	}

	username := sql.NullString{
		String: cmd.arguments[0],
		Valid:  true,
	}

	if _, err := s.db.GetUser(context.Background(), username); err != nil {
		return err
	}

	err := s.cfg.SetUser(cmd.arguments[0])
	if err != nil {
		return err
	}

	fmt.Println("User logged in:", cmd.arguments[0])

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("The register handler expects a single argument - username")
	}

	currentTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	userId := uuid.NullUUID{
		UUID:  uuid.New(),
		Valid: true,
	}

	username := sql.NullString{
		String: cmd.arguments[0],
		Valid:  true,
	}

	if _, err := s.db.CreateUser(
		context.Background(),
		database.CreateUserParams{
			ID:        userId,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			Name:      username,
		},
	); err != nil {
		return err
	}

	fmt.Println("User created:", username.String)

	return handlerLogin(s, command{name: "login", arguments: []string{username.String}})
}

func handlerReset(s *state, cmd command) error {
	if err := s.db.DeleteUser(context.Background()); err != nil {
		return err
	}

	fmt.Println("Users deleted")

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())

	if err != nil {
		return err
	}

	for _, user := range users {
		if user.Name.String == s.cfg.Current_user_name {
			fmt.Printf("* %v (current)\n", user.Name.String)
		} else {
			fmt.Printf("* %v\n", user.Name.String)
		}
	}

	return nil
}

func handlerAgg(s *state, cmd command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}

	fmt.Printf("Returned feed: %v\n", *feed)

	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("The Add Feed handler expects a two arguments - name, url")
	}

	feedId := uuid.NullUUID{
		UUID:  uuid.New(),
		Valid: true,
	}

	feedName := sql.NullString{
		String: cmd.arguments[0],
		Valid:  true,
	}

	url := sql.NullString{
		String: cmd.arguments[1],
		Valid:  true,
	}

	if _, err := s.db.CreateFeed(
		context.Background(),
		database.CreateFeedParams{
			ID:     feedId,
			Name:   feedName,
			Url:    url,
			UserID: user.ID,
		},
	); err != nil {
		return err
	}

	// Auto-follow newly created feed
	fmt.Println("Auto-following feed...")
	followArgs := command{
		name:      "follow",
		arguments: []string{url.String},
	}
	return handlerFollow(s, followArgs, user)

}

func handlerFeeds(s *state, cmd command) error {
	feedsWithUsers, err := s.db.GetFeedsWithUserDetails(context.Background())
	if err != nil {
		return err
	}

	for _, feedWithUser := range feedsWithUsers {
		fmt.Printf(
			"Name: %v, URL: %v, Created By: %v\n",
			feedWithUser.Name.String,
			feedWithUser.Url.String,
			feedWithUser.Username.String,
		)
	}

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		return fmt.Errorf("The Follow handler expects a one argument - url")
	}

	fmt.Println("Url arg: ", cmd.arguments[0])

	url := sql.NullString{
		String: cmd.arguments[0],
		Valid:  true,
	}

	// Get the feed record so we have the feed_id
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	// Prepare other new FeedFollow fields
	currentTime := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}

	feedFollowId := uuid.NullUUID{
		UUID:  uuid.New(),
		Valid: true,
	}

	// Create the new feed_follow record
	feedFollow, err := s.db.GetFeedFollow(
		context.Background(),
		database.GetFeedFollowParams{
			ID:        feedFollowId,
			CreatedAt: currentTime,
			UpdatedAt: currentTime,
			UserID:    user.ID,
			FeedID:    feed.ID,
		},
	)

	fmt.Println("FeedFollow: ", feedFollow)

	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	feeds, err := s.db.GetFeedFollowsForUser(
		context.Background(),
		user.ID,
	)

	if err != nil {
		return err
	}

	for _, feed := range feeds {
		fmt.Println("Feed: ", feed.FeedsName)
	}

	return nil
}
