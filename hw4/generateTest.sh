#!/bin/bash

## generateTest.sh
##
## @author Tanya L. Crenshaw
## @since 10/25/2014
##
## This bash script creates a file that contains a collection of
## randomnly named companies formatted according to the input file
## rules of HW4.  For example,
## 
## * Linelotace Electrics
## No description
## http://Linelotace.com
## No address available
## None
## Portland
## OR
## 97203
## 45.521201
## -122.680772
##
## Example use: ./generateTest.sh 100 > 100.txt
##
## The above example use generates a list of 100 randmonly named
## companies and saves the result in a file called 100.txt
##

## Get the total number of companies that are to be generated from 
## the command line argument.
total=$1

## If the command line argument is unset, bail.
if [ -z "$1" ]; then
    echo "usage: ./generateTest.sh <number>"
    exit
fi

# Create two arrays of words that sound like international technobabble.
firstArray=(Nitro Volitude Black Dingy Zerfinnix Isfix Linelotace First Primer Second Segundo Biz Initech Pedestrian Swanson Betts Torus Belvin International Mega Giga Femto Nano Sustainable Palo Ace Tech Over Micro Fan Heroin Puppet Groove Candy White Transitron ACME Beta Alpha Massive Randomo Grande Massivo Dulce Computadores Tekeda SOL SERPA Gravify Cybria Amor Soompi Arena Super Nosotros Cellular Najo Culturo Pasemo Visteria Yemos Aviaty Turah Vipori Jiaro Rhini AQBY Calgara Aborio Globexo Sister Brother Mama Papa)

secondArray=(LLC CORP Associates Tech Camp Oxi Street Ware Laboratory DotCom Plus Circle Electrics Duo High Scale Zen House Vex Green Alto Zine Enterprise Oligarchs Warehouse Zap City State Dynamic Biz Corp llc) 

alphabet=(A B C D E F G H I J K L M N O P Q R S T U V W X Y Z)

# Calculate the sizes of the arrays
firstSize=${#firstArray[@]}
secondSize=${#secondArray[@]}


entry=0
while [ $entry -lt $total ]; do

    # flip a coin.  Either I'll choose a name randomly from the
    # first array, or I'll just make a three-letter acronym.
    ran=$RANDOM

    # if the random number is even, choose a name from
    # the first array.
    if [ $(( $ran % 2 )) -eq 0 ]
    then
	# Get the random number to be somewhere within the
	# size of the first array
	index=$ran%$firstSize
	name=${firstArray[$index]}
    else
	# Make a three-letter acronym
	counter=0
	name=  # Make an empty string
	while [ $counter -lt 3 ]; do
	    let ran=$RANDOM
	    index=$(($ran%26))
	    letter="${alphabet[$index]}"
	    name+=$letter
	    let counter=counter+1
	done    
    fi
    
    ## choose the second part of the name
    index=$ran%$secondSize
    ending=${secondArray[$index]}


    ## Produce the entry output
    echo "* $name $ending"
    echo "No description"
    echo "http://$name.com"
    echo "No address available"
    echo "None"
    echo "Portland"
    echo "OR"
    echo "97203"
    echo "45.521201"
    echo "-122.680772"
    echo ""

    ## increment entry counter
    let entry=entry+1
done
