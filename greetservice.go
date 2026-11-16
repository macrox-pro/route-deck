// Command route-deck is the RouteDeck desktop shell.
package main

type GreetService struct{}

func (g *GreetService) Greet(name string) string {
	return "Hello " + name + "!"
}
