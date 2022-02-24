const { Schema, model } = require("mongoose");

const blogSchema = new Schema({
  title: { type: String, required: true },
  author: {
    type: Schema.Types.ObjectId,
    ref: "User",
    required: true,
  },
  description: { type: String, required: true },
  article: { type: String, required: true, min: 500 },
  comments: [{ type: Schema.Types.ObjectId, ref: "Comment" }],
});
module.exports = model("Blog", blogSchema);
