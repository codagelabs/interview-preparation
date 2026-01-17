import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.annotation.JsonSerialize;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import org.json.JSONObject;
import org.json.JSONArray;

import java.time.LocalDateTime;
import java.util.List;
import java.util.ArrayList;

/**
 * Comprehensive JSON Serialization Examples in Java
 * Demonstrates different libraries: Jackson, Gson, and org.json
 */
public class JsonSerializationExample {
    
    // Sample data class for serialization
    static class Person {
        private String name;
        private int age;
        private String email;
        private List<String> hobbies;
        private LocalDateTime createdAt;
        
        // Default constructor (required for some libraries)
        public Person() {}
        
        // Constructor
        public Person(String name, int age, String email, List<String> hobbies) {
            this.name = name;
            this.age = age;
            this.email = email;
            this.hobbies = hobbies;
            this.createdAt = LocalDateTime.now();
        }
        
        // Getters and Setters
        public String getName() { return name; }
        public void setName(String name) { this.name = name; }
        
        public int getAge() { return age; }
        public void setAge(int age) { this.age = age; }
        
        public String getEmail() { return email; }
        public void setEmail(String email) { this.email = email; }
        
        public List<String> getHobbies() { return hobbies; }
        public void setHobbies(List<String> hobbies) { this.hobbies = hobbies; }
        
        public LocalDateTime getCreatedAt() { return createdAt; }
        public void setCreatedAt(LocalDateTime createdAt) { this.createdAt = createdAt; }
        
        @Override
        public String toString() {
            return "Person{name='" + name + "', age=" + age + ", email='" + email + 
                   "', hobbies=" + hobbies + ", createdAt=" + createdAt + "}";
        }
    }
    
    public static void main(String[] args) {
        // Create sample data
        List<String> hobbies = new ArrayList<>();
        hobbies.add("Reading");
        hobbies.add("Coding");
        hobbies.add("Gaming");
        
        Person person = new Person("John Doe", 30, "john@example.com", hobbies);
        
        System.out.println("=== JSON Serialization Examples ===\n");
        
        // 1. Jackson (Most popular)
        jacksonExample(person);
        
        // 2. Gson (Google's library)
        gsonExample(person);
        
        // 3. org.json (Simple, no dependencies)
        orgJsonExample(person);
        
        // 4. Serializing Collections
        serializeCollections();
        
        // 5. Custom Serialization
        customSerializationExample();
    }
    
    /**
     * Jackson JSON Serialization
     * Most feature-rich and widely used
     */
    public static void jacksonExample(Person person) {
        System.out.println("1. JACKSON JSON SERIALIZATION");
        System.out.println("==============================");
        
        try {
            ObjectMapper mapper = new ObjectMapper();
            
            // Serialize object to JSON string
            String jsonString = mapper.writeValueAsString(person);
            System.out.println("Serialized JSON:");
            System.out.println(jsonString);
            
            // Pretty print JSON
            String prettyJson = mapper.writerWithDefaultPrettyPrinter().writeValueAsString(person);
            System.out.println("\nPretty printed JSON:");
            System.out.println(prettyJson);
            
            // Deserialize JSON back to object
            Person deserializedPerson = mapper.readValue(jsonString, Person.class);
            System.out.println("\nDeserialized object:");
            System.out.println(deserializedPerson);
            
        } catch (Exception e) {
            System.err.println("Jackson error: " + e.getMessage());
        }
        
        System.out.println("\n" + "=".repeat(50) + "\n");
    }
    
    /**
     * Gson JSON Serialization
     * Google's library, simple and lightweight
     */
    public static void gsonExample(Person person) {
        System.out.println("2. GSON JSON SERIALIZATION");
        System.out.println("===========================");
        
        try {
            Gson gson = new GsonBuilder()
                .setPrettyPrinting()
                .create();
            
            // Serialize object to JSON string
            String jsonString = gson.toJson(person);
            System.out.println("Serialized JSON:");
            System.out.println(jsonString);
            
            // Deserialize JSON back to object
            Person deserializedPerson = gson.fromJson(jsonString, Person.class);
            System.out.println("\nDeserialized object:");
            System.out.println(deserializedPerson);
            
        } catch (Exception e) {
            System.err.println("Gson error: " + e.getMessage());
        }
        
        System.out.println("\n" + "=".repeat(50) + "\n");
    }
    
    /**
     * org.json JSON Serialization
     * Simple library with no external dependencies
     */
    public static void orgJsonExample(Person person) {
        System.out.println("3. ORG.JSON JSON SERIALIZATION");
        System.out.println("===============================");
        
        try {
            // Manual JSON object creation
            JSONObject jsonObject = new JSONObject();
            jsonObject.put("name", person.getName());
            jsonObject.put("age", person.getAge());
            jsonObject.put("email", person.getEmail());
            
            // Add hobbies array
            JSONArray hobbiesArray = new JSONArray();
            for (String hobby : person.getHobbies()) {
                hobbiesArray.put(hobby);
            }
            jsonObject.put("hobbies", hobbiesArray);
            jsonObject.put("createdAt", person.getCreatedAt().toString());
            
            // Serialize to JSON string
            String jsonString = jsonObject.toString(2); // Pretty print with 2 spaces
            System.out.println("Serialized JSON:");
            System.out.println(jsonString);
            
            // Parse JSON back to object (manual)
            JSONObject parsedJson = new JSONObject(jsonString);
            Person deserializedPerson = new Person();
            deserializedPerson.setName(parsedJson.getString("name"));
            deserializedPerson.setAge(parsedJson.getInt("age"));
            deserializedPerson.setEmail(parsedJson.getString("email"));
            
            // Parse hobbies array
            JSONArray parsedHobbies = parsedJson.getJSONArray("hobbies");
            List<String> hobbiesList = new ArrayList<>();
            for (int i = 0; i < parsedHobbies.length(); i++) {
                hobbiesList.add(parsedHobbies.getString(i));
            }
            deserializedPerson.setHobbies(hobbiesList);
            
            System.out.println("\nDeserialized object:");
            System.out.println(deserializedPerson);
            
        } catch (Exception e) {
            System.err.println("org.json error: " + e.getMessage());
        }
        
        System.out.println("\n" + "=".repeat(50) + "\n");
    }
    
    /**
     * Serializing Collections
     */
    public static void serializeCollections() {
        System.out.println("4. SERIALIZING COLLECTIONS");
        System.out.println("===========================");
        
        try {
            // Create a list of persons
            List<Person> people = new ArrayList<>();
            people.add(new Person("Alice", 25, "alice@example.com", List.of("Swimming", "Dancing")));
            people.add(new Person("Bob", 35, "bob@example.com", List.of("Cooking", "Photography")));
            
            // Jackson - Serialize list
            ObjectMapper mapper = new ObjectMapper();
            String jsonArray = mapper.writerWithDefaultPrettyPrinter().writeValueAsString(people);
            System.out.println("Serialized List of Persons:");
            System.out.println(jsonArray);
            
            // Deserialize list back
            List<Person> deserializedPeople = mapper.readValue(jsonArray, 
                mapper.getTypeFactory().constructCollectionType(List.class, Person.class));
            System.out.println("\nDeserialized List:");
            deserializedPeople.forEach(System.out::println);
            
        } catch (Exception e) {
            System.err.println("Collection serialization error: " + e.getMessage());
        }
        
        System.out.println("\n" + "=".repeat(50) + "\n");
    }
    
    /**
     * Custom Serialization Example
     */
    public static void customSerializationExample() {
        System.out.println("5. CUSTOM SERIALIZATION");
        System.out.println("=======================");
        
        try {
            // Using Jackson with custom annotations
            ObjectMapper mapper = new ObjectMapper();
            
            // Create a simple object for custom serialization
            SimpleData data = new SimpleData("Test", 42, true);
            
            String json = mapper.writerWithDefaultPrettyPrinter().writeValueAsString(data);
            System.out.println("Custom serialized JSON:");
            System.out.println(json);
            
        } catch (Exception e) {
            System.err.println("Custom serialization error: " + e.getMessage());
        }
    }
    
    // Simple class for custom serialization example
    static class SimpleData {
        private String name;
        private int value;
        private boolean active;
        
        public SimpleData() {}
        
        public SimpleData(String name, int value, boolean active) {
            this.name = name;
            this.value = value;
            this.active = active;
        }
        
        // Getters and Setters
        public String getName() { return name; }
        public void setName(String name) { this.name = name; }
        
        public int getValue() { return value; }
        public void setValue(int value) { this.value = value; }
        
        public boolean isActive() { return active; }
        public void setActive(boolean active) { this.active = active; }
    }
}

