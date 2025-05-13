db = db.getSiblingDB('animalsys');

db.createUser({
  user: 'animalsysuser',
  pwd: 'animalsyspass',
  roles: [{ role: 'readWrite', db: 'animalsys' }]
});

// Users collection schema
const userSchema = {
  $jsonSchema: {
    bsonType: 'object',
    required: ['username', 'email', 'password_hash', 'role'],
    properties: {
      username: { bsonType: 'string' },
      email: { bsonType: 'string' },
      password_hash: { bsonType: 'string' },
      role: { enum: ['admin', 'employee', 'volunteer', 'user'] },
      ldap_id: { bsonType: ['string', 'null'] }
    }
  }
};
db.createCollection('users', { validator: userSchema });

// Animals collection schema
const animalSchema = {
  $jsonSchema: {
    bsonType: 'object',
    required: ['name', 'species', 'breed', 'age', 'health_history', 'status'],
    properties: {
      name: { bsonType: 'string' },
      species: { bsonType: 'string' },
      breed: { bsonType: 'string' },
      age: { bsonType: 'int' },
      health_history: { bsonType: 'array', items: { bsonType: 'string' } },
      status: { enum: ['available', 'adopted', 'deceased'] }
    }
  }
};
db.createCollection('animals', { validator: animalSchema });

// Adoptions collection schema
const adoptionSchema = {
  $jsonSchema: {
    bsonType: 'object',
    required: ['animal_id', 'user_id', 'status', 'application_data'],
    properties: {
      animal_id: { bsonType: 'objectId' },
      user_id: { bsonType: 'objectId' },
      status: { enum: ['pending', 'approved', 'rejected'] },
      application_data: { bsonType: 'object' },
      contract_document_id: { bsonType: ['objectId', 'null'] }
    }
  }
};
db.createCollection('adoptions', { validator: adoptionSchema });

// Schedules collection schema
const scheduleSchema = {
  $jsonSchema: {
    bsonType: 'object',
    required: ['employee_id', 'shift_date', 'shift_time', 'tasks'],
    properties: {
      employee_id: { bsonType: 'objectId' },
      shift_date: { bsonType: 'date' },
      shift_time: { bsonType: 'string' },
      tasks: { bsonType: 'array', items: { bsonType: 'string' } }
    }
  }
};
db.createCollection('schedules', { validator: scheduleSchema }); 