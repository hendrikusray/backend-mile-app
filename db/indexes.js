db = db.getSiblingDB(process.env.MONGO_INITDB_DATABASE || "local");

print("==> Creating MongoDB indexes...");

db.tasks.createIndex({ owner_id: 1, status: 1, created_at: -1 });
db.tasks.createIndex({ title: "text", description: "text" });

// users
db.users.createIndex({ username: 1 }, { unique: true });

print("==> Index creation complete ✅");
