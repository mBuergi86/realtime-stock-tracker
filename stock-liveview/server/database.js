const { MongoClient } = require("mongodb");

let client;

const main = async () => {
  const uri =
    "mongodb://127.0.0.1:27017/,127.0.0.1:27018,127.0.0.1:27019/?replicaSet=rs0";
  client = new MongoClient(uri);

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
