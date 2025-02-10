const puppeteer = require("puppeteer");
const { MongoClient, ObjectId } = require("mongodb");
const { isErrnoException } = require("puppeteer");
const _ = require("lodash");

const uri =
  "mongodb://localhost:27017/test_db"; // Replace with your MongoDB connection string
const client = new MongoClient(uri);

const productDetailBody = {
  name: "Zip-Pocket Joggers Navy Blue",
  imgLink:
    "https://nobero.com/cdn/shop/files/NavyBlue2_1331027c-0619-4aee-9a3a-08c9105412fd.jpg?v=1735808866",
  description: {
    productDetails:
      "<p><strong>Product Description</strong></p><ul><li>Look dapper and stay comfortable with the Nobero Men's Jogger.</li><li>Made with Cotton Rich Blend, 240 GSM, French Terry that goes through rigorous tests to meet high-quality standards, it can be worn all year round.</li><li>Fabrics used are procured from locals to support small businesses &amp; our notion \"VOCAL FOR LOCAL\".</li><li>Exclusively designed and manufactured by Nobero in Tirupur, India.</li><li>Get 100% money-back if you do not like our product. No questions asked. Shop with confidence!</li></ul><p><strong>Fabric Material</strong></p><ul><li>Cotton Rich Blend: We use the softest, most breathable and finest quality combed cotton.</li><li>240gsm: With a high GSM number, the fabric will be denser. 240 GSM is thick, weighs more and will make the apparel durable for better use.</li><li>French terry: It is a knitted terry cloth fabric that is soft and smooth on one side (generally, outer) and has loops and yarns on the other side. It is a wrinkle-free, lightweight, absorbent, moisture-wicking fabric that is extremely comfortable to be worn in any season</li></ul><p><strong>Style &amp; Fit</strong></p><ul><li>The&nbsp;Nobero Men's&nbsp;joggers may look and feel like they're ready to lounge in, but don't let their soft side fool you. These classic-meets-modern must-haves are engineered to keep up with even the toughest athletes.</li></ul>",
    specifications: {
      sleeveLength: "sleeveLength",
      collar: "collar",
      fit: "fit",
      patternType: "patternType",
      occasion: "occasion",
      length: "length",
      hemline: "hemline",
      placket: "placket",
      placketLength: "placketLength",
      cuff: "cuff",
      transparency: "transparency",
      weavePattern: "weavePattern",
      mainTrend: "mainTrend",
      numberOfItems: 1,
      packageContains: "packageContains",
    },
  },
  variations: [],
  imgLink: [],
  batchId: "975867",
  createdAt: new Date(),
  updatedAt: new Date(),
  createdBy: new ObjectId("67987164cf21db640ba92dc9"),
  modifiedBy: new ObjectId("67987164cf21db640ba92dc9"),
  slug: "zip-pocket-joggers-dark-blue",
  code: "8758697",
};

const productListBody = {
  name: "Zip-Pocket Joggers Light Blue",
  sizes: ["S", "M", "L", "XL", "XXL", "XXXL", "base_size"],
  color: "Blue",
  price: 800,
  imgLink: "https://nobero.com/cdn/shop/files/PowderBlue1.jpg?v=1735808835",
  stock: 10,
  discount: 50,
  detail: new ObjectId("679b66b8e9fb40ff1ec07acc"),
  category: new ObjectId("6798b36baeb9c16dd25aa383"),
  batchId: "975867",
  gender: "male",
  collection: ["joggers", "latest", "best-seller", "trending"],
  tags: ["name"],
  createdAt: new Date(),
  updatedAt: new Date(),
  createdBy: new ObjectId("67987164cf21db640ba92dc9"),
  modifiedBy: new ObjectId("67987164cf21db640ba92dc9"),
  slug: "zip-pocket-joggers-light-blue",
  code: "8758695",
};

const productBatchBody = {
  createdAt: new Date(),
  createdBy: new ObjectId("67987164cf21db640ba92dc9"),
  updatedAt: new Date(),
  modifiedBy: new ObjectId("67987164cf21db640ba92dc9"),
  code: "975867",
  name: "Zip Pocket Joggers",
};

const sizes = ["S", "M", "L", "XL", "XXL", "XXXL"];
const collection = ["latest", "best-sellers", "trending"];
const category = new Object("67a8a6876b836b42a9caa870")

async function start(websiteUrl) {
  await client.connect();
  const database = client.db();

  const browser = await puppeteer.launch({ headless: true });
  const page = await browser.newPage();

  try {
    const mainUrl = websiteUrl; // Replace with your website URL
    await page.goto(mainUrl, { waitUntil: "networkidle2" });

    // Get all anchor tags with class "product_link"
    const productLinks = await page.$$eval("a.product_link", (anchors) =>
      anchors.map((anchor) => anchor.href)
    );

    let imageUrls = [];

    for (const link of productLinks) {
      await page.goto(link, { waitUntil: "networkidle2" });

      // Wait for the first image inside the .prodLordSlider-element-container
      await page.waitForSelector(".prodLordSlider-element-container img", {
        timeout: 5000,
      });

      // Get image URLs
      const productData = await page.evaluate(() => {
        const name =
          document.querySelector("h1.product-title")?.innerText.trim() || "";
        const images = Array.from(
          document.querySelectorAll(".prodLordSlider-element-container img")
        ).map((img) => img.src); // Get src attribute of each image

        const colors = Array.from(
          document.querySelectorAll("fieldset.product-form-input label input")
        ).map((input) => input.value); // Get value attribute from each color input

        return { images, colors, name };
      });

      
      await formatDataAndSaveInDb(database, productData);
    }

    console.log("process has been ended")
  } catch (error) {
    console.error("Error:", error);
  } finally {
    await browser.close();
  }
}

async function formatDataAndSaveInDb(database, productData) {
  const productListCollection = database.collection("product_list");
  const productDetailsCollection = database.collection("product_detail");
  const productBatchCollection = database.collection("product_batch");

  let imgs = [];
  let finalData = [];
  let prevLength = 0;
  let name = productData.name;
  let colors = productData.colors;

  let prevBatchName = "";
  let productBatchObjectId = "";
  let batchId = "";

  for (let item of productData.images) {
    let itemLength = item.length;
    if (prevLength === itemLength) {
      imgs.push(item);
    } else {
      if (imgs.length >= 3) {
        const color = colors[finalData.length];
        finalData.push({
          imgs: imgs,
          name: `${name} ${color}`,
          batchName: name,
          color: color,
        });
      }
      imgs = [item];
    }
    prevLength = itemLength;
  }


  for (let item of finalData) {
    const productDetailData = _.cloneDeep(productDetailBody);
    const productListData = _.cloneDeep(productListBody);
    const productBatchData = _.cloneDeep(productBatchBody);

    const productcode = getRandomNumberOfLength(6).toString();
    const slug = sentenceToSlug(item.name);

    if (prevBatchName !== item.batchName) {
      batchId = getRandomNumberOfLength(5).toString();
      productBatchData.code = batchId;
      productBatchData.name = item.batchName;
      const res = await productBatchCollection.insertOne(productBatchData);
      prevBatchName = item.batchName
      productBatchObjectId = res.insertedId;
    }

    productDetailData.batchId = batchId;
    productDetailData.code = productcode;
    productDetailData.imgLink = item.imgs;
    productDetailData.name = item.name;
    productDetailData.slug = slug;
    productDetailData.variations = [];

    for (let item of sizes) {
      productDetailData.variations.push({
        size: item,
        price: getRandomNumberInRange(500, 3000),
        discount: getRandomNumberInRange(4, 8) * 10,
        stock: getRandomNumberInRange(10, 50),
      });
    }

    const productDetailRes = await productDetailsCollection.insertOne(
      productDetailData
    );
    const detailId = productDetailRes.insertedId;

    const colletion = collection[getRandomNumberInRange(0, 5)];

    productListData.batchId = batchId;
    productListData.sizes = sizes;
    productListData.category = category;
    productListData.code = productcode;
    productListData.collection = colletion ? [colletion] : [];
    productListData.color = item.color;
    productListData.discount = productDetailData.variations[0].discount;
    productListData.price = productDetailData.variations[0].price;
    productListData.gender = "male";
    productListData.imgLink = item.imgs[0];
    productListData.name = item.name;
    productListData.sizes = sizes;
    productListData.slug = slug;
    productListData.stock = productDetailData.variations[0].stock;
    productListData.tags = [];
    productListData.detail = detailId;

    await productListCollection.insertOne(productListData);

    console.log(productBatchObjectId.toString())

    await productBatchCollection.updateOne(
      { _id: new ObjectId(productBatchObjectId) }, // Match the document by _id
      {
        $push: {
          batchProductDetails: {
            productCode: productcode,
            imgLink: item.imgs[0],
            slug: slug,
          },
        },
      }
    );
  }
}

function getRandomNumberInRange(min, max) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function sentenceToSlug(sentence) {
  return sentence
    .toLowerCase() // Convert to lowercase
    .trim() // Remove leading/trailing spaces
    .replace(/[\s\W-]+/g, "-") // Replace spaces and non-word characters with a single hyphen
    .replace(/^-+|-+$/g, ""); // Remove hyphens from the start or end
}

function getRandomNumberOfLength(length) {
  if (length <= 0) return 0; // Return 0 if length is invalid
  const min = Math.pow(10, length - 1); // Minimum value for the given length
  const max = Math.pow(10, length) - 1; // Maximum value for the given length
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

const websiteUrlToScrape = "https://nobero.com/collections/fashion-joggers-men";

start(websiteUrlToScrape);
