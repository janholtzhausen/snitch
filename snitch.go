package main
import (
   "context"
   "fmt"
   "github.com/kubemq-io/kubemq-go"
   "log"
   "os"
   "bufio"
   "strings"
   "io"
   "strconv"
)

 type Config map[string]string

 func ReadConfig(filename string) (Config, error) {
     // init with some bogus data
 	config := Config{
 		"grpcserver":    "localhost",
 		"grpcport":      "50000",
 		"clientid":      "hello-world-subscriber",
 		"channelname":   "hello-world",
 	}
 	if len(filename) == 0 {
 		return config, nil
 	}
 	file, err := os.Open(filename)
 	if err != nil {
 		return nil, err
 	}
 	defer file.Close()
 	
 	reader := bufio.NewReader(file)
 	
 	for {
 		line, err := reader.ReadString('\n')
 		
 		// check if the line has = sign
             // and process the line. Ignore the rest.
 		if equal := strings.Index(line, "="); equal >= 0 {
 			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
 				value := ""
 				if len(line) > equal {
 					value = strings.TrimSpace(line[equal+1:])
 				}
                             // assign the config map
 				config[key] = value
 			}
 		}
 		if err == io.EOF {
 			break
 		}
 		if err != nil {
 			return nil, err
 		}
 	}
 	return config, nil
 }


func main() {

   config, err := ReadConfig(`/etc/snitch/snitch.conf`)
   if err != nil {
           fmt.Println(err)
   }

   grpcserver :=  config["grpcserver"]  
   grpcport :=  config["grpcport"]   
   clientid :=  config["clientid"]   
   channelname :=  config["channelname"]  

   intport, err := strconv.Atoi(grpcport) 

   ctx, cancel := context.WithCancel(context.Background())
   defer cancel()
   client, err := kubemq.NewClient(ctx,
      kubemq.WithAddress(grpcserver, intport),
      kubemq.WithClientId(clientid),
      kubemq.WithTransportType(kubemq.TransportTypeGRPC))
   if err != nil {
      log.Fatal(err)
   }
   defer client.Close()
   channelName := channelname
   errCh := make(chan error)
   eventsCh, err := client.SubscribeToEvents(ctx, channelName, "", errCh)
   if err != nil {
      log.Fatal(err)
      return

   }
   for {
      select {
      case err := <-errCh:
         log.Fatal(err)
         return
      case event, more := <-eventsCh:
         if !more {
            fmt.Println("Event Received, done")
            return
         }
         log.Printf("Event Received:\nEventID: %s\nChannel: %s\nMetadata: %s\nBody: %s\n", event.Id, event.Channel, event.Metadata, event.Body)
      case <-ctx.Done():
         return
      }
   }
}
