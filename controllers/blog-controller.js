const HttpError = require("../models/http-error");
const Blog = require("../models/blog");
const { validationResult } = require("express-validator");
const Comment = require("../models/comment");
const mongoose = require("mongoose");
const comment = require("../models/comment");
const getAllBlogs = async (req, res, next) => {
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
    const blog = await Blog.findById(blogId)
      .populate("author", "firstName lastName")
      .populate({
        path: "comments",
        populate: { path: "user", select: "firstName lastName" },
      });
    res.json({ blog });
  } catch (err) {
    const error = new HttpError(
      "Error Retreiving blog, please try again later.",
      500
    );
    return next(error);
  }
};
const makeComment = async (req, res, next) => {
  const errors = validationResult(req);

  if (!errors.isEmpty()) {
    return next(
      new HttpError("Invalid inputs passed, please check your data.", 422)
    );
  }
  const blogId = req.params.bid;
  let blog;
  let createdComment;
  try {
    blog = await Blog.findById(blogId);
    createdComment = new Comment({
      user: req.userData.userId,
      content: req.body.comment,
      date: new Date(),
      blog: blog,
    });
  } catch (err) {
    console.log(error.message);
    const error = new HttpError(
      " Could not find blog, please try again later.",
      500
    );
    return next(error);
  }

  try {
    const session = await mongoose.startSession();
    session.startTransaction();

    await createdComment.save({ session });
    blog.comments.push(createdComment);
    await blog.save({ session });
    await session.commitTransaction();
  } catch (err) {
    const error = new HttpError(
      "Adding comment failed, please try again later.",
      500
    );
    return next(error);
  }
  res.status(201).send("Comment Created!");
};
const updateComment = async (req, res, next) => {
  const errors = validationResult(req);

  if (!errors.isEmpty()) {
    return next(
      new HttpError("Invalid inputs passed, please check your data.", 422)
    );
  }
  const commentId = req.params.cid;
  try {
    const comment = await Comment.findById(commentId);
    if (comment.user != req.userData.userId) {
      return next(new HttpError("Unauthorized!", 401));
    }
    comment.content = req.body.comment;
    await comment.save();
    res.json(comment);
  } catch (err) {
    const error = new HttpError(
      "Updating comment failed, please try again later.",
      500
    );
    return next(error);
  }
};

const deleteComment = async (req, res, next) => {
  const commentId = req.params.cid;
  try {
    const comment = await Comment.findById(commentId);

    if (comment.user != req.userData.userId) {
      return next(new HttpError("Unauthorized!", 401));
    }
    const session = await mongoose.startSession();
    session.startTransaction();
    Comment.findByIdAndDelete(commentId).setOptions({ session });
    await Blog.updateOne(
      { comments: commentId },
      { $pull: { comments: commentId } }
    ).setOptions({ session });
    await session.commitTransaction();
    res.json(comment);
  } catch (err) {
    const error = new HttpError(
      "Deleting comment failed, please try again later.",
      500
    );
    return next(error);
  }
};
exports.getAllBlogs = getAllBlogs;
exports.getBlogById = getBlogById;
exports.makeComment = makeComment;
exports.updateComment = updateComment;
exports.deleteComment = deleteComment;
