const puppeteer = require("puppeteer");
const { MongoClient, ObjectId } = require("mongodb");
const { isErrnoException } = require("puppeteer");
const _ = require("lodash");

const uri =
  "mongodb+srv://opwebdev:Omkar^100@omkar.iuqcpfi.mongodb.net/test_db"; // Replace with your MongoDB connection string
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
const category = new Object("67a8b6106b836b42a9caa881");

async function start(websiteUrl) {
  await client.connect();
  const database = client.db();

  const browser = await puppeteer.launch({ headless: true });
  const page = await browser.newPage();

  try {
    const mainUrl = websiteUrl;
    await page.goto(mainUrl, { waitUntil: "networkidle2" });

    let productLinks = new Set();

    let previousHeight = 0;
    let maxScrollAttempts = 20; // Maximum attempts to scroll and fetch new content
    let scrollAttempts = 0;
    let breakk = true

    while ((scrollAttempts < maxScrollAttempts) && breakk) {
      // Scroll down
      const currentHeight = await page.evaluate(() => {
        window.scrollBy(0, window.innerHeight);
        return document.body.scrollHeight;
      });

      // Wait for content to load
      await new Promise((resolve) => setTimeout(resolve, 2000));

      // Get all product links
      const newLinks = await page.$$eval("a.product_link", (anchors) =>
        anchors.map((anchor) => anchor.href)
      );

      newLinks.forEach((link) => productLinks.add(link));

      // Check if we've reached the end
      if (previousHeight === currentHeight) {
        scrollAttempts++;
      } else {
        scrollAttempts = 0; // Reset attempts if new content is found
      }
      previousHeight = currentHeight;
      // breakk = false

      console.log(`Fetched ${productLinks.size} links so far...`);
    }

    console.log(`Total product links fetched: ${productLinks.size}`);

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

    console.log("process has been ended");
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

  let finalData = [];
  let name = productData.name;
  let colors = productData.colors;
  let count = 0;


  const obj = {};

  for (let item of productData.images) {
    const myURL = new URL(item);
    const query = myURL.searchParams.get("v");
    if (obj[query]) {
      obj[query].imgs.push(item);
    } else {
      obj[query] = {
        imgs: [item],
        name: `${name} ${colors[count]}`,
        batchName: name,
        color: colors[count],
      };
      count++
    }
  }

  for (let key in obj) {
    finalData.push(obj[key]);
  }

  const productBatchData = _.cloneDeep(productBatchBody);
  const batchId = getRandomNumberOfLength(5).toString();
  productBatchData.code = batchId;
  productBatchData.name = name;
  const res = await productBatchCollection.insertOne(productBatchData);
  const productBatchObjectId = res.insertedId;

  for (let item of finalData) {
    const productDetailData = _.cloneDeep(productDetailBody);
    const productListData = _.cloneDeep(productListBody);

    const productcode = getRandomNumberOfLength(6).toString();
    const slug = sentenceToSlug(item.name);


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

    console.log(productBatchObjectId.toString());

    await productBatchCollection.updateOne(
      { _id: productBatchObjectId }, // Match the document by _id
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

const websiteUrlToScrape =
  "https://nobero.com/collections/pick-printed-t-shirts";

start(websiteUrlToScrape);

const productData = {
  images: [
    "https://nobero.com/cdn/shop/files/5_86dd2609-2ecd-4ba1-8485-e82c47d7496b.jpg?v=1735322751",
    "https://nobero.com/cdn/shop/files/WhatsAppImage2024-07-29at7.13.12PM_1_cd5fcdc8-de0b-4788-99da-4553f3e41049.jpg?v=1736680332",
    "https://nobero.com/cdn/shop/files/WhatsAppImage2024-07-29at7.13.12PM_df1865e1-e237-44ca-ae9a-801142235b08.jpg?v=1735322751",
    "https://nobero.com/cdn/shop/files/WhatsAppImage2024-08-05at1.45.51PM_1_fdbac966-157c-48c2-b5be-2bbd2d7e1a3b.jpg?v=1736680361",
    "https://nobero.com/cdn/shop/files/2-2_1027c3dd-0ddb-4d02-af4f-763a947bd9f7.jpg?v=1736680361",
    "https://nobero.com/cdn/shop/files/10_100da693-3493-42b4-ac03-efce292eec0d.jpg?v=1736680361",
    "https://nobero.com/cdn/shop/files/1_b3544e25-d2b1-4be7-aa24-b1563d2c8d76.jpg?v=1736680361",
    "https://nobero.com/cdn/shop/files/2_47f0e462-0bd0-4b8e-ac5e-6345ebf6262b.jpg?v=1736680361",
    "https://nobero.com/cdn/shop/files/15_22d23a67-583c-44e5-a4db-371f78375c10.jpg?v=1736680361",
    "https://nobero.com/cdn/shop/files/2-11_09157623-8438-4e26-9ac8-6cb9880f956e.jpg?v=1736680332",
    "https://nobero.com/cdn/shop/files/2-8_ebb94411-8f0e-4961-9a24-77b0ce95460e.jpg?v=1736680332",
    "https://nobero.com/cdn/shop/files/3_1a6bc7ae-510f-4b33-9662-8375d347b16f.jpg?v=1735322751",
    "https://nobero.com/cdn/shop/files/5_86dd2609-2ecd-4ba1-8485-e82c47d7496b.jpg?v=1735322751",
    "https://nobero.com/cdn/shop/files/WhatsAppImage2024-07-29at7.13.12PM_1_cd5fcdc8-de0b-4788-99da-4553f3e41049.jpg?v=1736680332",
  ],
  colors: [
    "Authentic/ And/ Never Say No",
    "Wandersoul/ Take a Break/ Be Free",
    "Authentic/ Wandersoul/ Sunset",
    "Wabi Sabi/ Sunset/ Wild",
    "Black/ Olive Green/ Wine Red",
    "Black/ White/ Marine",
    "Olive Green/ Powder Blue/ Wine Red",
    "Wine Red/ Grey Melange/ Sand",
    "Restart/ Balance V3/ Mountains",
    "Authentic/ Chase/ Discover",
    "White/ Marine/ Wine Red",
    "Olive Green/ Wine Red/ White",
    "S",
    "M",
    "L",
    "XL",
    "XXL",
    "XXXL",
    "S",
    "M",
    "L",
    "XL",
    "XXL",
    "XXXL",
  ],
  name: "The Classic - 3 Pack",
};

