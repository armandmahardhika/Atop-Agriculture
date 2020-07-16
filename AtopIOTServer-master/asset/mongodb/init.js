// $ mongo localhost:27017/atop init.js
db.users.insert({
  name: "admin",
  password: "default",
  role: "admin",
  tags: "#admin",
});

db.searchOption.insert({
  collection: "users",
  searchableFields: ["name", "tags", "role"],
});

db.tableSetting.insert({
  name: "users",
  unique: ["name"],
});
// uniqu field settings
// db.users.createIndex({ ip: 1 }, { unique: true });
