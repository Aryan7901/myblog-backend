const { Schema, model } = require("mongoose");

const commentSchema = new Schema({
  user: { type: Schema.Types.ObjectId, required: true, ref: "User" },
  content: { type: String, required: true },
  date: {
    type: Date,
    required: true,
  },
  blog: {
    type: Schema.Types.ObjectId,
    required: true,
    ref: "Blog",
  },
});
module.exports = model("Comment", commentSchema);
