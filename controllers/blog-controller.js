const HttpError = require("../models/http-error");
const Blog = require("../models/blog");
const { validationResult } = require("express-validator");

const getAllBlogs = async (req, res, next) => {
  const errors = validationResult(req);
  let blogs;
  try {
    blogs = await Blog.find({}).populate("author", "firstName lastName");
  } catch (err) {
    const error = new HttpError(
      "Error Retreiving data, please try again later.",
      500
    );
    return next(error);
  }
  res.json({ blogs });
};

const getBlogById = async (req, res, next) => {
  const blogId = req.params.bid;
  try {
    const blog = await Blog.findById(blogId).populate(
      "author",
      "firstName lastName"
    );
    res.json({ blog });
  } catch (err) {
    const error = new HttpError(
      "Error Retreiving blog, please try again later.",
      500
    );
    return next(error);
  }
};
exports.getAllBlogs = getAllBlogs;
exports.getBlogById = getBlogById;
