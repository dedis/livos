# Livos - liquid federate voting system:

Semester project fall 2021. Done at EPFL with the DEDIS lab.
Project directed by N.Kocher and realised by Guillaume Tabard and Etienne Boisson.

## Abstract 
This project aims to build a proof of concept that liquid Democracy can indeed work in the a real life application and obtains some results about the impact of a liquid system compared to a more traditional one.
Liquid Democracy is the democracy model fit for today’s society. The technology is ready, the only thing missing is more effort into working on concrete implementations that researches this area in more detail. Most importantly, we need to determine which models are applicable for the actual governance (directive or administrative) of a country. The urge of having an implementation of this LD system is a huge strength of our project.

You can find a backend implementation of all the liquid system with a delegation and a vote process in the votings files. 

The server runs in the local port 9000 (http://localhost:9000/). 
To start it, run in the terminal 
>go run mod.go

## Website application
The webpackage contains all the frontend implementation.
![This is the Homepage](/web/images/HomepageEmpty.PNG)
![This is the Homepage with rooms](/web/images/HomePageManyRooms.PNG)

You can create Voting rooms, enter by clicking on the room number, manage the voting session to change the parameters and vote/delegate choosing the user and the candidates/referendum choice the user votes for or the voter the user delegates to.
![This is the Creation page](/web/images/CreationPage.PNG)
![This is the Voting room](/web/images/ElectionRoomEmpty.PNG)
![This is the Managing page](/web/images/ManageOneRoom.PNG)

In order to see the results of the election/referendum, you need close the room and enter it.
![This is the Result page](/web/images/ResultsOfRoom.PNG)

You also need to already have installed viz.js in order to be able to displays the graphs of the current state of the vote.
![This is the graph of the state of the vote](/web/images/GraphRenderedOnWebSite.PNG)

## Simulations
You can also run simulations. You need to be in the simulation directory : cd .\simulation.
The file simulation_test.go contains the simulations you can run. You can modify it's content to modify the number of simulations you want to run, or what type of voting session you want to experiment (election of referendum) and other parameters as the number of voters or candidates.
Then run the command in the terminal
>go test

If you run only one simulation, this will generate an output text file that needs to be fed to graphviz (you need to have imported Graphviz before) with the command :
> dot -Tpdf outpuSimulation -o outputSimulation.pdf
This will generate a pdf representing the graph of the vote.
![This is a example of the final graph of the vote](/web/images/Graph_best_faculty.PNG)

You can learn more about our project in our [Report](/LIVOSProjectReport.pdf)