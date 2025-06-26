# Handy commands

```
> kotlinc -help

> kotlinc -version

> kotlinc -e 'println("hello")'

```

For Building single file

```
> echo 'fun main() { println("Hello") }' > Hello.kt
> kotlinc Hello.kt
> kotlin HelloKt
Hello
```

For building src dir with files

```
kotlinc src/**/*.kt -d build
```

For running with classpath 

```
kotlin -cp build project003.GuessKt <num>
```


