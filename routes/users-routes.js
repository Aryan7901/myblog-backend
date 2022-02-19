const express = require("express");
const { check } = require("express-validator");
const router = require("./blog-routes");
const checkAuth = require("../middleware/check-auth");
const userController = require("../controllers/user-controller");
const blogCrudController = require("../controllers/blog-crud-controller");
router.post(
  "/signup",
  [
    check("firstName").not().isEmpty(),
    check("lastName").not().isEmpty(),
    check("email").normalizeEmail().isEmail(),
    check("password").isLength({ min: 6 }),
  ],
  userController.signup
);

router.post("/login", userController.login);
router.use(checkAuth);

router.get("/list", blogCrudController.getUsersPages);
router.post(
  "/new-blog",
  [
    check("title").notEmpty(),
    check("description").notEmpty(),
    check("article").isLength({ min: 50 }),
  ],
  blogCrudController.createBlogPage
);

router.patch(
  "/:bid",
  [
    check("title").notEmpty(),
    check("description").notEmpty(),
    check("article").isLength({ min: 50 }),
  ],
  blogCrudController.updateBlogPage
);
router.delete("/:bid", blogCrudController.deleteBlogPage);
module.exports = router;
