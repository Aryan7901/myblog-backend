const express = require("express");
const router = express.Router();
const blogControllers = require("../controllers/blog-controller");
router.get("/all", blogControllers.getAllBlogs);
router.get("/blog/:bid", blogControllers.getBlogById);

module.exports = router;
