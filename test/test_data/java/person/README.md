# Testing Java Project

This Java project is designed for end-to-end (e2e) testing and serves as a refactoring example. All Java classes are located in the `src` directory.

To run the java program you can invoke the following command:


```sh
javac -d out src/com/example/model/Person.java \
           src/com/example/service/GreetingService.java \
           src/com/example/MainApp.java \
&& java -cp out com.example.MainApp
```
