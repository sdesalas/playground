package project004

fun main(args: Array<String>) {
  val secret = System.getenv("SECRET_NUMBER").toInt()
  if (args.size == 0) {
    println("incorrect, please input a number")
    return
  }
  val guess = args.first().toInt()

  if (guess == secret) println("thats it, you guessed right!")
  else println("whoops try again!")
}