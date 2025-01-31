package test

// {{
// 	Key: "$match", Value: bson.D{
// 		{Key: "name", Value: nameToMatch},
// 	},
// }},

// {{
// 	Key: "$lookup", Value: bson.D{
// 		{Key: "from", Value: "state"},                  // The collection to join with
// 		{Key: "localField", Value: "stateId"},          // Field in the current collection
// 		{Key: "foreignField", Value: "_id"},            // Field in the "state" collection
// 		{Key: "as", Value: "stateDetails"},             // Output array field
// 	},
// }},

// {{
// 	Key: "$unwind", Value: bson.D{
// 		{Key: "path", Value: "$stateDetails"},
// 		{Key: "preserveNullAndEmptyArrays", Value: true},
// 	},
// }},

// {{
// 	Key: "$project", Value: bson.D{
// 		{Key: "name", Value: 1},                        // Include name
// 		{Key: "stateId", Value: 1},                     // Include stateId
// 		{Key: "stateDetails.state_name", Value: 1},     // Include state_name from stateDetails
// 	},
// }},
