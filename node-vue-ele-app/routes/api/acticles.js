// @login & register
const express = require("express");
const router = express.Router();
const passport = require('passport');

const Acticle = require("../../models/Acticle");

// $route GET api/acticles/test
// @desc  返回請求的json數據
// @access public
router.get("/test", (req, res) => {
    res.json({msg: "acticles works"})
});

// $route POST api/acticles/add
// @desc  創建信息接口
// @access private
router.post("/add", passport.authenticate("jwt", {session: false}),
(req, res) => {
    const acticleFields = {};

    if (req.body.type) acticleFields.type = req.body.type;
    if (req.body.describe) acticleFields.describe = req.body.describe;
    if (req.body.remark) acticleFields.remark = req.body.remark;

    new Acticle(acticleFields).save().then(acticle => {
        res.json(acticle);
    })
});

// $route GET api/acticles
// @desc  獲取所有訊息
// @access private
router.get("/", passport.authenticate("jwt", {session: false}),
(req, res) => {
    Acticle.find()
        .then(acticle => {
            if (!acticle) {
                return res.status(404).json('沒有任何內容');
            }
            res.json(acticle);
        })
        .catch((err) => {
            res.status(404).json(err);
        });
});

// $route GET api/acticles/:id
// @desc  獲取單一訊息
// @access private
router.get("/:id", passport.authenticate("jwt", {session: false}),
(req, res) => {
    Acticle.findOne({_id: req.params.id})
        .then(acticle => {
            if (!acticle) {
                return res.status(404).json('沒有任何內容');
            }
            res.json(acticle);
        })
        .catch((err) => {
            res.status(404).json(err);
        });
});

// $route POST api/acticles/edit
// @desc  編輯信息接口
// @access private
router.post("/edit/:id", passport.authenticate("jwt", {session: false}),
(req, res) => {
    const acticleFields = {};

    if (req.body.type) acticleFields.type = req.body.type;
    if (req.body.describe) acticleFields.describe = req.body.describe;
    if (req.body.remark) acticleFields.remark = req.body.remark;

    Acticle.findOneAndUpdate(
        {_id: req.params.id},
        {$set: acticleFields},
        {new: true}
    ).then(acticle => res.json(acticle))
});

// $route POST api/acticles/delete
// @desc  刪除信息接口
// @access private
router.delete("/delete/:id", passport.authenticate("jwt", {session: false}),
(req, res) => {
    Acticle.findOneAndRemove({_id: req.params.id})
        .then(acticle => {
            acticle.save().then(acticle => res.json(acticle));
        })
        .catch((err) => res.status(404).json('刪除失敗!'));
});

module.exports = router;