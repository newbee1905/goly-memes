package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dotenv "github.com/Netflix/go-env"
)

type Env struct {
	Db                      string `env:"DB"`
	DiscordWebHook          string `env:"DiscordWebHook"`
	DiscordWebHookAvatarUrl string `env:"DiscordWebHookAvatarUrl"`
}

type DiscordHookEmbedAuthor struct {
	Name string `json:"name"`
}

type DiscordHookEmbedImage struct {
	Url string `json:"url"`
}

type DiscordHookEmbed struct {
	Colour int                    `json:"color"`
	Author DiscordHookEmbedAuthor `json:"author"`
	Image  DiscordHookEmbedImage  `json:"image"`
}

type DiscordHookMsg struct {
	Content   string             `json:"content"`
	Embeds    []DiscordHookEmbed `json:"embeds"`
	AvatarUrl string             `json:"avatar_url"`
	Username  string             `json:"username"`
}

type RedditData struct {
	Kind string `json:"kind"`
	Data struct {
		Children []struct {
			Data struct {
				Title string `json:"title"`
				Url   string `json:"url_overridden_by_dest"`
			} `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

func main() {
	var env Env
	_, err := dotenv.UnmarshalFromEnviron(&env)
	if err != nil {
		log.Fatal(err)
	}

	var redditData RedditData

	req, err := http.Get("https://reddit.com/r/memes/top.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.NewDecoder(req.Body).Decode(&redditData)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Println(redditData.Data.Children)

	redditEmbed := DiscordHookEmbed{
		0,
		DiscordHookEmbedAuthor{redditData.Data.Children[0].Data.Title},
		DiscordHookEmbedImage{redditData.Data.Children[0].Data.Url},
	}

	msg := DiscordHookMsg{
		"# Daily Reddit Meme",
		[]DiscordHookEmbed{redditEmbed},
		env.DiscordWebHookAvatarUrl,
		"Memes Daily",
	}
	fmt.Println(msg)
	b, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(b))
	bodyReader := bytes.NewReader(b)

	fmt.Println(env.DiscordWebHook)
	req, err = http.Post(env.DiscordWebHook, "application/json; charset=UTF-8", bodyReader)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(req)
}
