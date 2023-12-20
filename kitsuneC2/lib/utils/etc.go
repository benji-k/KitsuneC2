package utils

import (
	"fmt"
	"math/rand"
)

var adjectives []string = []string{
	"Greedy",
	"Hysterical",
	"Ambitious",
	"Flagrant",
	"Smelly",
	"Tidy",
	"Thinkable",
	"Obscene",
	"Scandalous",
	"Thankful",
	"Eight",
	"Slim",
	"Tranquil",
	"Fascinated",
	"Heavenly",
	"Erratic",
	"High",
	"Sulky",
	"Encouraging",
	"Superficial",
}

var nouns []string = []string{
	"Office",
	"Police",
	"Love",
	"Housing",
	"Beer",
	"Contract",
	"Blood",
	"Contribution",
	"Mall",
	"Explanation",
	"Throat",
	"Pollution",
	"Client",
	"Emphasis",
	"Promotion",
	"Road",
	"Feedback",
	"Physics",
	"Method",
	"Teaching",
}

func GenerateRandomName() string {
	adjective := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	return adjective + noun
}

func PrintBanner() {
	fmt.Println(`
                       @@@@@                                                  @@@@                                                
                       @@@@@@@@                                             @@@@@@                                                
                       @@@@@@@@@@                                         @@@@@@@@                                                
                       @@@@@@@@@@@@                                    @@@@@@@@@@@                                                
                       @@@@@@@@@@@@@@@                               @@@@@@@@@@@@@                                                
                       @@@@@@@@@@ @@@@@@                          @@@@@@@@@@@@@@@@                                                
                       @@@@@@@@@@@   @@@@@                      @@@@@@ @@@@@@@@@@@                                                
                       @@@@@@@@@@@     @@@@@@@@@@@@@@@@@@@@@@@@@@@@   @@@@@@@@ @@@                                                
                       @@@@ @@@ @@@    @@@@@@@@@@@@@@@@@@@@@@@@@@    @@@@@@@@ @@@@                                                
                       @@@@ @@@@@@@@@@@@@@@        @@@        @@@@@@@@@@@@@@  @@@@                                                
                       @@@@  @@@@@@@@@@@           @@@          @@@@@@@@@@@   @@@@                                                
                       @@@@   @@@@@@@              @@@             @@@@@@@    @@@@                                                
                       @@@@@@@@@@                  @@@                @@@@@@  @@@@                                                
                       @@@@@@@                     @@@                  @@@@@@@@@@                                                
                       @@@@                        @@@                     @@@@@@@                                                
                       @@@@@                       @@@                        @@@@                                                
                       @@@@@@                      @@@                       @@@@@@                                               
                       @@@@@@@                     @@@                      @@@@@@@                                               
                      @@@@ @@@@@@@                 @@@                     @@@@ @@@                                               
                      @@@@   @@@@@@@@              @@@                   @@@@@  @@@@                                              
                      @@@@    @@@@@@@@@@           @@@                @@@@@@    @@@@                                              
                      @@@      @@@@ @@@@@@@        @@@            @@@@@@@@@      @@@                                              
                      @@@       @@@@    @@@@@@     @@@       @@@@@@@@@@@@@       @@@                                              
                      @@@        @@@@      @@@@@@@ @@@   @@@@@@@@@   @@@@       @@@@@                                             
                      @@@@@       @@@@        @@@@@@@@@@@@@@@@     @@@@       @@@@@                                               
                        @@@@@      @@@@          @@@@@@@@@        @@@@      @@@@@                                                 
                          @@@@@@    @@@@                         @@@@     @@@@@                                                   
                             @@@@@    @@@@                      @@@@    @@@@@                                                     
                              @@@@@@   @@@@                    @@@@   @@@@@                                                     
                                @@@@@@  @@@@                 @@@@@  @@@@@                                                  
                                  @@@@@@ @@@@               @@@@ @@@@@@                                                 
                                     @@@@@@@@@             @@@@@@@@@@                                                  
                                       @@@@@@@@    @@@@   @@@@@@@@@                                                    
                                         @@@@@@@@@@@@@@@@@@@@@@@@                                                    
                                            @@@@@@@    @@@@@@@                                                         
                                              @@@@@    @@@@@                                                          
                                                @@@@@@@@@@                                          
                                                  @@@@@@      
	`)
}
