package main

import(
  "log"
  "net/http"
  "github.com/nlopes/slack"
)

type envConfig struct{
  Port string `default:"6767"`
  BotToken string `default:"xoxb-211403098598-421812270979-xRlm4BxYQqw2we9vqdTmgD1R"`
  VerificationToken string `default:"JFMdRmB5YbsdtLgWsOWAnWvQ"`
  BotID string `default:"UCDPW7YUT"`
  ChannelID string `default:"CCTU1PSJU"`
}

func main(){
  //var env envConfig

  log.Printf("[INFO] Bot slack listening")
  client := slack.New("xoxb-211403098598-421812270979-xRlm4BxYQqw2we9vqdTmgD1R")
  slackListener := &SlackListener{
    client: client,
    botID: "UCDPW7YUT",
    channelID: "CCTU1PSJU",
  }
  go slackListener.ListenAndResponse()

  ////////
  // Routes
  ///////
  router := NewRouter()
  log.Printf("[INFO] Server listening on PORT 6767")
  server := http.ListenAndServe(":6767", router)
  log.Fatal(server)
}
