package main

import(
  "fmt"
  "log"
  "github.com/nlopes/slack"
  "strings"
  "strconv"
  // "github.com/mongodb/mongo-go-driver/mongo"
  // "github.com/mongodb/mongo-go-driver/bson"
  // "github.com/mongodb/mongo-go-driver/mongo/options"
)

type SlackListener struct{
  client *slack.Client
  botID string
  channelID string
}

//Crear conexion a mongo


// ListenAndResponse escucha eventos de Slack y responde
// Mensajes particulares, y responde a través de un boton

func (s *SlackListener) ListenAndResponse(){
  rtm := s.client.NewRTM()

  //Conexion con los eventos de Slack
  go rtm.ManageConnection()

  // Manejo de los eventos de slack

  for msg := range rtm.IncomingEvents{
    switch ev := msg.Data.(type){
    case *slack.MessageEvent:
      if err := s.handleMessageEvent(ev); err != nil{
        log.Printf("[ERROR] Failed to handle message: %s", err)
      }
    }
  }
}

func (s *SlackListener) handleMessageEvent(ev *slack.MessageEvent) error {
  // Manejo de los eventos de Mensajes.

  //Solo responder cuando el bot es mencionado.
  if !strings.HasPrefix(ev.Msg.Text, fmt.Sprintf("<@%s> ", s.botID)){
    return nil
  }

  //Parseo del mensaje
  m := strings.Split(strings.TrimSpace(ev.Msg.Text), " ")[1:]
  if len(m) == 0 || m[0] != "cpn"{
    return fmt.Errorf("Mensaje invalido")
  }
  i, err := strconv.Atoi(m[1]); if err != nil{
    return fmt.Errorf("CPN invalido")
  }
  fmt.Println(i)


  attachment := slack.Attachment{
    Text: "¿que te gustaria hacer?",
    Color: "#f9a41b",
    CallbackID: "inicio",
    Actions: []slack.AttachmentAction{
      {
        Name: "cpn_info",
        Type: "select",
        Text: "Elige una",
        Options: []slack.AttachmentActionOption{
          {
            Text: "Datos",
            Value: "1",
          },
          {
            Text: "Cpe",
            Value: "2",
          },
        },
      },
    },
  }

  if _, _, err := s.client.PostMessage(ev.Channel, slack.MsgOptionText("Hola!", false), slack.MsgOptionAttachments(attachment)); err != nil{
    return fmt.Errorf("Failed to post message: %s", err)
  }
  return nil
}
//
