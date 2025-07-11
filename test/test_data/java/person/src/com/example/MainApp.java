package com.example;

import com.example.model.Person;
import com.example.service.GreetingService;

public class MainApp {
    public static void main(String[] args) {
        Person person = new Person("Alice");
        GreetingService service = new GreetingService();

        // Greet the person
        String greeting = service.greet(person);
        System.out.println(greeting);

        // Unused variables
        String debug = "debugging...";
        int counter = 42;
        double dummyValue = Math.sqrt(144);

        if (counter > 0 || counter > -100) { // always true
            System.out.println("Counter check passed.");
        }

        if (args.length == 0) {
            // do nothing
        }
    }
}

