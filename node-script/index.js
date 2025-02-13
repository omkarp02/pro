const { MongoClient } = require('mongodb');

const sourceURI = 'mongodb+srv://opwebdev:Omkar^100@omkar.iuqcpfi.mongodb.net/test_db';      // Change to your source DB URI
const destinationURI = 'mongodb://localhost:27017/test_db';  // Change to your destination DB URI

async function copyCollections() {
  const sourceClient = new MongoClient(sourceURI);
  const destinationClient = new MongoClient(destinationURI);

  try {
    await sourceClient.connect();
    await destinationClient.connect();

    const sourceDB = sourceClient.db();
    const destinationDB = destinationClient.db();

    // Get all collection names from the source database
    const collections = await sourceDB.listCollections().toArray();

    console.log(`Found ${collections.length} collections in sourceDB.`);

    // Drop all collections in the destination database
    const destCollections = await destinationDB.listCollections().toArray();
    for (let col of destCollections) {
      console.log(`Dropping collection: ${col.name}`);
      await destinationDB.collection(col.name).drop();
    }

    // Copy each collection from source to destination
    for (let { name: collectionName } of collections) {
      console.log(`Copying collection: ${collectionName}`);
      const sourceCollection = sourceDB.collection(collectionName);
      const destinationCollection = destinationDB.collection(collectionName);

      const documents = await sourceCollection.find().toArray();
      if (documents.length > 0) {
        await destinationCollection.insertMany(documents);
      }
    }

    console.log('Successfully copied all collections.');
  } catch (error) {
    console.error('Error during database copy:', error);
  } finally {
    await sourceClient.close();
    await destinationClient.close();
  }
}

copyCollections();
