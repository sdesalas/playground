FROM maven:3.6.0-jdk-8-alpine

# download dependencies
ADD pom.xml /
RUN mvn verify clean

# build
ADD . /
RUN mvn package

# copy to bare jdk box
FROM openjdk:8-jdk-alpine
WORKDIR /root/
# → this relies on your MAVEN package command building *-jar-with-dependencies.jar with a single match
COPY --from=0 /target/*-jar-with-dependencies.jar app.jar 
ENTRYPOINT ["java","-Djava.security.egd=file:/dev/./urandom","-jar","./app.jar"]
