const { Schema, model } = require("mongoose");
const uniqueValidator = require("mongoose-unique-validator");
const userSchema = new Schema({
  firstName: { type: String, required: true },
  lastName: { type: String, required: true },
  blogs: [
    {
      type: Schema.Types.ObjectId,
      ref: "Blog",
    },
  ],
  email: {
    type: String,
    required: true,
    match: [
      /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/,
      "Please fill a valid email address",
    ],
  },
  password: { type: String, required: true, min: 8 },
});
userSchema.plugin(uniqueValidator);
module.exports = model("User", userSchema);
