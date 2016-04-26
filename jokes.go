package jokes

import (
    "log"
    "os"
    "io/ioutil"
    "encoding/json"
    "net/http"

    "github.com/itsabot/abot/shared/plugin"
    "github.com/itsabot/abot/shared/nlp"
    "github.com/itsabot/abot/shared/datatypes"
)

var p *dt.Plugin
var mashapeKey = os.Getenv("MASHAPE_KEY")

type APIResponse struct {
 Title string
 Joke string
 Category string
}

func init() {
    // Abot should route messages to this plugin that contain any combination
    // of the below words. The stems of the words below are used, so you don't
    // need to include duplicates (e.g. there's no need to include both "stock"
    // and "stocks"). Everything will be lowercased as well, so there's no
    // difference between "ETF" and "etf".
    trigger := &nlp.StructuredInput{
        Commands: []string{"tell"},
        Objects: []string{"joke"},
    }

    // Tell Abot how this plugin will respond to new conversations and follow-up
    // requests.
    fns := &dt.PluginFns{Run: Run, FollowUp: FollowUp}

    // Create the plugin.
    var err error
    pluginPath := "github.com/chuckwagoncomputing/abot-plugin-jokes"
    p, err = plugin.New(pluginPath, trigger, fns)
    if err != nil {
        log.Fatalln("building", err)
    }

    p.Vocab = dt.NewVocab(
        dt.VocabHandler{
            Fn: tellJoke,
            Trigger: &nlp.StructuredInput{
                Commands: []string{"tell"},
                Objects: []string{"joke"},
            },
        },
    )
}

// Abot calls Run the first time a user interacts with a plugin
func Run(in *dt.Msg) (string, error) {
    return FollowUp(in)
}

// Abot calls FollowUp every subsequent time a user interacts with the plugin
// as long as the messages hit this plugin consecutively. As soon as Abot sends
// a message for this user to a different plugin, this plugin's Run function
// will be called the next it's triggered.  This Run/FollowUp design allows us
// to reset a plugin's state when a user changes conversations.
func FollowUp(in *dt.Msg) (string, error) {
    return p.Vocab.HandleKeywords(in), nil
}

func tellJoke(in *dt.Msg) string {
 client := &http.Client{ }

 req, err := http.NewRequest("GET", "http://webknox-jokes.p.mashape.com/jokes/random?category=clean", nil)
 req.Header.Add("X-Mashape-Key", mashapeKey)
 req.Header.Add("Accept", "application/json")
 resp, err := client.Do(req)

 if (err == nil) {
  body, err := ioutil.ReadAll(resp.Body)
  if (err == nil) {
   var data APIResponse
   err := json.Unmarshal(body, &data)
   if (err == nil) {
    return data.Joke
   } else {
    return "Hmmm, That joke wasn't very funny."
   }
  } else {
   return "Hmmm, That joke didn't make any sense."
  }
 } else {
  return "Oops, That joke didn't make any sense."
 }
}
