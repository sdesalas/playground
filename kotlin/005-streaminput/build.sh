#! /bin/bash

kotlinc src/**/*.kt -d build

kotlin -cp build project005.MainKt 5
