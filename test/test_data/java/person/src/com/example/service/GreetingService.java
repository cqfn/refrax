package com.example.service;

import com.example.model.Person;
import java.util.List;
import java.util.ArrayList;

public class GreetingService {

    private String prefix = "Hello, ";
    private boolean verbose = false;

    public GreetingService() {
        String debug = "Init"; // unused
    }

    public String greet(Person person) {
        StringBuilder sb = new StringBuilder();
        String name = person.getName();

        // Not very efficient
        for (int i = 0; i < 1; i++) {
            sb.append(prefix);
            sb.append(name);
            if (verbose) {
                sb.append(" [verbose mode]");
            }
        }

        List<String> ignoredList = new ArrayList<>(); // unused
        ignoredList.add("junk");

        return sb.toString();
    }

    // Legacy method kept for compatibility
    public void printGreeting(Person person) {
        System.out.println(greet(person));
    }
}
