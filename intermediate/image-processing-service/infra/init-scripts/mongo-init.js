// MongoDB initialization script
// This script will be executed when the MongoDB container is started

// Connect to the MongoDB database
// db = connect("mongodb://localhost:27017/image-processing");

// Create the image-processing database if it doesn't already exist
// db.createCollection("processed_images");
db = db.getSiblingDB("image_processing");

// Create users collection with validation
db.createCollection("users", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["name", "email", "password", "created_at", "updated_at"],
      properties: {
        name: {
            bsonType: "string",
            description: "must be a string and is required"
        },
        email: {
          bsonType: "string",
          pattern: '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$',
          description: "must be a valid email and is required"
        },
        password: {
          bsonType: "string",
          description: "must be a string and is required"
        },
        created_at: {
          bsonType: "date",
          description: "must be a date and is required"
        },
        updated_at: {
          bsonType: "date",
          description: "must be a date and is required"
        }
      }
    }
  }
});

// Create unique index on email field
db.users.createIndex({ email: 1 }, { unique: true });

print("Database initialized");
