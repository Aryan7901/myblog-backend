const { validationResult } = require("express-validator");
const HttpError = require("../models/http-error");
const User = require("../models/user");
const Blog = require("../models/blog");
const mongoose = require("mongoose");
const createBlogPage = async (req, res, next) => {
  const errors = validationResult(req);
  if (!errors.isEmpty()) {
    return next(
      new HttpError("Invalid inputs passed, please check your data.", 422)
    );
  }
  const { title, description, article } = req.body;
  const userId = req.userData.userId;
  const createdBlog = new Blog({
    title,
    description,
    article,
    author: userId,
  });
  try {
    const session = await mongoose.startSession();
    session.startTransaction();
    await createdBlog.save({ session });
    const user = await User.findById(userId);
    user.blogs.push(createdBlog);
    await user.save({ session });
    session.commitTransaction();
  } catch (err) {
    const error = new HttpError(
      "Creating new blog failed, please try again later.",
      500
    );
    return next(error);
  }
  res.send({ createdBlog });
};
const updateBlogPage = async (req, res, next) => {
  const errors = validationResult(req);
  if (!errors.isEmpty()) {
    return next(
      new HttpError("Invalid inputs passed, please check your data.", 422)
    );
  }
  const { title, description, article } = req.body;
  const blogId = req.params.bid;
  try {
    const blog = await Blog.findById(blogId);
    if (blog.author != req.userData.userId) {
      return next(new HttpError("Unauthorized!", 401));
    }
    blog.title = title;
    blog.description = description;
    blog.article = article;
    await blog.save();
    res.send({ blog });
  } catch (err) {
    return next(
      new HttpError("Updating  blog failed, please try again later.", 500)
    );
  }
};
const deleteBlogPage = async (req, res, next) => {
  const blogId = req.params.bid;
  try {
    const session = await mongoose.startSession();
    const blog = await Blog.findById(blogId);
    if (blog.author != req.userData.userId) {
      return next(new HttpError("Unauthorized!", 401));
    }

    session.startTransaction();
    await Blog.findByIdAndDelete(blogId).setOptions({ session });
    await User.updateOne(
      { blogs: blogId },
      { $pull: { blogs: blogId } }
    ).setOptions({ session });
    await session.commitTransaction();

    res.send({ blog });
  } catch (err) {
    console.log(err);
    return next(
      new HttpError("Deleting blog failed, please try again later.", 500)
    );
  }
};
const getUsersPages = async (req, res, next) => {
  const id = req.userData.userId;
  try {
    const user = await User.findById(id).populate("blogs");
    res.send({
      blogs: user.blogs,
      author: { firstName: user.firstName, lastName: user.lastName },
    });
  } catch (error) {
    return next(new HttpError("failed to get user's blogs", 500));
  }
};
exports.createBlogPage = createBlogPage;
exports.updateBlogPage = updateBlogPage;
exports.deleteBlogPage = deleteBlogPage;
exports.getUsersPages = getUsersPages;
