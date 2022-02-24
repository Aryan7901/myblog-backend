const express = require("express");
const { check } = require("express-validator");
const router = express.Router();
const blogControllers = require("../controllers/blog-controller");
const checkAuth = require("../middleware/check-auth");
router.get("/all", blogControllers.getAllBlogs);
router.get("/blog/:bid", blogControllers.getBlogById);
router.use(checkAuth);
router.post(
  "/comment/:bid",
  check("comment").notEmpty(),
  blogControllers.makeComment
);
router.patch(
  "/comment/:cid",
  check("comment").notEmpty(),
  blogControllers.updateComment
);
router.delete("/comment/:cid", blogControllers.deleteComment);
module.exports = router;
