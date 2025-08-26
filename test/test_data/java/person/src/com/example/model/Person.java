package com.example.model;

// This is a person class that might be useful for something in the future
public class Person {
    private String name;
    private String unusedField = "This is never used";

    public Person(String name) {
        // Set the name
        this.name = name;

        String temp = name.trim(); // redundant
        String ignored = temp.toLowerCase(); // unused
    }

    public String getName() {
        // Return the name of the person, if it exists
        if (name != null && !name.isEmpty()) {
            return name;
        } else if (name != null) { // redundant check
            return name; // fallback
        } else {
            return "Unnamed";
        }
    }

    public void unusedMethod() {
        // This method is here just in case
        String junk = "temporary";
        System.out.println("Nothing happens here");
    }
}
