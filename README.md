# Livos - liquid federate voting system:

Semester project fall 2021. Done at EPFL with the DEDIS lab.
Project directed by N.Kocher and realised by Guillaume Tabard and Etienne Boisson.

This project aims to build a proof of concept that liquid Democracy can indeed work in the real life.
Liquid Democracy is the democracy model fit for today’s society. The technology is ready, the only thing missing is more effort into working on concrete implementations that research this area in more detail. Most importantly, we need to determine which models are applicable for the actual governance (directive or administrative) of a country. The urge of having an implementation of this LD system is a huge strength of our project, answering this problem.

That's why you can find a backend implementation of all the liquid system with a delegation and a vote process in the votings files. 

In the web package you can run in local(http://localhost:9000/) our website. Just write in your terminal : "go run mod.go". The webpackage contains all the frontend implementation, you also need to have installed viz.js in order to be able to displays the graph of the state of the vote. 
You can create Voting rooms, vote by clicking on the room number, manage the voting session with the parameters, close the room to see the results of the election.

You can also run simulations, just go with the terminal command: cd .\simulation\, inside the package simulation. Then the file simulation_test.go allows you to choose which simulation you want to run. Then write in the terminal "go test", this will generate an output text file that needs to be fed to graphviz (you need to have imported Graphviz before) with the command : "dot -Tpdf outpuSimulation -o outputSimulation.pdf". This will generate a pdf representing your graph (you therefore needs a pdf reader).
These simulations aims to proove that there exist a difference beetween a democracy with or without liquidity. 

Check out our results there in the Report: 