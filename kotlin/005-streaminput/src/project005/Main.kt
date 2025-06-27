package project005

import java.util.Scanner

fun main(args: Array<String>) {
  val number = (Math.random() * 10).toInt()
  println("Guess the number")

  Scanner(System.`in`).use { input -> 
    do {
      val guess = input.nextInt()
      if (guess != number) {
        println("$guess is NOT the right number. Try again!")
      } else {
        println("Wohoo! $guess is the right number!")
      }
    } while (guess != number)
  }
}