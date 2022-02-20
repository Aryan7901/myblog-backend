const { Schema, model } = require("mongoose");

const blogSchema = new Schema({
  title: { type: String, required: true },
  author: {
    type: Schema.Types.ObjectId,
    ref: "User",
  },
  description: { type: String, required: true },
  article: { type: String, required: true, min: 500 },
});
module.exports = model("Blog", blogSchema);
