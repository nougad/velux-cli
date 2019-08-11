package main

import "os"
import "fmt"
import "flag"
import "strings"

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join([]string(*i), ",")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("usage: velux-cli <command> [<args>]")
		fmt.Println("The most commonly used commands are: ")
		fmt.Println(" print   Shows status")
		fmt.Println(" dump    Writs json into file")
		fmt.Println(" moveShutters moves shutters")
		return
	}

	switch os.Args[1] {
	case "print":
		printCommand := flag.NewFlagSet("print", flag.ExitOnError)
		tokenpath := printCommand.String("tokenfile", "/openhab/conf/token.json", "file with access token")
		printCommand.Parse(os.Args[2:])
		state := fetchData(*tokenpath)
		PrintStatus(state)
	case "dump":
		dumpCommand := flag.NewFlagSet("print", flag.ExitOnError)
		tokenpath := dumpCommand.String("tokenfile", "/openhab/conf/token.json", "file with access token")
		jsonout := dumpCommand.String("outfile", "-", "file writing output to")
		dumpCommand.Parse(os.Args[2:])
		state := fetchData(*tokenpath)
		DumpJSON(state, *jsonout)
	case "moveShutters":
		cmd := flag.NewFlagSet("moveShutter", flag.ExitOnError)
		tokenpath := cmd.String("tokenfile", "/openhab/conf/token.json", "file with access token")
		position := cmd.Int("pos", 0, "move shutter to position")
		var shutters arrayFlags
		cmd.Var(&shutters, "shutters", "which sutters to control")
		cmd.Parse(os.Args[2:])

		state := fetchData(*tokenpath)
		Move(state, shutters, int64(*position))
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}
