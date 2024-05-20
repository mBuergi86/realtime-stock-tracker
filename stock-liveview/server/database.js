const { MongoClient } = require("mongodb");

let client;

const main = async () => {
  const uri =
    "mongodb://stockmarket:supersecret123@mongodb:27017,mongodb:27018,mongodb:27019/?replicaSet=rs0";
  client = new MongoClient(uri, {
    useNewUrlParser: true,
    useUnifiedTopology: true,
    maxPoolSize: 10,
    serverSelectionTimeoutMS: 5000,
    socketTimeoutMS: 45000,
  });

  try {
    console.log("Connecting to MongoDB");
    await client.connect();
    console.log("Connected to MongoDB");
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

module.exports = { main, getMongoClient };
