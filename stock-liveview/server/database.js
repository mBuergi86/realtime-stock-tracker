const { MongoClient } = require("mongodb");

let client;

const initializeDatabase = async () => {
  const uri =
    "mongodb://mongodb:27017,mongodb:27018,mongodb:27019/?replicaSet=rs0";
  client = new MongoClient(uri, {
    useNewUrlParser: true,
    useUnifiedTopology: true,
    maxPoolSize: 10,
    serverSelectionTimeoutMS: 5000,
    socketTimeoutMS: 45000,
  });

  try {
    console.log("Processing connection to MongoDB");
    await client.connect();
    console.log("Successfully connected to MongoDB");
    await listDatabases(client);
  } catch (error) {
    console.error(error);
  }
};

const listDatabases = async (client) => {
  const databasesList = await client.db().admin().listDatabases();
  console.log("Databases:");
  databasesList.databases.forEach((db) => console.log(` - ${db.name}`));
};

const getMongoClient = () => {
  return client;
};

module.exports = { initializeDatabase, getMongoClient };
